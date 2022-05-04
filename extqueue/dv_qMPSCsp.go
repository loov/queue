package extqueue

import (
	"sync/atomic"
)

// MPSCqspDV is a MPMC queue based on http://www.1024cores.net/home/lock-free-algorithms/queues/bounded-mpmc-queue
type MPSCqspDV[T any] struct {
	_      [8]int64
	mask   int64
	buffer []seqPaddedValue[T]
	_      [4]int64
	sendx  int64
	_      [7]int64
	recvx  int64
	_      [7]int64
}

// NewMPSCqspDV creates a new queue.
func NewMPSCqspDV[T any](size int) *MPSCqspDV[T] {
	if size <= 1 {
		size = 2
	}
	size = int(nextPowerOfTwo(uint32(size)))

	q := &MPSCqspDV[T]{}
	q.buffer = make([]seqPaddedValue[T], size)
	q.mask = int64(size) - 1
	for i := range q.buffer {
		q.buffer[i].sequence = int64(i)
	}

	return q
}

// Cap returns number of elements this queue can hold before blocking
func (q *MPSCqspDV[T]) Cap() int { return len(q.buffer) }

// MultipleProducers makes this a MP queue
func (q *MPSCqspDV[T]) MultipleProducers() {}

// Send sends a value to the queue and blocks when it is full
func (q *MPSCqspDV[T]) Send(v T) bool {
	for wait := 0; ; spin(&wait) {
		if q.TrySend(v) {
			return true
		}
	}
}

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *MPSCqspDV[T]) TrySend(v T) bool {
	var cell *seqPaddedValue[T]
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
func (q *MPSCqspDV[T]) Recv(v *T) bool {
	for wait := 0; ; spin(&wait) {
		if q.TryRecv(v) {
			return true
		}
	}
}

// TryRecv receives a value from the queue and returns when it is empty
func (q *MPSCqspDV[T]) TryRecv(v *T) bool {
	var cell *seqPaddedValue[T]
	pos := q.recvx
	for {
		cell = &q.buffer[pos&q.mask]
		seq := atomic.LoadInt64(&cell.sequence)
		df := seq - (pos + 1)
		if df == 0 {
			q.recvx = pos + 1
			break
		} else if df < 0 {
			// empty
			return false
		}
	}

	*v = cell.value
	atomic.StoreInt64(&cell.sequence, pos+q.mask+1)
	return true
}
