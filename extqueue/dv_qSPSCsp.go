package extqueue

import (
	"sync/atomic"
)

// SPSCqspDV is a SPSC queue based on http://www.1024cores.net/home/lock-free-algorithms/queues/bounded-mpmc-queue
//
// This modifies the base algorithm by removing relevant atomic ops for Send and Recv
type SPSCqspDV struct {
	_ [8]int64

	buffer []seqPaddedValue
	mask   int64
	_      [4]int64

	sendx int64
	_     [7]int64

	recvx int64
	_     [7]int64
}

// NewSPSCqspDV creates a new SPSCqspDV queue
func NewSPSCqspDV(size int) *SPSCqspDV {
	if size <= 1 {
		size = 2
	}
	size = int(nextPowerOfTwo(uint32(size)))

	q := &SPSCqspDV{}
	q.buffer = make([]seqPaddedValue, size)
	q.mask = int64(size) - 1
	for i := range q.buffer {
		q.buffer[i].sequence = int64(i)
	}

	return q
}

// Cap returns number of elements this queue can hold before blocking
func (q *SPSCqspDV) Cap() int { return len(q.buffer) }

// Send sends a value to the queue and blocks when it is full
func (q *SPSCqspDV) Send(v Value) bool {
	for wait := 0; ; spin(&wait) {
		if q.TrySend(v) {
			return true
		}
	}
	return false
}

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *SPSCqspDV) TrySend(v Value) bool {
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
func (q *SPSCqspDV) Recv(v *Value) bool {
	for wait := 0; ; spin(&wait) {
		if q.TryRecv(v) {
			return true
		}
	}
	return false
}

// TryRecv receives a value from the queue and returns when it is empty
func (q *SPSCqspDV) TryRecv(v *Value) bool {
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
