package extqueue

import (
	"sync/atomic"
)

// MPSCrwMC is a MPSC queue using disruptor style waiting on the producer side
// and MCRingBuffer style consumer batching.
type MPSCrsMC struct {
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
}

// NewMPSCrsMC creates a new MPSCrsMC queue
func NewMPSCrsMC(batchSize, size int) *MPSCrsMC {
	q := &MPSCrsMC{}
	q.batchSize = int64(batchSize)
	if size < batchSize {
		size = batchSize
	}
	q.buffer = make([]Value, int(nextPowerOfTwo(uint32(size))))
	q.mask = int64(len(q.buffer) - 1)

	return q
}

// Cap returns number of elements this queue can hold before blocking
func (q *MPSCrsMC) Cap() int { return len(q.buffer) }

// MultipleProducers makes this a MP queue
func (q *MPSCrsMC) MultipleProducers() {}

// Send sends a value to the queue and blocks when it is full
func (q *MPSCrsMC) Send(v Value) bool {
	// grab a write location
	writeTo := atomic.AddInt64(&q.writeTo, 1) - 1

	// channel is full, wait for it to drain
	for try := 0; atomic.LoadInt64(&q.nextRead)+q.mask < writeTo; spin(&try) {
	}

	// now we can write
	q.buffer[writeTo&q.mask] = v

	// wait for previous writes to complete
	for try := 0; writeTo != atomic.LoadInt64(&q.unwritten); spin(&try) {
	}

	atomic.StoreInt64(&q.unwritten, writeTo+1)

	return true
}

// FlushSend is to implement interface, on this queue this is a nop
func (q *MPSCrsMC) FlushSend() {}

// Recv receives a value from the queue and blocks when it is empty
func (q *MPSCrsMC) Recv(v *Value) bool { return q.recv(v, true) }

// TryRecv receives a value from the queue and returns when it is empty
func (q *MPSCrsMC) TryRecv(v *Value) bool { return q.recv(v, false) }

func (q *MPSCrsMC) recv(v *Value, block bool) bool {
	localUnwritten := q.localUnwritten
	for try := 0; q.localNextRead >= localUnwritten; spin(&try) {
		localUnwritten = atomic.LoadInt64(&q.unwritten)
		if !block {
			return false
		}
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
func (q *MPSCrsMC) FlushRecv() {
	atomic.StoreInt64(&q.nextRead, q.localNextRead)
	q.localReadBatch = 0
}
