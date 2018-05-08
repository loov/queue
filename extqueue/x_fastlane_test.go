// +build external

package extqueue

import (
	"github.com/loov/queue/testsuite"
	"github.com/tidwall/fastlane"
)

type MPSCnwFL struct {
	ch fastlane.ChanUint64
}

func NewMPSCnwFL() *MPSCnwFL {
	return &MPSCnwFL{}
}

func (q *MPSCnwFL) MultipleProducers() {}

func (q *MPSCnwFL) Send(v Value) bool {
	q.ch.Send(uint64(v))
	return true
}

func (q *MPSCnwFL) Recv(v *Value) bool {
	*v = Value(q.ch.Recv())
	return true
}

func init() {
	All.Append(
		&testsuite.Desc{"MPSCnwFL", testsuite.ParamNone, func(bs, s int) testsuite.Queue { return NewMPSCnwFL() }},
	)
}
