package queue

var _ MPMC = (*MPMCc_go)(nil)

// MPMCc_go is a wrapper around go standard channel implementing Queue interfaces
type MPMCc_go struct {
	ch chan Value
}

// NewMPMCc_go creates a new MPMCc_go queue
func NewMPMCc_go(size int) *MPMCc_go {
	return &MPMCc_go{make(chan Value, size)}
}

// Cap returns number of elements this queue can hold before blocking
func (q *MPMCc_go) Cap() int { return cap(q.ch) }

// MultipleProducers makes this a MP queue
func (q *MPMCc_go) MultipleProducers() {}

// MultipleConsumers makes this a MC queue
func (q *MPMCc_go) MultipleConsumers() {}

// Close closes this queue
func (q *MPMCc_go) Close() { close(q.ch) }

// Send sends a value to the queue and blocks when it is full
func (q *MPMCc_go) Send(v Value) bool {
	select {
	case q.ch <- v:
		return true
	}
}

// Recv receives a value from the queue and blocks when it is empty
func (q *MPMCc_go) Recv(v *Value) bool {
	select {
	case x := <-q.ch:
		*v = x
		return true
	}
}

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *MPMCc_go) TrySend(v Value) bool {
	select {
	case q.ch <- v:
		return true
	default:
		return false
	}
}

// TryRecv receives a value from the queue and returns when it is empty
func (q *MPMCc_go) TryRecv(v *Value) bool {
	select {
	case x := <-q.ch:
		*v = x
		return true
	default:
		return false
	}
}
