package extqueue

import (
	"sync"
	"sync/atomic"
)

// MPMCqpGo is an lock-free MPMC queue based on https://docs.google.com/document/d/1yIAYmbvL3JxOKOjuCyon7JhW4cSv1wy5hC0ApeGMV9s/pub
type MPMCqpGo[T any] struct {
	sendx  uint64
	_      [7]uint64
	recvx  uint64
	_      [7]uint64
	buffer []seqPaddedValue32[T]

	mu    sync.Mutex
	sendq sync.Cond
	recvq sync.Cond

	sendw, recvw int
}

// NewMPMCqpGo creates a new MPMCqpGo queue
func NewMPMCqpGo[T any](size int) *MPMCqpGo[T] {
	if size < 2 {
		size = 2
	}
	q := &MPMCqpGo[T]{
		sendx:  0,
		recvx:  0,
		buffer: make([]seqPaddedValue32[T], size),
	}
	q.sendq.L = &q.mu
	q.recvq.L = &q.mu
	return q
}

// Cap returns number of elements this queue can hold before blocking
func (q *MPMCqpGo[T]) Cap() int { return len(q.buffer) }

// MultipleConsumers makes this a MC queue
func (q *MPMCqpGo[T]) MultipleConsumers() {}

// MultipleProducers makes this a MP queue
func (q *MPMCqpGo[T]) MultipleProducers() {}

func (q *MPMCqpGo[T]) cap() uint32 { return uint32(len(q.buffer)) }

// Send sends a value to the queue and blocks when it is full
func (q *MPMCqpGo[T]) Send(value T) bool { return q.trySend(&value, true) }

// TrySend tries to send a value to the queue and returns immediately when it is full
func (q *MPMCqpGo[T]) TrySend(value T) bool { return q.trySend(&value, false) }

// Recv receives a value from the queue and blocks when it is empty
func (q *MPMCqpGo[T]) Recv(value *T) bool { return q.tryRecv(value, true) }

// TryRecv receives a value from the queue and returns when it is empty
func (q *MPMCqpGo[T]) TryRecv(value *T) bool { return q.tryRecv(value, false) }

func (q *MPMCqpGo[T]) trySend(value *T, block bool) bool {
	for loopCount := 0; ; backoff(&loopCount) {
		x := atomic.LoadUint64(&q.sendx)
		seq, pos := uint32(x>>32), uint32(x)
		elem := &q.buffer[pos]
		eseq := atomic.LoadUint32(&elem.sequence)
		//fmt.Printf("send: state %v %v %v\n", seq, pos, eseq)
		if seq == eseq {
			// The element is ready for writing on this seq.
			// Try to claim the right to write to this element.
			var newx uint64
			if pos+1 < q.cap() {
				newx = x + 1 // just increase the pos
			} else {
				newx = uint64(seq+2) << 32
			}

			if atomic.CompareAndSwapUint64(&q.sendx, x, newx) {
				// We own the element, do non-atomic write.
				elem.value = *value
				// Make the element available for reading.
				atomic.StoreUint32(&elem.sequence, eseq+1)

				// try to release a receiver
				q.mu.Lock()
				if q.recvw > 0 {
					q.recvq.Signal()
				}
				q.mu.Unlock()
				return true
			}
			// Lost the race, retry
		} else if int32(seq-eseq) > 0 {
			if !block {
				return false
			}

			if x-atomic.LoadUint64(&q.recvx) != 2<<32 {
				waitcount := 0
				//fmt.Printf("send: busy wait %v\n", pos)
				for int32(seq-atomic.LoadUint32(&elem.sequence)) > 0 {
					backoff(&waitcount)
				}
				continue
			}

			q.mu.Lock()
			if x-atomic.LoadUint64(&q.recvx) != 2<<32 {
				q.mu.Unlock()
				continue
			}
			q.sendw++
			//fmt.Printf("send: sleep %v\n", pos)
			q.sendq.Wait()
			q.sendw--
			q.mu.Unlock()
		}
		// The element has already been written on this seq,
		// this means that q.sendx has been changed as well,
		// retry.
	}
}

func (q *MPMCqpGo[T]) tryRecv(result *T, block bool) bool {
	var empty T
	for loopCount := 0; ; backoff(&loopCount) {
		// if closed return false

		x := atomic.LoadUint64(&q.recvx)
		seq, pos := uint32(x>>32), uint32(x)
		elem := &q.buffer[pos]
		eseq := atomic.LoadUint32(&elem.sequence) - 1
		//fmt.Printf("recv: state %v %v %v\n", seq, pos, eseq)
		if seq == eseq {
			// The element is ready for writing on this seq.
			// Try to claim the right to write to this element.
			var newx uint64
			if pos+1 < q.cap() {
				newx = x + 1 // just increase the pos
			} else {
				newx = uint64(seq+2) << 32
			}

			if atomic.CompareAndSwapUint64(&q.recvx, x, newx) {
				*result, elem.value = elem.value, empty
				atomic.StoreUint32(&elem.sequence, eseq+2)
				// try to release a sender
				q.mu.Lock()
				if q.sendw > 0 {
					q.sendq.Signal()
				}
				q.mu.Unlock()
				return true
			}
			// Lost the race, retry
		} else if int32(seq-eseq) > 0 {
			if !block {
				return false
			}

			if x != atomic.LoadUint64(&q.sendx) {
				waitcount := 0
				//fmt.Printf("recv: busy wait %v\n", pos)
				for int32(seq-atomic.LoadUint32(&elem.sequence)+1) > 0 {
					backoff(&waitcount)
				}
				continue
			}

			//fmt.Printf("recv: sleep %v\n", pos)
			q.mu.Lock()
			if x != atomic.LoadUint64(&q.sendx) {
				q.mu.Unlock()
				continue
			}
			q.recvw++
			q.recvq.Wait()
			q.recvw--
			q.mu.Unlock()
		}
		// The element has already been read on this seq,
		// this means that q.recvx has been changed as well,
		// retry.
	}
}
