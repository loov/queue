package extqueue

// MPMCcGo is a wrapper around go standard channel implementing Queue interfaces
type MPMCcGo struct {
	ch chan Value
}

// NewMPMCcGo creates a new MPMCcGo queue
func NewMPMCcGo(size int) *MPMCcGo {
	return &MPMCcGo{make(chan Value, size)}
}

// Cap returns number of elements this queue can hold before blocking
func (q *MPMCcGo) Cap() int { return cap(q.ch) }

// MultipleProducers makes this a MP queue
func (q *MPMCcGo) MultipleProducers() {}

// MultipleConsumers makes this a MC queue
func (q *MPMCcGo) MultipleConsumers() {}

// Close closes this queue
func (q *MPMCcGo) Close() { close(q.ch) }

// Send sends a value to the queue and blocks when it is full
func (q *MPMCcGo) Send(v Value) bool {
	q.ch <- v
	return true
}

// Recv receives a value from the queue and blocks when it is empty
func (q *MPMCcGo) Recv(v *Value) bool {
	*v = <-q.ch
	return true
}

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *MPMCcGo) TrySend(v Value) bool {
	select {
	case q.ch <- v:
		return true
	default:
		return false
	}
}

// TryRecv receives a value from the queue and returns when it is empty
func (q *MPMCcGo) TryRecv(v *Value) bool {
	select {
	case x := <-q.ch:
		*v = x
		return true
	default:
		return false
	}
}
