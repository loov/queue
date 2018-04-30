package queue

import "testing"

func testSPSC(t *testing.T, caps Capability, ctor func() Queue) {}
func testMPSC(t *testing.T, caps Capability, ctor func() Queue) {}
func testSPMC(t *testing.T, caps Capability, ctor func() Queue) {}
func testMPMC(t *testing.T, caps Capability, ctor func() Queue) {}
