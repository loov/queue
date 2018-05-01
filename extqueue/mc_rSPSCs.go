package extqueue

import (
	"sync/atomic"
)

// SPSCrsMC is a SPSC queue based on MCRingBuffer http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.577.960&rep=rep1&type=pdf
type SPSCrsMC struct {
	_ [8]uint64
	// volatile
	read  int64
	write int64
	_     [8 - 2]uint64
	// consumer
	localWrite int64
	nextRead   int64
	readBatch  int64
	_          [8 - 3]uint64
	// producer
	localRead  int64
	nextWrite  int64
	writeBatch int64
	_          [8 - 3]uint64
	// constant
	batchSize int64
	buffer    []Value
}

// NewSPSCrsMC creates a new SPSCrsMC queue
func NewSPSCrsMC(batchSize, size int) *SPSCrsMC {
	q := &SPSCrsMC{}
	q.batchSize = int64(batchSize)
	q.buffer = make([]Value, ceil(size+1, batchSize))
	return q
}

// Cap returns number of elements this queue can hold before blocking
func (q *SPSCrsMC) Cap() int { return len(q.buffer) - 1 }

func (q *SPSCrsMC) next(i int64) int64 {
	r := i + 1
	if r >= int64(len(q.buffer)) {
		return 0
	}
	return r
}

// Send sends a value to the queue and blocks when it is full
func (q *SPSCrsMC) Send(v Value) bool { return q.send(v, true) }

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *SPSCrsMC) TrySend(v Value) bool { return q.send(v, false) }

// Recv receives a value from the queue and blocks when it is empty
func (q *SPSCrsMC) Recv(v *Value) bool { return q.recv(v, true) }

// TryRecv receives a value from the queue and returns when it is empty
func (q *SPSCrsMC) TryRecv(v *Value) bool { return q.recv(v, false) }

func (q *SPSCrsMC) send(v Value, block bool) bool {
	afterNextWrite := q.next(q.nextWrite)
	if afterNextWrite == q.localRead {
		for try := 0; afterNextWrite == atomic.LoadInt64(&q.read); spin(&try) {
			if !block {
				return false
			}
		}
		q.localRead = atomic.LoadInt64(&q.read)
	}

	q.buffer[q.nextWrite] = v
	q.nextWrite = afterNextWrite
	q.writeBatch++
	if q.writeBatch >= q.batchSize {
		q.FlushSend()
	}
	return true
}

// FlushSend propagates pending send operations to the receiver.
func (q *SPSCrsMC) FlushSend() {
	atomic.StoreInt64(&q.write, q.nextWrite)
	q.writeBatch = 0
}

func (q *SPSCrsMC) recv(v *Value, block bool) bool {
	if q.nextRead == q.localWrite {
		for try := 0; q.nextRead == atomic.LoadInt64(&q.write); spin(&try) {
			if !block {
				return false
			}
		}
		q.localWrite = atomic.LoadInt64(&q.write)
	}

	*v = q.buffer[q.nextRead]
	// q.buffer[q.nextRead] = 0 clear value, only needed for pointers

	q.nextRead = q.next(q.nextRead)
	q.readBatch++
	if q.readBatch >= q.batchSize {
		q.FlushRecv()
	}

	return true
}

// FlushRecv propagates pending receive operations to the sender.
func (q *SPSCrsMC) FlushRecv() {
	atomic.StoreInt64(&q.read, q.nextRead)
	q.readBatch = 0
}
