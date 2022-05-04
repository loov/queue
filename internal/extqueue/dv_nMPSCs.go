package extqueue

import (
	"sync/atomic"
	"unsafe"
)

// MPSCnsDV is a MPSC queue based on http://www.1024cores.net/home/lock-free-algorithms/queues/non-intrusive-mpsc-node-based-queue
type MPSCnsDV[T any] struct {
	stub Node[T]
	_    [7]uint64
	head unsafe.Pointer
	_    [7]uint64
	tail unsafe.Pointer
	_    [7]uint64
}

// NewMPSCnsDV creates a MPSCnsDV queue
func NewMPSCnsDV[T any]() *MPSCnsDV[T] {
	q := &MPSCnsDV[T]{}
	q.head = unsafe.Pointer(&q.stub)
	q.tail = unsafe.Pointer(&q.stub)
	return q
}

// MultipleProducers makes this a MP queue
func (q *MPSCnsDV[T]) MultipleProducers() {}

// Send sends a value to the queue, always suceeds
func (q *MPSCnsDV[T]) Send(value T) bool {
	n := &Node[T]{Value: value}
	prev := atomic.SwapPointer(&q.head, unsafe.Pointer(n))
	prevn := (*Node[T])(prev)
	atomic.StorePointer(&prevn.next, unsafe.Pointer(n))
	return true
}

// TrySend sends a value to the queue, always suceeds
func (q *MPSCnsDV[T]) TrySend(value T) bool { return q.Send(value) }

// Recv receives a value from the queue and blocks when it is empty
func (q *MPSCnsDV[T]) Recv(value *T) bool {
	for wait := 0; ; spin(&wait) {
		if q.TryRecv(value) {
			return true
		}
	}
}

// TryRecv receives a value from the queue and returns when it is empty
func (q *MPSCnsDV[T]) TryRecv(value *T) bool {
	tail := (*Node[T])(q.tail)
	next := atomic.LoadPointer(&tail.next)
	if next == nil {
		return false

	}
	q.tail = next
	*value = (*Node[T])(next).Value
	return true
}
