package queue

import "testing"

var TestProcs = 16
var TestSizes = [...]int{
	1, 2, 3,
	7, 8, 9,
	15, 16, 17,
	127, 128, 129,
	8191, 8192, 8193,
}
var TestCount = [...]int{
	1, 2, 3,
	7, 8, 9,
	15, 16, 17,
	127, 128, 129,
	8191, 8192, 8193,
}

type Capability uint64

func (caps *Capability) Add(cap Capability)     { *caps |= cap }
func (caps Capability) Any(cap Capability) bool { return caps&cap != 0 }
func (caps Capability) Has(cap Capability) bool { return caps&cap == cap }

const (
	IsBlockSPSC = Capability(1 << iota)
	IsBlockMPSC
	IsBlockSPMC
	IsNonblockSPSC
	IsNonblockMPSC
	IsNonblockSPMC
	IsBounded
	// BatchReceiver

	IsBlockMPMC    = IsBlockMPSC | IsBlockSPMC
	IsNonblockMPMC = IsNonblockMPSC | IsNonblockSPMC

	AnyQueue = IsBlockMPMC | IsNonblockMPMC
)

func Detect(q Queue) Capability {
	var caps Capability
	if _, ok := q.(SPSC); ok {
		caps.Add(IsBlockSPSC)
	}
	if _, ok := q.(MPSC); ok {
		caps.Add(IsBlockMPSC)
	}
	if _, ok := q.(SPMC); ok {
		caps.Add(IsBlockSPMC)
	}
	if _, ok := q.(NonblockingSPSC); ok {
		caps.Add(IsNonblockSPSC)
	}
	if _, ok := q.(NonblockingMPSC); ok {
		caps.Add(IsNonblockMPSC)
	}
	if _, ok := q.(NonblockingSPMC); ok {
		caps.Add(IsNonblockSPMC)
	}
	if _, ok := q.(Bounded); ok {
		caps.Add(IsBounded)
	}
	return caps
}

// RunTest runs queue tests for queues
func RunTest(t *testing.T, ctor func() Queue) {
	caps := Detect(ctor())
	if !caps.Any(AnyQueue) {
		t.Fatal("does not implement any of queue interfaces")
	}
	t.Helper()

	if caps.Has(IsBlockSPSC) {
		t.Run("SPSC", func(t *testing.T) { t.Helper(); testSPSC(t, ctor) })
	}
	if caps.Has(IsBlockMPSC) {
		t.Run("MPSC", func(t *testing.T) { t.Helper(); testMPSC(t, ctor) })
	}
	if caps.Has(IsBlockSPMC) {
		t.Run("SPMC", func(t *testing.T) { t.Helper(); testSPMC(t, ctor) })
	}
	if caps.Has(IsBlockMPMC) {
		t.Run("MPMC", func(t *testing.T) { t.Helper(); testMPMC(t, ctor) })
	}
}

// RunBenchmark runs queue benchmarks for queues
func RunBenchmark(b *testing.B, ctor func() Queue) {
	caps := Detect(ctor())
	if !caps.Any(AnyQueue) {
		b.Fatal("does not implement any of queue interfaces")
	}
	b.Helper()
}

func testSPSC(t *testing.T, ctor func() Queue) {}
func testMPSC(t *testing.T, ctor func() Queue) {}
func testSPMC(t *testing.T, ctor func() Queue) {}
func testMPMC(t *testing.T, ctor func() Queue) {}
