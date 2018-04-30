package queue

import (
	"flag"
	"testing"
)

var shake = flag.Int("shake", 1, "run tests multiple times")
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
	if !caps.Any(CapQueue) {
		t.Fatal("does not implement any of queue interfaces")
	}
	t.Helper()

	// TODO: add system noise when shaking

	if caps.Has(CapBlockSPSC) {
		for i := 0; i < *shake; i++ {
			t.Run("SPSC", func(t *testing.T) { t.Helper(); testSPSC(t, caps, ctor) })
		}
	}
	if caps.Has(CapBlockMPSC) {
		for i := 0; i < *shake; i++ {
			t.Run("MPSC", func(t *testing.T) { t.Helper(); testMPSC(t, caps, ctor) })
		}
	}
	if caps.Has(CapBlockSPMC) {
		for i := 0; i < *shake; i++ {
			t.Run("SPMC", func(t *testing.T) { t.Helper(); testSPMC(t, caps, ctor) })
		}
	}
	if caps.Has(CapBlockMPMC) {
		for i := 0; i < *shake; i++ {
			t.Run("MPMC", func(t *testing.T) { t.Helper(); testMPMC(t, caps, ctor) })
		}
	}
}

// RunBenchmarks runs queue benchmarks for queues
func RunBenchmarks(b *testing.B, ctor func() Queue) {
	caps := Detect(ctor())
	if !caps.Any(CapQueue) {
		b.Fatal("does not implement any of queue interfaces")
	}
	b.Helper()

	if caps.Has(CapBlockSPSC) {
		b.Run("SPSC", func(b *testing.B) { b.Helper(); benchSPSC(b, caps, ctor) })
	}
	if caps.Has(CapBlockMPSC) {
		b.Run("MPSC", func(b *testing.B) { b.Helper(); benchMPSC(b, caps, ctor) })
	}
	if caps.Has(CapBlockSPMC) {
		b.Run("SPMC", func(b *testing.B) { b.Helper(); benchSPMC(b, caps, ctor) })
	}
	if caps.Has(CapBlockMPMC) {
		b.Run("MPMC", func(b *testing.B) { b.Helper(); benchMPMC(b, caps, ctor) })
	}
}
