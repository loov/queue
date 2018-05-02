package extqueue

import (
	"sync/atomic"
)

// SPSCqsDV is a SPSC queue based on http://www.1024cores.net/home/lock-free-algorithms/queues/bounded-mpmc-queue
// The base algorithm is modified by removing some atomic operations on the consumer and producer side.
type SPSCqsDV struct {
	_ [8]int64

	buffer []seqValue
	mask   int64
	_      [4]int64

	sendx int64
	_     [7]int64

	recvx int64
	_     [7]int64
}

// NewSPSCqsDV creates a new SPSCqsDV queue
func NewSPSCqsDV(size int) *SPSCqsDV {
	if size <= 1 {
		size = 2
	}
	size = int(nextPowerOfTwo(uint32(size)))

	q := &SPSCqsDV{}
	q.buffer = make([]seqValue, size)
	q.mask = int64(size) - 1
	for i := range q.buffer {
		q.buffer[i].sequence = int64(i)
	}

	return q
}

// Cap returns number of elements this queue can hold before blocking
func (q *SPSCqsDV) Cap() int { return len(q.buffer) }

// Send sends a value to the queue and blocks when it is full
func (q *SPSCqsDV) Send(v Value) bool {
	for wait := 0; ; spin(&wait) {
		if q.TrySend(v) {
			return true
		}
	}
}

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *SPSCqsDV) TrySend(v Value) bool {
	var cell *seqValue
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
func (q *SPSCqsDV) Recv(v *Value) bool {
	for wait := 0; ; spin(&wait) {
		if q.TryRecv(v) {
			return true
		}
	}
}

// TryRecv receives a value from the queue and returns when it is empty
func (q *SPSCqsDV) TryRecv(v *Value) bool {
	var cell *seqValue
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
