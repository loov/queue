package extqueue

import (
	"sync"
)

// SPSCrMC is a SPSC queue based on MCRingBuffer http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.577.960&rep=rep1&type=pdf
type SPSCrMC struct {
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
	mu        sync.Mutex
	reader    sync.Cond
	writer    sync.Cond
	batchSize int64
	buffer    []Value
}

// NewSPSCrMC creates a new SPSCrMC queue
func NewSPSCrMC(batchSize, size int) *SPSCrMC {
	q := &SPSCrMC{}
	q.reader.L = &q.mu
	q.writer.L = &q.mu
	q.batchSize = int64(batchSize)
	q.buffer = make([]Value, ceil(size+1, batchSize))
	return q
}

// Cap returns number of elements this queue can hold before blocking
func (q *SPSCrMC) Cap() int { return len(q.buffer) - 1 }

func (q *SPSCrMC) next(i int64) int64 {
	r := i + 1
	if r >= int64(len(q.buffer)) {
		return 0
	}
	return r
}

// Send sends a value to the queue and blocks when it is full
func (q *SPSCrMC) Send(v Value) bool { return q.send(v, true) }

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *SPSCrMC) TrySend(v Value) bool { return q.send(v, false) }

// Recv receives a value from the queue and blocks when it is empty
func (q *SPSCrMC) Recv(v *Value) bool { return q.recv(v, true) }

// TryRecv receives a value from the queue and returns when it is empty
func (q *SPSCrMC) TryRecv(v *Value) bool { return q.recv(v, false) }

func (q *SPSCrMC) send(v Value, block bool) bool {
	afterNextWrite := q.next(q.nextWrite)
	if afterNextWrite == q.localRead {
		q.mu.Lock()
		if afterNextWrite == q.read {
			if !block {
				q.mu.Unlock()
				return false
			}
			q.writer.Wait()
		}
		q.localRead = q.read
		q.mu.Unlock()
	}

	q.buffer[q.nextWrite] = v
	q.nextWrite = afterNextWrite
	q.writeBatch++
	if q.writeBatch >= q.batchSize {
		q.FlushSend()
	}
	return true
}

func (q *SPSCrMC) FlushSend() {
	q.mu.Lock()
	q.write = q.nextWrite
	q.writeBatch = 0
	q.reader.Signal()
	q.mu.Unlock()
}

func (q *SPSCrMC) recv(v *Value, block bool) bool {
	if q.nextRead == q.localWrite {
		q.mu.Lock()
		if q.nextRead == q.write {
			if !block {
				q.mu.Unlock()
				return false
			}
			q.reader.Wait()
		}
		q.localWrite = q.write
		q.mu.Unlock()
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

func (q *SPSCrMC) FlushRecv() {
	q.mu.Lock()
	q.read = q.nextRead
	q.readBatch = 0
	q.writer.Signal()
	q.mu.Unlock()
}
