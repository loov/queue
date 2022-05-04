package extqueue

import (
	"sync/atomic"
	"unsafe"
)

// SPSCnsDV is a SPSC queue based on http://www.1024cores.net/home/lock-free-algorithms/queues/unbounded-spsc-queue
type SPSCnsDV[T any] struct {
	stub Node[T]
	_    [7]uint64
	// producer
	head     unsafe.Pointer
	first    unsafe.Pointer
	tailCopy unsafe.Pointer
	_        [7]uint64
	// consumer
	tail unsafe.Pointer
	_    [7]uint64
}

// NewSPSCnsDV creates a new SPSCnsDV queue
func NewSPSCnsDV[T any]() *SPSCnsDV[T] {
	q := &SPSCnsDV[T]{}
	q.head = unsafe.Pointer(&q.stub)
	q.tail = unsafe.Pointer(&q.stub)
	q.first = unsafe.Pointer(&q.stub)
	q.tailCopy = unsafe.Pointer(&q.stub)
	return q
}

// Send sends a value to the queue, always succeeds
func (q *SPSCnsDV[T]) Send(value T) bool {
	n := q.alloc()
	n.Value = value
	n.next = nil
	atomic.StorePointer(&(*Node[T])(q.head).next, unsafe.Pointer(n))
	q.head = unsafe.Pointer(n)
	return true
}

// TrySend tries to send a value to the queue, always succeeds
func (q *SPSCnsDV[T]) TrySend(value T) bool { return q.Send(value) }

// Recv receives a value from the queue and blocks when it is empty
func (q *SPSCnsDV[T]) Recv(value *T) bool {
	for wait := 0; ; spin(&wait) {
		if q.TryRecv(value) {
			return true
		}
	}
}

// TryRecv receives a value from the queue and returns when it is full
func (q *SPSCnsDV[T]) TryRecv(value *T) bool {
	tail := (*Node[T])(q.tail)
	next := atomic.LoadPointer(&tail.next)
	if next == nil {
		return false
	}
	q.tail = next
	*value = (*Node[T])(next).Value
	return true
}

func (q *SPSCnsDV[T]) alloc() *Node[T] {
	// first tries to allocate node from internal node cache,
	// if attempt fails, allocates node via ::operator new()

	if q.first != q.tailCopy {
		n := (*Node[T])(q.first)
		q.first = n.next
		return n
	}

	q.tailCopy = atomic.LoadPointer(&q.tail)
	if q.first != q.tailCopy {
		n := (*Node[T])(q.first)
		q.first = n.next
		return n
	}

	return &Node[T]{}
}
