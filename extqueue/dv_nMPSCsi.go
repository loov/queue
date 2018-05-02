package extqueue

import (
	"sync/atomic"
	"unsafe"
)

// MPSCnsiDV is a MPSC queue based on http://www.1024cores.net/home/lock-free-algorithms/queues/intrusive-mpsc-node-based-queue
type MPSCnsiDV struct {
	stub Node
	_    [7]uint64
	head unsafe.Pointer
	_    [7]uint64
	tail unsafe.Pointer
	_    [7]uint64
}

// NewMPSCnsiDV creates a MPSCnsDV queue
func NewMPSCnsiDV() *MPSCnsiDV {
	q := &MPSCnsiDV{}
	q.head = unsafe.Pointer(&q.stub)
	q.tail = unsafe.Pointer(&q.stub)
	return q
}

// MultipleProducers makes this a MP queue
func (q *MPSCnsiDV) MultipleProducers() {}

// Send sends a value to the queue, always suceeds
func (q *MPSCnsiDV) Send(value Value) bool { return q.SendNode(&Node{Value: value}) }

// TrySend sends a value to the queue, always suceeds
func (q *MPSCnsiDV) TrySend(value Value) bool { return q.SendNode(&Node{Value: value}) }

// SendNode sends a node to the queue, always suceeds
func (q *MPSCnsiDV) SendNode(node *Node) bool {
	node.next = nil
	prev := atomic.SwapPointer(&q.head, unsafe.Pointer(node))
	prevn := (*Node)(prev)
	atomic.StorePointer(&prevn.next, unsafe.Pointer(node))
	return true
}

// Recv receives a value from the queue and blocks when it is empty
func (q *MPSCnsiDV) Recv(value *Value) bool {
	node, ok := q.RecvNode()
	if ok {
		*value = node.Value
		return true
	}
	return false
}

// TryRecv receives a value from the queue and returns when it is empty
func (q *MPSCnsiDV) TryRecv(value *Value) bool {
	node, ok := q.TryRecvNode()
	if ok {
		*value = node.Value
		return true
	}
	return false
}

// RecvNode receives a node from the queue and blocks when it is empty
func (q *MPSCnsiDV) RecvNode() (*Node, bool) {
	for wait := 0; ; spin(&wait) {
		if node, ok := q.TryRecvNode(); ok {
			return node, true
		}
	}
}

// TryRecvNode receives a node from the queue and returns when it is empty
func (q *MPSCnsiDV) TryRecvNode() (*Node, bool) {
	tail := (*Node)(q.tail)
	next := atomic.LoadPointer(&tail.next)
	if tail == &q.stub {
		if next == nil {
			return nil, false
		}
		q.tail = next
		tail = (*Node)(next)
		next = atomic.LoadPointer(&tail.next)
	}
	if next != nil {
		q.tail = next
		tail.next = nil
		return tail, true
	}

	head := atomic.LoadPointer(&q.head)
	if q.tail != head {
		return nil, false
	}

	q.SendNode(&q.stub)
	next = atomic.LoadPointer(&tail.next)
	if next != nil {
		q.tail = next
		tail.next = nil
		return tail, true
	}

	return nil, false
}
