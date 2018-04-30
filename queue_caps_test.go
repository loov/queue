package queue

import (
	"strings"
	"testing"
)

// Capability is bitflag for different queue capabilities
type Capability uint64

// Add adds a capability
func (caps *Capability) Add(cap Capability) { *caps |= cap }

// Any detects whether any of the flags are set
func (caps Capability) Any(cap Capability) bool { return caps&cap != 0 }

// Has detects whether has the specified capability
func (caps Capability) Has(cap Capability) bool { return caps&cap == cap }

func (caps Capability) String() string {
	xs := []string{}
	if caps.Has(CapBlockSPSC) {
		xs = append(xs, "BlockSPSC")
	}
	if caps.Has(CapBlockMPSC) {
		xs = append(xs, "BlockMPSC")
	}
	if caps.Has(CapBlockSPMC) {
		xs = append(xs, "BlockSPMC")
	}
	if caps.Has(CapNonblockSPSC) {
		xs = append(xs, "NonblockSPSC")
	}
	if caps.Has(CapNonblockMPSC) {
		xs = append(xs, "NonblockMPSC")
	}
	if caps.Has(CapNonblockSPMC) {
		xs = append(xs, "NonblockSPMC")
	}
	if caps.Has(CapBounded) {
		xs = append(xs, "Bounded")
	}
	return "[" + strings.Join(xs, ", ") + "]"
}

const (
	CapBlockSPSC = Capability(1 << iota)
	CapBlockMPSC = CapBlockSPSC | Capability(1<<iota)
	CapBlockSPMC
	CapNonblockSPSC = Capability(1 << iota)
	CapNonblockMPSC = CapNonblockSPSC | Capability(1<<iota)
	CapNonblockSPMC
	CapBounded = Capability(1 << iota)
	// BatchReceiver

	CapBlockMPMC    = CapBlockMPSC | CapBlockSPMC
	CapNonblockMPMC = CapNonblockMPSC | CapNonblockSPMC

	CapQueue = CapBlockMPMC | CapNonblockMPMC
)

func Detect(q Queue) Capability {
	var caps Capability
	if _, ok := q.(SPSC); ok {
		caps.Add(CapBlockSPSC)
	}
	if _, ok := q.(MPSC); ok {
		caps.Add(CapBlockMPSC)
	}
	if _, ok := q.(SPMC); ok {
		caps.Add(CapBlockSPMC)
	}
	if _, ok := q.(NonblockingSPSC); ok {
		caps.Add(CapNonblockSPSC)
	}
	if _, ok := q.(NonblockingMPSC); ok {
		caps.Add(CapNonblockMPSC)
	}
	if _, ok := q.(NonblockingSPMC); ok {
		caps.Add(CapNonblockSPMC)
	}
	if _, ok := q.(Bounded); ok {
		caps.Add(CapBounded)
	}
	return caps
}

func TestCapability(t *testing.T) {
	if !CapQueue.Any(CapQueue) {
		t.Fatal("!CapQueue.Any(CapQueue)")
	}
	if !CapBlockSPSC.Any(CapQueue) {
		t.Fatal("!CapBlockSPSC.Any(CapQueue)")
	}
}
