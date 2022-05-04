package extqueue

import (
	"sync/atomic"
)

// MPMCqsDV is a MPMC queue based on http://www.1024cores.net/home/lock-free-algorithms/queues/bounded-mpmc-queue
type MPMCqsDV[T any] struct {
	_      [8]int64
	mask   int64
	buffer []seqValue[T]
	_      [4]int64
	sendx  int64
	_      [7]int64
	recvx  int64
	_      [7]int64
}

// NewMPMCqsDV creates a NewMPMCqsDV queue
func NewMPMCqsDV[T any](size int) *MPMCqsDV[T] {
	if size <= 1 {
		size = 2
	}
	size = int(nextPowerOfTwo(uint32(size)))

	q := &MPMCqsDV[T]{}
	q.buffer = make([]seqValue[T], size)
	q.mask = int64(size) - 1
	for i := range q.buffer {
		q.buffer[i].sequence = int64(i)
	}

	return q
}

// Cap returns number of elements this queue can hold before blocking
func (q *MPMCqsDV[T]) Cap() int { return len(q.buffer) }

// MultipleConsumers makes this a MC queue
func (q *MPMCqsDV[T]) MultipleConsumers() {}

// MultipleProducers makes this a MP queue
func (q *MPMCqsDV[T]) MultipleProducers() {}

// Send sends a value to the queue and blocks when it is full
func (q *MPMCqsDV[T]) Send(v T) bool {
	for wait := 0; ; spin(&wait) {
		if q.TrySend(v) {
			return true
		}
	}
}

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *MPMCqsDV[T]) TrySend(v T) bool {
	var cell *seqValue[T]
	pos := atomic.LoadInt64(&q.sendx)
	for {
		cell = &q.buffer[pos&q.mask]
		seq := atomic.LoadInt64(&cell.sequence)
		df := seq - pos
		if df == 0 {
			if atomic.CompareAndSwapInt64(&q.sendx, pos, pos+1) {
				break
			}
		} else if df < 0 {
			// full
			return false
		} else {
			pos = atomic.LoadInt64(&q.sendx)
		}
	}

	cell.value = v
	atomic.StoreInt64(&cell.sequence, pos+1)
	return true
}

// Recv receives a value from the queue and blocks when it is empty
func (q *MPMCqsDV[T]) Recv(v *T) bool {
	for wait := 0; ; spin(&wait) {
		if q.TryRecv(v) {
			return true
		}
	}
}

// TryRecv receives a value from the queue and returns when it is full
func (q *MPMCqsDV[T]) TryRecv(v *T) bool {
	var cell *seqValue[T]
	pos := atomic.LoadInt64(&q.recvx)
	for {
		cell = &q.buffer[pos&q.mask]
		seq := atomic.LoadInt64(&cell.sequence)
		df := seq - (pos + 1)
		if df == 0 {
			if atomic.CompareAndSwapInt64(&q.recvx, pos, pos+1) {
				break
			}
		} else if df < 0 {
			// empty
			return false
		} else {
			pos = atomic.LoadInt64(&q.recvx)
		}
	}

	*v = cell.value
	atomic.StoreInt64(&cell.sequence, pos+q.mask+1)
	return true
}
