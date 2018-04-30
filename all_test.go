package queue

// Code generated by all_gen.go; DO NOT EDIT.

//go:generate go run all_gen.go -out all_test.go

import (
	"strconv"
	"testing"

	"github.com/loov/queue/testsuite"
)

var _ MPMC = (*MPMCc_go)(nil)
var _ NonblockingMPMC = (*MPMCc_go)(nil)

func TestMPMCc_go(t *testing.T) {
	batchSize := 0
	for _, size := range testsuite.TestSizes {
		name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		t.Run(name, func(t *testing.T) { testsuite.Tests(t, func() testsuite.Queue { return NewMPMCc_go(size) }) })
	}
}

func BenchmarkMPMCc_go(b *testing.B) {
	batchSize := 0
	for _, size := range testsuite.BenchSizes {
		name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		b.Run(name, func(b *testing.B) { testsuite.Benchmarks(b, func() testsuite.Queue { return NewMPMCc_go(size) }) })
	}
}

var _ MPMC = (*MPMCq_go)(nil)
var _ NonblockingMPMC = (*MPMCq_go)(nil)

func TestMPMCq_go(t *testing.T) {
	batchSize := 0
	for _, size := range testsuite.TestSizes {
		name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		t.Run(name, func(t *testing.T) { testsuite.Tests(t, func() testsuite.Queue { return NewMPMCq_go(size) }) })
	}
}

func BenchmarkMPMCq_go(b *testing.B) {
	batchSize := 0
	for _, size := range testsuite.BenchSizes {
		name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		b.Run(name, func(b *testing.B) { testsuite.Benchmarks(b, func() testsuite.Queue { return NewMPMCq_go(size) }) })
	}
}

var _ MPMC = (*MPMCqp_go)(nil)
var _ NonblockingMPMC = (*MPMCqp_go)(nil)

func TestMPMCqp_go(t *testing.T) {
	batchSize := 0
	for _, size := range testsuite.TestSizes {
		name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		t.Run(name, func(t *testing.T) { testsuite.Tests(t, func() testsuite.Queue { return NewMPMCqp_go(size) }) })
	}
}

func BenchmarkMPMCqp_go(b *testing.B) {
	batchSize := 0
	for _, size := range testsuite.BenchSizes {
		name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		b.Run(name, func(b *testing.B) { testsuite.Benchmarks(b, func() testsuite.Queue { return NewMPMCqp_go(size) }) })
	}
}

var _ SPSC = (*SPSCns_dv)(nil)
var _ NonblockingSPSC = (*SPSCns_dv)(nil)

func TestSPSCns_dv(t *testing.T) {
	batchSize := 0
	size := 0
	name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
	t.Run(name, func(t *testing.T) { testsuite.Tests(t, func() testsuite.Queue { return NewSPSCns_dv() }) })
}

func BenchmarkSPSCns_dv(b *testing.B) {
	batchSize := 0
	size := 0
	name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
	b.Run(name, func(b *testing.B) { testsuite.Benchmarks(b, func() testsuite.Queue { return NewSPSCns_dv() }) })
}

var _ MPSC = (*MPSCns_dv)(nil)
var _ NonblockingMPSC = (*MPSCns_dv)(nil)

func TestMPSCns_dv(t *testing.T) {
	batchSize := 0
	size := 0
	name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
	t.Run(name, func(t *testing.T) { testsuite.Tests(t, func() testsuite.Queue { return NewMPSCns_dv() }) })
}

func BenchmarkMPSCns_dv(b *testing.B) {
	batchSize := 0
	size := 0
	name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
	b.Run(name, func(b *testing.B) { testsuite.Benchmarks(b, func() testsuite.Queue { return NewMPSCns_dv() }) })
}

var _ MPMC = (*MPMCqs_dv)(nil)
var _ NonblockingMPMC = (*MPMCqs_dv)(nil)

func TestMPMCqs_dv(t *testing.T) {
	batchSize := 0
	for _, size := range testsuite.TestSizes {
		name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		t.Run(name, func(t *testing.T) { testsuite.Tests(t, func() testsuite.Queue { return NewMPMCqs_dv(size) }) })
	}
}

func BenchmarkMPMCqs_dv(b *testing.B) {
	batchSize := 0
	for _, size := range testsuite.BenchSizes {
		name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		b.Run(name, func(b *testing.B) { testsuite.Benchmarks(b, func() testsuite.Queue { return NewMPMCqs_dv(size) }) })
	}
}

var _ MPMC = (*MPMCqsp_dv)(nil)
var _ NonblockingMPMC = (*MPMCqsp_dv)(nil)

func TestMPMCqsp_dv(t *testing.T) {
	batchSize := 0
	for _, size := range testsuite.TestSizes {
		name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		t.Run(name, func(t *testing.T) { testsuite.Tests(t, func() testsuite.Queue { return NewMPMCqsp_dv(size) }) })
	}
}

func BenchmarkMPMCqsp_dv(b *testing.B) {
	batchSize := 0
	for _, size := range testsuite.BenchSizes {
		name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		b.Run(name, func(b *testing.B) { testsuite.Benchmarks(b, func() testsuite.Queue { return NewMPMCqsp_dv(size) }) })
	}
}

var _ SPMC = (*SPMCqs_dv)(nil)
var _ NonblockingSPMC = (*SPMCqs_dv)(nil)

func TestSPMCqs_dv(t *testing.T) {
	batchSize := 0
	for _, size := range testsuite.TestSizes {
		name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		t.Run(name, func(t *testing.T) { testsuite.Tests(t, func() testsuite.Queue { return NewSPMCqs_dv(size) }) })
	}
}

func BenchmarkSPMCqs_dv(b *testing.B) {
	batchSize := 0
	for _, size := range testsuite.BenchSizes {
		name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		b.Run(name, func(b *testing.B) { testsuite.Benchmarks(b, func() testsuite.Queue { return NewSPMCqs_dv(size) }) })
	}
}

var _ SPMC = (*SPMCqsp_dv)(nil)
var _ NonblockingSPMC = (*SPMCqsp_dv)(nil)

func TestSPMCqsp_dv(t *testing.T) {
	batchSize := 0
	for _, size := range testsuite.TestSizes {
		name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		t.Run(name, func(t *testing.T) { testsuite.Tests(t, func() testsuite.Queue { return NewSPMCqsp_dv(size) }) })
	}
}

func BenchmarkSPMCqsp_dv(b *testing.B) {
	batchSize := 0
	for _, size := range testsuite.BenchSizes {
		name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		b.Run(name, func(b *testing.B) { testsuite.Benchmarks(b, func() testsuite.Queue { return NewSPMCqsp_dv(size) }) })
	}
}

var _ SPSC = (*SPSCqs_dv)(nil)
var _ NonblockingSPSC = (*SPSCqs_dv)(nil)

func TestSPSCqs_dv(t *testing.T) {
	batchSize := 0
	for _, size := range testsuite.TestSizes {
		name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		t.Run(name, func(t *testing.T) { testsuite.Tests(t, func() testsuite.Queue { return NewSPSCqs_dv(size) }) })
	}
}

func BenchmarkSPSCqs_dv(b *testing.B) {
	batchSize := 0
	for _, size := range testsuite.BenchSizes {
		name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		b.Run(name, func(b *testing.B) { testsuite.Benchmarks(b, func() testsuite.Queue { return NewSPSCqs_dv(size) }) })
	}
}

var _ SPSC = (*SPSCqsp_dv)(nil)
var _ NonblockingSPSC = (*SPSCqsp_dv)(nil)

func TestSPSCqsp_dv(t *testing.T) {
	batchSize := 0
	for _, size := range testsuite.TestSizes {
		name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		t.Run(name, func(t *testing.T) { testsuite.Tests(t, func() testsuite.Queue { return NewSPSCqsp_dv(size) }) })
	}
}

func BenchmarkSPSCqsp_dv(b *testing.B) {
	batchSize := 0
	for _, size := range testsuite.BenchSizes {
		name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
		b.Run(name, func(b *testing.B) { testsuite.Benchmarks(b, func() testsuite.Queue { return NewSPSCqsp_dv(size) }) })
	}
}
