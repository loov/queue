package queue

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func testSPSC(t *testing.T, caps Capability, ctor func() Queue) {
	t.Run("SendRecv", func(t *testing.T) {
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

	t.Run("Concurrent", func(t *testing.T) {
		for _, count := range TestCount {
			q := ctor().(SPSC)
			fmt.Println(Cap(q))
			Parallel(t, func() error {
				for i := 0; i < count; i++ {
					if !q.Send(Value(i + 1)) {
						return fmt.Errorf("failed to send %v", i)
					}
				}
				FlushSend(q)
				return nil
			}, func() error {
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

func testMPSC(t *testing.T, caps Capability, ctor func() Queue) {}

func testSPMC(t *testing.T, caps Capability, ctor func() Queue) {}

func testMPMC(t *testing.T, caps Capability, ctor func() Queue) {}
