package queue

import "testing"

func benchSPSC(b *testing.B, caps Capability, ctor func() Queue) {}
func benchMPSC(b *testing.B, caps Capability, ctor func() Queue) {}
func benchSPMC(b *testing.B, caps Capability, ctor func() Queue) {}
func benchMPMC(b *testing.B, caps Capability, ctor func() Queue) {}
