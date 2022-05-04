package extqueue

// MPMCcGo is a wrapper around go standard channel implementing Queue interfaces
type MPMCcGo[T any] struct {
	ch chan T
}

// NewMPMCcGo creates a new MPMCcGo queue
func NewMPMCcGo[T any](size int) *MPMCcGo[T] {
	return &MPMCcGo[T]{make(chan T, size)}
}

// Cap returns number of elements this queue can hold before blocking
func (q *MPMCcGo[T]) Cap() int { return cap(q.ch) }

// MultipleProducers makes this a MP queue
func (q *MPMCcGo[T]) MultipleProducers() {}

// MultipleConsumers makes this a MC queue
func (q *MPMCcGo[T]) MultipleConsumers() {}

// Close closes this queue
func (q *MPMCcGo[T]) Close() { close(q.ch) }

// Send sends a value to the queue and blocks when it is full
func (q *MPMCcGo[T]) Send(v T) bool {
	q.ch <- v
	return true
}

// Recv receives a value from the queue and blocks when it is empty
func (q *MPMCcGo[T]) Recv(v *T) bool {
	*v = <-q.ch
	return true
}

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *MPMCcGo[T]) TrySend(v T) bool {
	select {
	case q.ch <- v:
		return true
	default:
		return false
	}
}

// TryRecv receives a value from the queue and returns when it is empty
func (q *MPMCcGo[T]) TryRecv(v *T) bool {
	select {
	case x := <-q.ch:
		*v = x
		return true
	default:
		return false
	}
}
