package queue

import (
	"strconv"
	"testing"

	. "github.com/loov/queue/testsuite"
	"github.com/pltr/onering"
)

// SPSC wrapper
var _ SPSC = (*SPSCrs_one)(nil)

type SPSCrs_one struct {
	onering.SPSC
	capacity int
}

func NewSPSCrs_one(size int) *SPSCrs_one {
	q := &SPSCrs_one{}
	q.Init(uint(size))
	q.capacity = size
	return q
}

func (q *SPSCrs_one) Cap() int { return q.capacity - 1 }

func (q *SPSCrs_one) Send(v Value) bool {
	q.Put(v)
	return true
}

func (q *SPSCrs_one) Recv(v *Value) bool {
	return q.Get(v)
}

func (q *SPSCrs_one) BatchRecv(fn func(Value)) bool {
	q.Consume(fn)
	return true
}

// SPMC wrapper

var _ SPMC = (*SPMCrs_one)(nil)

type SPMCrs_one struct {
	onering.SPMC
	capacity int
}

func NewSPMCrs_one(size int) *SPMCrs_one {
	q := &SPMCrs_one{}
	q.Init(uint(size))
	q.capacity = size
	return q
}

func (q *SPMCrs_one) MultipleConsumers() {}

func (q *SPMCrs_one) Cap() int { return q.capacity }

func (q *SPMCrs_one) Send(v Value) bool {
	q.Put(v)
	return true
}

func (q *SPMCrs_one) Recv(v *Value) bool {
	return q.Get(v)
}

// MPSC wrapper

var _ MPSC = (*MPSCrs_one)(nil)

type MPSCrs_one struct {
	onering.MPSC
	capacity int
}

func NewMPSCrs_one(size int) *MPSCrs_one {
	q := &MPSCrs_one{}
	q.Init(uint(size))
	q.capacity = size
	return q
}

func (q *MPSCrs_one) MultipleProducers() {}

func (q *MPSCrs_one) Cap() int { return q.capacity - 1 }

func (q *MPSCrs_one) Send(v Value) bool {
	q.Put(v)
	return true
}

func (q *MPSCrs_one) Recv(v *Value) bool {
	return q.Get(v)
}

func (q *MPSCrs_one) BatchRecv(fn func(Value)) bool {
	q.Consume(fn)
	return true
}

// suites

func supportedSize(size int) bool {
	if size <= 1 {
		return false
	}
	return size&(size-1) == 0
}

var _ SPSC = (*SPSCrs_one)(nil)

func TestSPSCrs_one(t *testing.T) {
	batchSize := 0
	for _, size := range TestSizes {
		if !supportedSize(size) {
			continue
		}
		name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		t.Run(name, func(t *testing.T) { Tests(t, func() Queue { return NewSPSCrs_one(size) }) })
	}
}

func BenchmarkSPSCrs_one(b *testing.B) {
	batchSize := 0
	for _, size := range BenchSizes {
		if !supportedSize(size) {
			continue
		}
		name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		b.Run(name, func(b *testing.B) { Benchmarks(b, func() Queue { return NewSPSCrs_one(size) }) })
	}
}

var _ SPMC = (*SPMCrs_one)(nil)

func TestSPMCrs_one(t *testing.T) {
	batchSize := 0
	for _, size := range TestSizes {
		if !supportedSize(size) {
			continue
		}
		name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		t.Run(name, func(t *testing.T) { Tests(t, func() Queue { return NewSPMCrs_one(size) }) })
	}
}

func BenchmarkSPMCrs_one(b *testing.B) {
	batchSize := 0
	for _, size := range BenchSizes {
		if !supportedSize(size) {
			continue
		}
		name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		b.Run(name, func(b *testing.B) { Benchmarks(b, func() Queue { return NewSPMCrs_one(size) }) })
	}
}

var _ MPSC = (*MPSCrs_one)(nil)

func TestMPSCrs_one(t *testing.T) {
	batchSize := 0
	for _, size := range TestSizes {
		if !supportedSize(size) {
			continue
		}
		name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		t.Run(name, func(t *testing.T) { Tests(t, func() Queue { return NewMPSCrs_one(size) }) })
	}
}

func BenchmarkMPSCrs_one(b *testing.B) {
	batchSize := 0
	for _, size := range BenchSizes {
		if !supportedSize(size) {
			continue
		}
		name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		b.Run(name, func(b *testing.B) { Benchmarks(b, func() Queue { return NewMPSCrs_one(size) }) })
	}
}
