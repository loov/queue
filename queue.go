package queue

// Value is the data being sent
type Value = int64

// Queue is the general interface
type Queue interface {
	// SPSC and/or NonblockingSPSC
	// Closer?
	// Flusher?
	// Capped?
}

// Bounded returns number of elements that can be added until the queue
// either Send blocks or TrySend fails.
type Bounded interface {
	Cap() int
}

// Flusher implements API for queues that need explicit flushing for the sending
// or receiving to be propagated to the other end.
type Flusher interface {
	// FlushSend propagates pending send operations to the receiver.
	FlushSend()
	// FlushRecv propagates pending receive operations to the sender.
	FlushRecv()
}

// Closer is implemented by queues that can be closed
type Closer interface {
	// Close the queue for further sending or receiving, depening on the
	// implementation
	Close()
}

// TODO:
// type BatchReceiver interface {
// 	// BatchRecv receives a batch
// 	BatchRecv(func(v Value)) bool
// }

// SPSC is a blocking single-producer and single-consumer queue,
// which waits until Send or Recv succeeds
type SPSC interface {
	Queue
	// Send puts a value to a queue,
	// returns false when the queue has been closed
	Send(v Value) bool
	// Recv takes a value from the queue
	// returns false when the queue has been closed
	Recv(v *Value) bool
}

// MPSC is a blocking multi-producer and single-consumer queue,
// which waits until Send or Recv succeeds
type MPSC interface {
	SPSC
	MultipleProducers()
}

// SPMC is a blocking single-producer and multi-consumer queue,
// which waits until Send or Recv succeeds
type SPMC interface {
	SPSC
	MultipleConsumers()
}

// MPMC is a blocking multi-producer and multi-consumer queue,
// which waits until Send or Recv succeeds
type MPMC interface {
	SPSC
	MultipleProducers()
	MultipleConsumers()
}

// NonblockingSPSC is a non-blocking single-producer and single-consumer queue, which
// which returns in case Send or Recv cannot be completed
type NonblockingSPSC interface {
	Queue
	// TrySend tries to put a value to a queue
	// returns false when the queue if full or closed
	TrySend(v Value) bool
	// TrySend tries to take a value to a queue
	// returns false when the queue if empty or closed
	TryRecv(v *Value) bool
}

// NonblockingMPSC is a non-blocking multi-producer and single-consumer queue, which
// which returns in case Send or Recv cannot be completed
type NonblockingMPSC interface {
	NonblockingSPSC
	MultipleProducers()
}

// NonblockingSPMC is a non-blocking single-producer and multi-consumer queue, which
// which returns in case Send or Recv cannot be completed
type NonblockingSPMC interface {
	NonblockingSPSC
	MultipleConsumers()
}

// NonblockingMPMC is a non-blocking multi-producer and multi-consumer queue, which
// which returns in case Send or Recv cannot be completed
type NonblockingMPMC interface {
	NonblockingSPSC
	MultipleProducers()
	MultipleConsumers()
}
