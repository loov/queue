package testsuite

import (
	"flag"
	"testing"
	"time"
)

var NonblockThreshold = time.Second / 8

var TestProcs = 16

var BatchSizes = []int{1, 4, 8, 16}
var TestSizes = []int{
	1, 2, 3,
	7, 8, 9,
	127, 128, 129,
	1023, 1024, 1025,
}

var BenchBatchSizes = []int{32, 256}
var BenchSizes = []int{256, 8192}
var BenchWork = []int{0}

var TestCount = []int{
	1, 2, 3,
	7, 8, 9,
	127, 128, 129,
	1023, 1024, 1025,
}

var shake = flag.Int("shake", 1, "run tests multiple times")

func skipRedundant(q Queue, testsize int) bool {
	if b, ok := q.(Bounded); ok {
		if b.Cap()*10 < testsize {
			// fmt.Println("skip", b.Cap(), testsize)
			return true
		}
	}
	return false
}

// Tests runs queue tests for queues
func Tests(t *testing.T, ctor func() Queue) {
	q := ctor()
	caps := Detect(q)
	if !caps.Any(CapQueue) {
		t.Fatal("does not implement any of queue interfaces")
	}
	t.Helper()

	// TODO: add system noise when shaking

	if caps.Has(CapBlockSPSC) {
		for i := 0; i < *shake; i++ {
			t.Run("b/SPSC", func(t *testing.T) { t.Helper(); testSPSC(t, caps, ctor) })
		}
	}
	if caps.Has(CapBlockMPSC) {
		for i := 0; i < *shake; i++ {
			t.Run("b/MPSC", func(t *testing.T) { t.Helper(); testMPSC(t, caps, ctor) })
		}
	}
	if caps.Has(CapBlockSPMC) {
		for i := 0; i < *shake; i++ {
			t.Run("b/SPMC", func(t *testing.T) { t.Helper(); testSPMC(t, caps, ctor) })
		}
	}
	if caps.Has(CapBlockMPMC) {
		for i := 0; i < *shake; i++ {
			t.Run("b/MPMC", func(t *testing.T) { t.Helper(); testMPMC(t, caps, ctor) })
		}
	}

	if caps.Has(CapNonblockSPSC) {
		for i := 0; i < *shake; i++ {
			t.Run("n/SPSC", func(t *testing.T) { t.Helper(); testNonblockSPSC(t, caps, ctor) })
		}
	}
	if caps.Has(CapNonblockMPSC) {
		for i := 0; i < *shake; i++ {
			t.Run("n/MPSC", func(t *testing.T) { t.Helper(); testNonblockMPSC(t, caps, ctor) })
		}
	}
	if caps.Has(CapNonblockSPMC) {
		for i := 0; i < *shake; i++ {
			t.Run("n/SPMC", func(t *testing.T) { t.Helper(); testNonblockSPMC(t, caps, ctor) })
		}
	}
	if caps.Has(CapNonblockMPMC) {
		for i := 0; i < *shake; i++ {
			t.Run("n/MPMC", func(t *testing.T) { t.Helper(); testNonblockMPMC(t, caps, ctor) })
		}
	}
}

// Benchmarks runs queue benchmarks for queues
func Benchmarks(b *testing.B, ctor func() Queue) {
	caps := Detect(ctor())
	if !caps.Any(CapQueue) {
		b.Fatal("does not implement any of queue interfaces")
	}
	b.Helper()

	benchCommon(b, caps, ctor)

	// blocking implementations

	if caps.Has(CapBlockSPSC) {
		b.Run("b/SPSC", func(b *testing.B) { b.Helper(); benchSPSC(b, caps, ctor) })
	}
	if caps.Has(CapBlockMPSC) {
		b.Run("b/MPSC", func(b *testing.B) { b.Helper(); benchMPSC(b, caps, ctor) })
	}
	if caps.Has(CapBlockSPMC) {
		b.Run("b/SPMC", func(b *testing.B) { b.Helper(); benchSPMC(b, caps, ctor) })
	}
	if caps.Has(CapBlockMPMC) {
		b.Run("b/MPMC", func(b *testing.B) { b.Helper(); benchMPMC(b, caps, ctor) })
	}

	// non-blocking implementations

	if caps.Has(CapNonblockSPSC) {
		b.Run("n/SPSC", func(b *testing.B) { b.Helper(); benchNonblockSPSC(b, caps, ctor) })
	}
	if caps.Has(CapNonblockMPSC) {
		b.Run("n/MPSC", func(b *testing.B) { b.Helper(); benchNonblockMPSC(b, caps, ctor) })
	}
	if caps.Has(CapNonblockSPMC) {
		b.Run("n/SPMC", func(b *testing.B) { b.Helper(); benchNonblockSPMC(b, caps, ctor) })
	}
	if caps.Has(CapNonblockMPMC) {
		b.Run("n/MPMC", func(b *testing.B) { b.Helper(); benchNonblockMPMC(b, caps, ctor) })
	}
}
