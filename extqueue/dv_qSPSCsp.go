package extqueue

import (
	"sync/atomic"
)

// SPSCqspDV is a SPSC queue based on http://www.1024cores.net/home/lock-free-algorithms/queues/bounded-mpmc-queue
// The base algorithm is modified by removing some atomic operations on the consumer and producer side.
type SPSCqspDV[T any] struct {
	_      [8]int64
	sendx  int64
	_      [7]int64
	recvx  int64
	_      [7]int64
	mask   int64
	buffer []seqPaddedValue[T]
}

// NewSPSCqspDV creates a new SPSCqspDV queue
func NewSPSCqspDV[T any](size int) *SPSCqspDV[T] {
	if size <= 1 {
		size = 2
	}
	size = int(nextPowerOfTwo(uint32(size)))

	q := &SPSCqspDV[T]{}
	q.buffer = make([]seqPaddedValue[T], size)
	q.mask = int64(size) - 1
	for i := range q.buffer {
		q.buffer[i].sequence = int64(i)
	}

	return q
}

// Cap returns number of elements this queue can hold before blocking
func (q *SPSCqspDV[T]) Cap() int { return len(q.buffer) }

// Send sends a value to the queue and blocks when it is full
func (q *SPSCqspDV[T]) Send(v T) bool {
	for wait := 0; ; spin(&wait) {
		if q.TrySend(v) {
			return true
		}
	}
}

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *SPSCqspDV[T]) TrySend(v T) bool {
	var cell *seqPaddedValue[T]
	pos := q.sendx
	for {
		cell = &q.buffer[pos&q.mask]
		seq := atomic.LoadInt64(&cell.sequence)
		df := seq - pos
		if df == 0 {
			q.sendx = pos + 1
			break
		} else if df < 0 {
			// full
			return false
		}
	}

	cell.value = v
	atomic.StoreInt64(&cell.sequence, pos+1)
	return true
}

// Recv receives a value from the queue and blocks when it is empty
func (q *SPSCqspDV[T]) Recv(v *T) bool {
	for wait := 0; ; spin(&wait) {
		if q.TryRecv(v) {
			return true
		}
	}
}

// TryRecv receives a value from the queue and returns when it is empty
func (q *SPSCqspDV[T]) TryRecv(v *T) bool {
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
