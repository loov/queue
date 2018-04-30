package testsuite

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func testSPSC(t *testing.T, caps Capability, ctor func() Queue) {
	t.Run("Single", func(t *testing.T) {
		for _, count := range TestCount {
			q := ctor().(SPSC)
			for i := 0; i < count; i++ {
				exp := Value(i)
				q.Send(exp)
				var got Value
				q.Recv(&got)
				if exp != got {
					t.Fatalf("expected %v got %v", exp, got)
				}
			}
		}
	})

	t.Run("Basic", func(t *testing.T) {
		for _, count := range TestCount {
			q := ctor().(SPSC)
			ProducerConsumer(t, 1, 1, func(int) error {
				for i := 0; i < count; i++ {
					if !q.Send(Value(i + 1)) {
						return fmt.Errorf("failed to send %v", i)
					}
				}
				FlushSend(q)
				return nil
			}, func(int) error {
				for i := 0; i < count; i++ {
					exp := Value(i + 1)

					var got Value
					if !q.Recv(&got) {
						return fmt.Errorf("recv failed")
					}

					if got != exp {
						return fmt.Errorf("invalid value got %v, expected %v", got, exp)
					}
				}
				return nil
			})
		}
	})

	if caps.Has(CapBounded) {
		t.Run("BlockOnFull", func(t *testing.T) {
			q := ctor().(interface {
				SPSC
				Bounded
			})
			capacity := q.Cap()

			for i := 0; i < capacity; i++ {
				if !q.Send(0) {
					t.Fatal("failed to send")
				}
			}

			FlushSend(q)
			sent := uint32(0)
			go func() {
				if !q.Send(0) {
					t.Fatal("failed to send")
				}
				FlushSend(q)
				atomic.StoreUint32(&sent, 1)
			}()
			runtime.Gosched()
			time.Sleep(time.Millisecond)
			if atomic.LoadUint32(&sent) != 0 {
				t.Fatalf("send to full queue")
			}

			var v Value
			if !q.Recv(&v) {
				t.Fatal("failed to recv from full")
			}

			FlushRecv(q)

			runtime.Gosched()
			if atomic.LoadUint32(&sent) != 1 {
				runtime.Gosched()
				time.Sleep(time.Millisecond)
				if atomic.LoadUint32(&sent) != 1 {
					t.Fatalf("did not unblock blocked channel")
				}
			}
		})
	}
}

func testMPSC(t *testing.T, caps Capability, ctor func() Queue) {
	t.Run("Basic", func(t *testing.T) {
		for _, count := range TestCount {
			q := ctor().(MPSC)
			ProducerConsumer(t,
				TestProcs, 1,
				func(id int) error {
					for i := 0; i < count; i++ {
						if !q.Send(Value(id)<<32 | Value(i)) {
							return fmt.Errorf("failed to send %v", i)
						}
					}
					return nil
				}, func(int) error {
					exps := make([]Value, TestProcs)
					for i := 0; i < count*TestProcs; i++ {
						var val Value
						if !q.Recv(&val) {
							return fmt.Errorf("failed to get")
						}
						id, got := val>>32, val&0xFFFFFFFF
						exp := exps[id]
						exps[id]++
						if exp != got {
							return fmt.Errorf("invalid order got %v, expected %v", got, exp)
						}
					}
					return nil
				})
		}
	})
}

func testSPMC(t *testing.T, caps Capability, ctor func() Queue) {
	t.Run("Basic", func(t *testing.T) {
		for _, count := range TestCount {
			q := ctor().(SPMC)
			ProducerConsumer(t,
				1, TestProcs,
				func(int) error {
					for i := 0; i < count*TestProcs; i++ {
						if !q.Send(Value(i + 1)) {
							return fmt.Errorf("failed to send %v", i)
						}
					}
					FlushSend(q)
					return nil
				}, func(int) error {
					var lastexp Value
					for i := 0; i < count; i++ {
						var got Value
						if !q.Recv(&got) {
							return fmt.Errorf("failed to get")
						}
						exp := lastexp
						lastexp = got
						if got <= exp {
							return fmt.Errorf("invalid order got %v, expected at least %v", got, exp)
						}
					}
					return nil
				})
		}
	})
}

func testMPMC(t *testing.T, caps Capability, ctor func() Queue) {
	t.Run("SendRecv", func(t *testing.T) {
		for _, count := range TestCount {
			q := ctor().(MPMC)
			ProducerConsumer(t,
				TestProcs, 0,
				func(id int) error {
					latest := make([]Value, TestProcs)
					for i := 0; i < count; i++ {
						if !q.Send(Value(id)<<32 | Value(i+1)) {
							return fmt.Errorf("failed to send %v", i)
						}

						var val Value
						if !q.Recv(&val) {
							return fmt.Errorf("failed to get")
						}

						id, got := int(val>>32), val&0xFFFFFFFF
						exp := latest[id]
						latest[id] = got

						if val <= exp {
							return fmt.Errorf("expected larger %v got %v", exp, got)
						}
					}
					return nil
				}, nil)
		}
	})

	t.Run("Basic", func(t *testing.T) {
		for _, count := range TestCount {
			q := ctor().(MPMC)
			ProducerConsumer(t,
				TestProcs, TestProcs,
				func(id int) error {
					for i := 0; i < count; i++ {
						if !q.Send(Value(id)<<32 | Value(i+1)) {
							return fmt.Errorf("failed to send %v", i)
						}
					}
					FlushSend(q)
					return nil
				}, func(id int) error {
					latest := make([]Value, TestProcs)
					for i := 0; i < count; i++ {
						var val Value
						if !q.Recv(&val) {
							return fmt.Errorf("failed to get")
						}

						id, got := val>>32, val&0xFFFFFFFF
						exp := latest[id]
						latest[id] = got
						if got <= exp {
							return fmt.Errorf("invalid order got %v, expected %v", got, exp)
						}
					}
					return nil
				})
		}
	})
}

func testNonblockSPSC(t *testing.T, caps Capability, ctor func() Queue) {
	t.Run("Single", func(t *testing.T) {
		for _, count := range TestCount {
			q := ctor().(NonblockingSPSC)
			for i := 0; i < count; i++ {
				exp := Value(i)
				if !q.TrySend(exp) {
					t.Fatalf("send failed")
				}
				var got Value
				if !q.TryRecv(&got) {
					t.Fatalf("recv failed")
				}
				if exp != got {
					t.Fatalf("expected %v got %v", exp, got)
				}
			}
		}
	})

	t.Run("Basic", func(t *testing.T) {
		for _, count := range TestCount {
			q := ctor().(NonblockingSPSC)
			ProducerConsumer(t,
				1, 1,
				func(id int) error {
					for i := 0; i < count; i++ {
						if !MustSendIn(q, Value(i+1), time.Microsecond) {
							return fmt.Errorf("failed to send %v", i)
						}
					}
					FlushSend(q)
					return nil
				},
				func(id int) error {
					for i := 0; i < count; i++ {
						exp := Value(i + 1)

						var got Value
						if !MustRecvIn(q, &got, time.Microsecond) {
							return fmt.Errorf("recv timed out")
						}

						if got != exp {
							return fmt.Errorf("invalid value got %v, expected %v", got, exp)
						}
					}
					return nil
				})
		}
	})

	if caps.Has(CapBounded) {
		t.Run("NonblockOnFull", func(t *testing.T) {
			q := ctor().(NonblockingSPSC)
			capacity := Cap(q)
			for i := 0; i < capacity; i++ {
				if !q.TrySend(0) {
					t.Fatal("failed to send")
				}
			}
			FlushSend(q)
			if q.TrySend(0) {
				t.Fatal("send succeeded")
			}
			FlushSend(q)
		})
	}
}

func testNonblockMPSC(t *testing.T, caps Capability, ctor func() Queue) {}
func testNonblockSPMC(t *testing.T, caps Capability, ctor func() Queue) {}
func testNonblockMPMC(t *testing.T, caps Capability, ctor func() Queue) {}

/*
func testNonblockMPSC(t *testing.T, caps Capability, ctor func() Queue) {
	t.Run("Basic", func(t *testing.T) {
		for _, size := range TestSizes {
			q := ctor(size)
			count := Value(size/2) + 1

			result := async.All(func() error {
				return async.SpawnWithResult(TestProcs, func(id int) error {
					for i := Value(0); i < count; i++ {
						if !mustSend(q, Value(id)<<32|i) {
							return fmt.Errorf("%v: failed to send %v", size, i)
						}
					}
					flushsend(q)
					return nil
				}).WaitError()
			}, func() error {
				var exps [TestProcs]Value
				for i := Value(0); i < count*TestProcs; i++ {
					var val Value
					if !mustRecv(q, &val) {
						return fmt.Errorf("%v: failed to get", size)
					}
					id, got := val>>32, val&0xFFFFFFFF
					exp := exps[id]
					exps[id]++
					if got != exp {
						return fmt.Errorf("%v: invalid order got %v, expected %v", size, got, exp)
					}
				}
				return nil
			})

			if errs := result.Wait(); errs != nil {
				t.Fatal(errs)
			}
		}
	})
}
func testNonblockSPMC(t *testing.T, caps Capability, ctor func() Queue) {
	t.Run("Basic", func(t *testing.T) {
		for _, size := range TestSizes {
			q := ctor(size)
			count := Value(size/2) + 1

			result := async.All(func() error {
				for i := Value(0); i < count*TestProcs; i++ {
					if !mustSend(q, i) {
						return fmt.Errorf("%v: failed to send %v", size, i)
					}
				}
				flushsend(q)
				return nil
			}, func() error {
				return async.SpawnWithResult(TestProcs, func(int) error {
					var lastexp Value = -1
					for i := Value(0); i < count; i++ {
						var got Value
						if !mustRecv(q, &got) {
							return fmt.Errorf("%v: failed to get", size)
						}
						exp := lastexp
						lastexp = got
						if got <= exp {
							return fmt.Errorf("%v: invalid order got %v, expected %v", size, got, exp)
						}
					}
					return nil
				}).WaitError()
			})

			if errs := result.Wait(); errs != nil {
				t.Fatal(errs)
			}
		}
	})
}
func testNonblockMPMC(t *testing.T, caps Capability, ctor func() Queue) {
	t.Run("Basic", func(t *testing.T) {
		for _, size := range TestSizes {
			q := ctor(size)
			count := Value(size/2) + 1

			result := async.All(func() error {
				return async.SpawnWithResult(TestProcs, func(id int) error {
					for i := Value(0); i < count; i++ {
						if !mustSend(q, Value(id)<<32|i) {
							return fmt.Errorf("%v: failed to send %v", size, i)
						}
					}
					flushsend(q)
					return nil
				}).WaitError()
			}, func() error {
				return async.SpawnWithResult(TestProcs, func(int) error {
					var exps [TestProcs]Value
					for i := range exps {
						exps[i] = -1
					}
					for i := Value(0); i < count; i++ {
						var val Value
						if !mustRecv(q, &val) {
							return fmt.Errorf("%v: failed to get", size)
						}
						id, got := val>>32, val&0xFFFFFFFF
						exp := exps[id]
						exps[id] = got
						if got <= exp {
							return fmt.Errorf("%v: invalid order got %v, expected %v", size, got, exp)
						}
					}
					return nil
				}).WaitError()
			})

			if errs := result.Wait(); errs != nil {
				t.Fatal(errs)
			}
		}
	})
}
*/
