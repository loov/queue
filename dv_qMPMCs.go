package queue

import (
	"sync/atomic"

	"github.com/egonelbre/exp/sync2/spin"
)

var _ MPMC = (*MPMCqs_dv)(nil)

// MPMCqs_dv is a MPMC queue based on http://www.1024cores.net/home/lock-free-algorithms/queues/bounded-mpmc-queue
type MPMCqs_dv struct {
	_ [8]int64

	buffer []seqValue
	mask   int64
	_      [4]int64

	sendx int64
	_     [7]int64

	recvx int64
	_     [7]int64
}

// NewMPMCqs_dv creates a NewMPMCqs_dv queue
func NewMPMCqs_dv(size int) *MPMCqs_dv {
	if size <= 1 {
		size = 2
	}
	size = int(nextPowerOfTwo(uint32(size)))

	q := &MPMCqs_dv{}
	q.buffer = make([]seqValue, size)
	q.mask = int64(size) - 1
	for i := range q.buffer {
		q.buffer[i].sequence = int64(i)
	}

	return q
}

// Cap returns number of elements this queue can hold before blocking
func (q *MPMCqs_dv) Cap() int { return len(q.buffer) }

// MultipleConsumers makes this a MC queue
func (q *MPMCqs_dv) MultipleConsumers() {}

// MultipleProducers makes this a MP queue
func (q *MPMCqs_dv) MultipleProducers() {}

// Send sends a value to the queue and blocks when it is full
func (q *MPMCqs_dv) Send(v Value) bool {
	var s spin.T256
	for s.Spin() {
		if q.TrySend(v) {
			return true
		}
	}
	return false
}

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *MPMCqs_dv) TrySend(v Value) bool {
	var cell *seqValue
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
func (q *MPMCqs_dv) Recv(v *Value) bool {
	var s spin.T256
	for s.Spin() {
		if q.TryRecv(v) {
			return true
		}
	}
	return false
}

// TryRecv receives a value from the queue and returns when it is full
func (q *MPMCqs_dv) TryRecv(v *Value) bool {
	var cell *seqValue
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
