package testsuite

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

// LocalWork simulates work
func LocalWork(amount int) {
	foo := 1
	for i := 0; i < amount; i++ {
		foo *= 2
		foo /= 2
	}
}

// Cap returns queue capacity, when supported
// otherwise returns 1 << 32
func Cap(q Queue) int {
	if bounded, ok := q.(Bounded); ok {
		return bounded.Cap()
	}
	return 1 << 32
}

// FlushSend flushes send, if queue supports it
func FlushSend(q Queue) {
	if flusher, ok := q.(Flusher); ok {
		flusher.FlushSend()
	}
}

// FlushRecv flushes recv, if queue supports it
func FlushRecv(q Queue) {
	if flusher, ok := q.(Flusher); ok {
		flusher.FlushRecv()
	}
}

func MustSendIn(q NonblockingSPSC, v Value, dur time.Duration) bool {
	if q.TrySend(v) {
		return true
	}

	start := time.Now()
	for try := 0; ; try++ {
		if q.TrySend(v) {
			return true
		}

		if try > 256 {
			try = 0
			runtime.Gosched()
			if time.Since(start) > dur {
				return false
			}
		}
	}
	return false
}

func MustRecvIn(q NonblockingSPSC, v *Value, dur time.Duration) bool {
	if q.TryRecv(v) {
		return true
	}

	start := time.Now()
	for try := 0; ; try++ {
		if q.TryRecv(v) {
			return true
		}

		if try > 256 {
			try = 0
			runtime.Gosched()
			if time.Since(start) > dur {
				return false
			}
		}
	}
	return false
}

func ProducerConsumer(t *testing.T, NP, NC int, producer, consumer func(id int) error) {
	t.Helper()

	var wg sync.WaitGroup
	wg.Add(NP + NC)

	errs := make(chan error, NP+NC)
	for i := 0; i < NP+NC; i++ {
		go func(id int) {
			var err error

			defer func() {
				if rerr := recover(); rerr != nil {
					if err, iserr := rerr.(error); iserr {
						errs <- err
					} else {
						errs <- fmt.Errorf("%v", rerr)
					}
					return
				}
				errs <- err
			}()

			if id < NP {
				err = producer(id)
			} else {
				err = consumer(id)
			}
		}(i)
	}

	for i := 0; i < NP+NC; i++ {
		if err := <-errs; err != nil {
			t.Fatal(err)
		}
	}
}
