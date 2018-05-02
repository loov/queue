// +build external

package extqueue

import (
	"github.com/loov/queue/testsuite"
	"github.com/pltr/onering"
)

func init() {
	All.Append(
		&Desc{"SPSCrsOR", Default, func(bs, s int) testsuite.Queue { return NewSPSCrsOR(s) }},
		&Desc{"SPMCrsOR", Default, func(bs, s int) testsuite.Queue { return NewSPMCrsOR(s) }},
		&Desc{"MPSCrsOR", Default, func(bs, s int) testsuite.Queue { return NewMPSCrsOR(s) }},
		&Desc{"MPMCrsOR", Default, func(bs, s int) testsuite.Queue { return NewMPMCrsOR(s) }},
	)
}

type SPSCrsOR struct {
	onering.SPSC
	capacity int
}

func NewSPSCrsOR(size int) *SPSCrsOR {
	q := &SPSCrsOR{}
	if size < 2 {
		size = 2
	}
	q.Init(uint32(size))
	q.capacity = size
	return q
}

func (q *SPSCrsOR) Cap() int { return q.capacity - 1 }

func (q *SPSCrsOR) Send(v Value) bool {
	q.Put(v)
	return true
}

func (q *SPSCrsOR) Recv(v *Value) bool {
	return q.Get(v)
}

// SPMC wrapper

type SPMCrsOR struct {
	onering.SPMC
	capacity int
}

func NewSPMCrsOR(size int) *SPMCrsOR {
	q := &SPMCrsOR{}
	if size < 2 {
		size = 2
	}
	q.Init(uint32(size))
	q.capacity = size
	return q
}

func (q *SPMCrsOR) MultipleConsumers() {}

func (q *SPMCrsOR) Cap() int { return q.capacity }

func (q *SPMCrsOR) Send(v Value) bool {
	q.Put(v)
	return true
}

func (q *SPMCrsOR) Recv(v *Value) bool {
	return q.Get(v)
}

// MPSC wrapper

type MPSCrsOR struct {
	onering.MPSC
	capacity int
}

func NewMPSCrsOR(size int) *MPSCrsOR {
	q := &MPSCrsOR{}
	if size < 2 {
		size = 2
	}
	q.Init(uint32(size))
	q.capacity = size
	return q
}

func (q *MPSCrsOR) MultipleProducers() {}

func (q *MPSCrsOR) Cap() int { return q.capacity - 1 }

func (q *MPSCrsOR) Send(v Value) bool {
	q.Put(v)
	return true
}

func (q *MPSCrsOR) Recv(v *Value) bool {
	return q.Get(v)
}

// MPMC wrapper

type MPMCrsOR struct {
	onering.MPMC
	capacity int
}

func NewMPMCrsOR(size int) *MPMCrsOR {
	q := &MPMCrsOR{}
	if size < 2 {
		size = 2
	}
	q.Init(uint32(size))
	q.capacity = size
	return q
}

func (q *MPMCrsOR) MultipleProducers() {}
func (q *MPMCrsOR) MultipleConsumers() {}

func (q *MPMCrsOR) Cap() int { return q.capacity - 1 }

func (q *MPMCrsOR) Send(v Value) bool {
	q.Put(v)
	return true
}

func (q *MPMCrsOR) Recv(v *Value) bool {
	return q.Get(v)
}

// suites

func supportedSize(size int) bool {
	if size <= 1 {
		return false
	}
	return size&(size-1) == 0
}
