package extqueue

import (
	"sync/atomic"
)

// MPSCqsDV is a MPMC queue based on http://www.1024cores.net/home/lock-free-algorithms/queues/bounded-mpmc-queue
type MPSCqsDV struct {
	buffer []seqValue
	mask   int64
	_      [4]int64
	sendx  int64
	_      [7]int64
	recvx  int64
	_      [7]int64
}

// NewMPSCqsDV creates a NewMPSCqsDV queue
func NewMPSCqsDV(size int) *MPSCqsDV {
	if size <= 1 {
		size = 2
	}
	size = int(nextPowerOfTwo(uint32(size)))

	q := &MPSCqsDV{}
	q.buffer = make([]seqValue, size)
	q.mask = int64(size) - 1
	for i := range q.buffer {
		q.buffer[i].sequence = int64(i)
	}

	return q
}

// Cap returns number of elements this queue can hold before blocking
func (q *MPSCqsDV) Cap() int { return len(q.buffer) }

// MultipleProducers makes this a MP queue
func (q *MPSCqsDV) MultipleProducers() {}

// Send sends a value to the queue and blocks when it is full
func (q *MPSCqsDV) Send(v Value) bool {
	for wait := 0; ; spin(&wait) {
		if q.TrySend(v) {
			return true
		}
	}
}

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *MPSCqsDV) TrySend(v Value) bool {
	for {
		pos := atomic.LoadInt64(&q.sendx)
		cell := &q.buffer[pos&q.mask]
		seq := atomic.LoadInt64(&cell.sequence)
		if seq-pos == 0 {
			if atomic.CompareAndSwapInt64(&q.sendx, pos, pos+1) {
				cell.value = v
				atomic.StoreInt64(&cell.sequence, pos+1)
				return true
			}
		} else if seq-pos < 0 {
			// full
			return false
		}
	}
}

// Recv receives a value from the queue and blocks when it is empty
func (q *MPSCqsDV) Recv(v *Value) bool {
	for wait := 0; ; spin(&wait) {
		if q.TryRecv(v) {
			return true
		}
	}
}

// TryRecv receives a value from the queue and returns when it is empty
func (q *MPSCqsDV) TryRecv(v *Value) bool {
	pos := q.recvx
	cell := &q.buffer[pos&q.mask]
	for {
		seq := atomic.LoadInt64(&cell.sequence)
		df := seq - (pos + 1)
		if df == 0 {
			*v = cell.value
			atomic.StoreInt64(&cell.sequence, pos+q.mask+1)
			q.recvx = pos + 1
			return true
		} else if df < 0 {
			// empty
			return false
		}
	}
}
