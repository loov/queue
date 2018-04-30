package extqueue

import (
	"sync/atomic"
	"unsafe"

	"github.com/egonelbre/exp/sync2/spin"
)

// MPSCnsDV is a MPSC queue based on http://www.1024cores.net/home/lock-free-algorithms/queues/non-intrusive-mpsc-node-based-queue
type MPSCnsDV struct {
	_    [8]uint64
	stub Node
	_    [7]uint64
	head unsafe.Pointer
	_    [7]uint64
	tail unsafe.Pointer
	_    [7]uint64
}

// NewMPSCnsDV creates a MPSCnsDV queue
func NewMPSCnsDV() *MPSCnsDV {
	q := &MPSCnsDV{}
	q.head = unsafe.Pointer(&q.stub)
	q.tail = unsafe.Pointer(&q.stub)
	return q
}

// MultipleProducers makes this a MP queue
func (q *MPSCnsDV) MultipleProducers() {}

// Send sends a value to the queue, always suceeds
func (q *MPSCnsDV) Send(value Value) bool {
	n := &Node{Value: value}
	prev := atomic.SwapPointer(&q.head, unsafe.Pointer(n))
	prevn := (*Node)(prev)
	atomic.StorePointer(&prevn.next, unsafe.Pointer(n))
	return true
}

// TrySend sends a value to the queue, always suceeds
func (q *MPSCnsDV) TrySend(value Value) bool { return q.Send(value) }

// Recv receives a value from the queue and blocks when it is empty
func (q *MPSCnsDV) Recv(value *Value) bool {
	var s spin.T256
	for s.Spin() {
		if q.TryRecv(value) {
			return true
		}
	}
	return false
}

// TryRecv receives a value from the queue and returns when it is empty
func (q *MPSCnsDV) TryRecv(value *Value) bool {
	tail := (*Node)(q.tail)
	next := atomic.LoadPointer(&tail.next)
	if next == nil {
		return false

	}
	q.tail = next
	*value = (*Node)(next).Value
	return true
}
