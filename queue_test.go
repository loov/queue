package queue

import (
	"fmt"
	"testing"
)

var TestProcs = 16

var TestSizes = [...]int{
	1, 2, 3,
	7, 8, 9,
	15, 16, 17,
	127, 128, 129,
	8191, 8192, 8193,
}
var BenchSizes = [...]int{8, 64, 8192}

var TestCount = [...]int{
	1, 2, 3,
	7, 8, 9,
	15, 16, 17,
	127, 128, 129,
	8191, 8192, 8193,
}

// RunTests runs queue tests for queues
func RunTests(t *testing.T, ctor func() Queue) {
	q := ctor()
	caps := Detect(q)
	fmt.Printf("%T: %v\n", q, caps)
	if !caps.Any(CapQueue) {
		t.Fatal("does not implement any of queue interfaces")
	}
	t.Helper()

	if caps.Has(CapBlockSPSC) {
		t.Run("SPSC", func(t *testing.T) { t.Helper(); testSPSC(t, ctor) })
	}
	if caps.Has(CapBlockMPSC) {
		t.Run("MPSC", func(t *testing.T) { t.Helper(); testMPSC(t, ctor) })
	}
	if caps.Has(CapBlockSPMC) {
		t.Run("SPMC", func(t *testing.T) { t.Helper(); testSPMC(t, ctor) })
	}
	if caps.Has(CapBlockMPMC) {
		t.Run("MPMC", func(t *testing.T) { t.Helper(); testMPMC(t, ctor) })
	}
}

// RunBenchmarks runs queue benchmarks for queues
func RunBenchmarks(b *testing.B, ctor func() Queue) {
	caps := Detect(ctor())
	if !caps.Any(CapQueue) {
		b.Fatal("does not implement any of queue interfaces")
	}
	b.Helper()
}

func testSPSC(t *testing.T, ctor func() Queue) {}
func testMPSC(t *testing.T, ctor func() Queue) {}
func testSPMC(t *testing.T, ctor func() Queue) {}
func testMPMC(t *testing.T, ctor func() Queue) {}
