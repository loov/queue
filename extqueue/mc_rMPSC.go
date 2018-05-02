package extqueue

import (
	"sync"
	"sync/atomic"
)

// MPSCrMC is a MPSC queue using condition variables on the producer side
// and MCRingBuffer style consumer batching.
//
// Not recommended.
type MPSCrMC struct {
	_ [8]uint64
	// volatile
	writeTo   int64
	nextRead  int64
	unwritten int64
	_         [8 - 3]uint64
	// consumer
	localUnwritten int64
	localNextRead  int64
	localReadBatch int64
	_              [8 - 3]uint64
	// constant
	batchSize int64
	mask      int64
	buffer    []Value
	// sleeping
	mu      sync.Mutex
	reader  sync.Cond
	writers sync.Cond
	drain   sync.Cond
}

// NewMPSCrMC creates a new MPSCrMC queue
func NewMPSCrMC(batchSize, size int) *MPSCrMC {
	q := &MPSCrMC{}
	q.reader.L = &q.mu
	q.writers.L = &q.mu
	q.drain.L = &q.mu

	q.batchSize = int64(batchSize)
	if size < batchSize {
		size = batchSize
	}
	q.buffer = make([]Value, int(nextPowerOfTwo(uint32(size))))
	q.mask = int64(len(q.buffer) - 1)
	return q
}

// Cap returns number of elements this queue can hold before blocking
func (q *MPSCrMC) Cap() int { return len(q.buffer) }

// MultipleProducers makes this a MP queue
func (q *MPSCrMC) MultipleProducers() {}

// Send sends a value to the queue and blocks when it is full
func (q *MPSCrMC) Send(v Value) bool {
	// grab a write location
	writeTo := atomic.AddInt64(&q.writeTo, 1) - 1

	// channel is full, wait for it to drain
	if atomic.LoadInt64(&q.nextRead)+q.mask < writeTo {
		q.mu.Lock()
		for q.nextRead+q.mask < writeTo {
			q.writers.Wait()
		}
		q.mu.Unlock()
	}

	// now we can write
	q.buffer[writeTo&q.mask] = v

	q.mu.Lock()
	for writeTo != q.unwritten {
		q.drain.Wait()
	}
	q.unwritten = writeTo + 1
	q.reader.Signal()
	q.drain.Broadcast()
	q.mu.Unlock()

	return true
}

// FlushSend is to implement interface, on this queue this is a nop
func (q *MPSCrMC) FlushSend() {}

// Recv receives a value from the queue and blocks when it is empty
func (q *MPSCrMC) Recv(v *Value) bool { return q.recv(v, true) }

// TryRecv receives a value from the queue and returns when it is empty
func (q *MPSCrMC) TryRecv(v *Value) bool { return q.recv(v, false) }

func (q *MPSCrMC) recv(v *Value, block bool) bool {
	localUnwritten := q.localUnwritten

	if q.localNextRead >= localUnwritten {
		q.mu.Lock()
		localUnwritten = atomic.LoadInt64(&q.unwritten)
		for q.localNextRead >= localUnwritten {
			if !block {
				q.mu.Unlock()
				return false
			}
			q.reader.Wait()
			localUnwritten = atomic.LoadInt64(&q.unwritten)
		}
		q.mu.Unlock()
	}

	q.localUnwritten = localUnwritten

	*v = q.buffer[q.localNextRead&q.mask]
	// q.buffer[q.localNextRead] = 0 // clear value, only needed for pointers

	q.localNextRead++
	q.localReadBatch++
	if q.localReadBatch >= q.batchSize {
		q.FlushRecv()
	}

	return true
}

// FlushRecv propagates pending receive operations to the sender.
func (q *MPSCrMC) FlushRecv() {
	q.mu.Lock()
	atomic.StoreInt64(&q.nextRead, q.localNextRead)
	q.localReadBatch = 0
	q.writers.Broadcast()
	q.mu.Unlock()
}
