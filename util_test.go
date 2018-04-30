package queue

import (
	"fmt"
	"sync"
	"testing"
)

func Cap(q Queue) int {
	if bounded, ok := q.(Bounded); ok {
		return bounded.Cap()
	}
	return 1 << 32
}

func FlushSend(q Queue) {
	if flusher, ok := q.(Flusher); ok {
		flusher.FlushSend()
	}
}

func FlushRecv(q Queue) {
	if flusher, ok := q.(Flusher); ok {
		flusher.FlushRecv()
	}
}

func Parallel(t *testing.T, fns ...func() error) {
	var wg sync.WaitGroup
	wg.Add(len(fns))

	errs := make(chan error, len(fns))
	for i := 0; i < len(fns); i++ {
		go func(fn func() error) {
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

			err = fn()
		}(fns[i])
	}

	for i := 0; i < len(fns); i++ {
		if err := <-errs; err != nil {
			t.Fatal(err)
		}
	}
}
