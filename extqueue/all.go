package extqueue

import (
	"testing"

	"github.com/loov/queue/testsuite"
)

//go:generate go run all_gen.go -out all_test.go

type Desc struct {
	Name   string
	Flags  DescFlag
	Create func(batchSize int, size int) testsuite.Queue
}

func (desc *Desc) HasBatchSize() bool { return desc.Flags&Batched == Batched }
func (desc *Desc) IsUnbouned() bool   { return desc.Flags&Unbounded == Unbounded }

type DescFlag int

const (
	Default = DescFlag(1 << iota)
	Batched
	Unbounded
)

var Descs = []Desc{
	{"MPMCcGo", Default, func(bs, s int) testsuite.Queue { return NewMPMCcGo(s) }},
	{"MPMCqGo", Default, func(bs, s int) testsuite.Queue { return NewMPMCqGo(s) }},
	{"MPMCqpGo", Default, func(bs, s int) testsuite.Queue { return NewMPMCqpGo(s) }},

	// {"SPSCrMC", Batched, func(bs, s int) testsuite.Queue { return NewSPSCrMC(bs, s) }
	// {"SPSCrsMC", Batched, func(bs, s int) testsuite.Queue { return NewSPSCrsMC(bs, s) }
	// {"MPSCrMC", Batched, func(bs, s int) testsuite.Queue { return NewMPSCrMC(bs, s) }
	// {"MPSCrsMC", Batched, func(bs, s int) testsuite.Queue { return NewMPSCrsMC(bs, s) }

	{"SPSCnsDV", Unbounded, func(bs, s int) testsuite.Queue { return NewSPSCnsDV() }},
	{"MPSCnsDV", Unbounded, func(bs, s int) testsuite.Queue { return NewMPSCnsDV() }},
	{"MPSCnsiDV", Unbounded, func(bs, s int) testsuite.Queue { return NewMPSCnsiDV() }},

	{"MPMCqsDV", Default, func(bs, s int) testsuite.Queue { return NewMPMCqsDV(s) }},
	{"MPMCqspDV", Default, func(bs, s int) testsuite.Queue { return NewMPMCqspDV(s) }},
	{"SPMCqsDV", Default, func(bs, s int) testsuite.Queue { return NewSPMCqsDV(s) }},
	{"SPMCqspDV", Default, func(bs, s int) testsuite.Queue { return NewSPMCqspDV(s) }},
	{"SPSCqsDV", Default, func(bs, s int) testsuite.Queue { return NewSPSCqsDV(s) }},
	{"SPSCqspDV", Default, func(bs, s int) testsuite.Queue { return NewSPSCqspDV(s) }},
	// {"SPSCqspDV2", Default, func(bs, s int) testsuite.Queue { return NewSPSCqspDV2(s) }},
}
