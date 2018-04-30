package queue

import "testing"

var TestProcs = 16
var TestCount = [...]int{
	1, 2, 3,
	7, 8, 9,
	15, 16, 17,
	127, 128, 129,
	8191, 8192, 8193,
}

type Capability uint64

func (caps *Capability) Add(cap Capability)     { *caps |= cap }
func (caps Capability) Has(cap Capability) bool { return caps&cap == cap }

const (
	BlockSPSC = Capability(1 << iota)
	BlockMP
	BlockMC
	NonblockSPSC
	NonblockMP
	NonblockMC
	// BatchReceiver
)

func Detect(q Queue) Capability {
	var caps Capability
	if _, ok := q.(SPSC); ok {
		caps.Add(BlockSPSC)
	}
	if _, ok := q.(MPSC); ok {
		caps.Add(BlockMP)
	}
	if _, ok := q.(SPMC); ok {
		caps.Add(BlockMC)
	}
	if _, ok := q.(NonblockingSPSC); ok {
		caps.Add(NonblockSPSC)
	}
	if _, ok := q.(NonblockingMPSC); ok {
		caps.Add(NonblockMP)
	}
	if _, ok := q.(NonblockingSPMC); ok {
		caps.Add(NonblockMC)
	}
	return caps
}

// RunTest runs queue tests for queues that can be sized
func RunTest(t *testing.T, fn func() Queue) {
	caps := Detect(fn())
	if caps == 0 {
		t.Fatal("does not implement any of queue interfaces")
	}

}

// RunBenchmark runs queue benchmarks for queues that can be sized
func RunBenchmark(b *testing.B, fn func() Queue) {
	caps := Detect(fn())
	if caps == 0 {
		b.Fatal("does not implement any of queue interfaces")
	}

}
