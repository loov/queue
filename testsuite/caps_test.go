package testsuite

import (
	"testing"
)

func TestCapability(t *testing.T) {
	if !CapQueue.Any(CapQueue) {
		t.Fatal("!CapQueue.Any(CapQueue)")
	}
	if !CapBlockSPSC.Any(CapQueue) {
		t.Fatal("!CapBlockSPSC.Any(CapQueue)")
	}
}
