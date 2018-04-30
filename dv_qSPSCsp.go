package queue

import (
	"sync/atomic"

	"github.com/egonelbre/exp/sync2/spin"
)

var _ SPSC = (*SPSCqsp_dv)(nil)

// SPSCqsp_dv is a SPSC queue based on http://www.1024cores.net/home/lock-free-algorithms/queues/bounded-mpmc-queue
//
// This modifies the base algorithm by removing relevant atomic ops for Send and Recv
type SPSCqsp_dv struct {
	_ [8]int64

	buffer []seqPaddedValue
	mask   int64
	_      [4]int64

	sendx int64
	_     [7]int64

	recvx int64
	_     [7]int64
}

// NewSPSCqsp_dv creates a new SPSCqsp_dv queue
func NewSPSCqsp_dv(size int) *SPSCqsp_dv {
	if size <= 1 {
		size = 2
	}
	size = int(nextPowerOfTwo(uint32(size)))

	q := &SPSCqsp_dv{}
	q.buffer = make([]seqPaddedValue, size)
	q.mask = int64(size) - 1
	for i := range q.buffer {
		q.buffer[i].sequence = int64(i)
	}

	return q
}

// Cap returns number of elements this queue can hold before blocking
func (q *SPSCqsp_dv) Cap() int { return len(q.buffer) }

// Send sends a value to the queue and blocks when it is full
func (q *SPSCqsp_dv) Send(v Value) bool {
	var s spin.T256
	for s.Spin() {
		if q.TrySend(v) {
			return true
		}
	}
	return false
}

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *SPSCqsp_dv) TrySend(v Value) bool {
	var cell *seqPaddedValue
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
func (q *SPSCqsp_dv) Recv(v *Value) bool {
	var s spin.T256
	for s.Spin() {
		if q.TryRecv(v) {
			return true
		}
	}
	return false
}

// TryRecv receives a value from the queue and returns when it is empty
func (q *SPSCqsp_dv) TryRecv(v *Value) bool {
	var cell *seqPaddedValue
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
