package queue

import (
	"strconv"
	"testing"

	. "github.com/loov/queue/testsuite"
	"github.com/tidwall/fastlane"
)

// MPSC wrapper
var _ MPSC = (*MPSCnw_fl)(nil)

type MPSCnw_fl struct {
	ch fastlane.ChanUint64
}

func NewMPSCnw_fl() *MPSCnw_fl {
	return &MPSCnw_fl{}
}

func (q *MPSCnw_fl) MultipleProducers() {}

func (q *MPSCnw_fl) Send(v Value) bool {
	q.ch.Send(uint64(v))
	return true
}

func (q *MPSCnw_fl) Recv(v *Value) bool {
	*v = Value(q.ch.Recv())
	return true
}

func TestMPSCnw_fl(t *testing.T) {
	batchSize := 0
	size := 0
	name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
	t.Run(name, func(t *testing.T) {
		Tests(t, func() Queue {
			return NewMPSCnw_fl()
		})
	})
}

func BenchmarkMPSCnw_fl(b *testing.B) {
	batchSize := 0
	size := 0
	name := strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
	b.Run(name, func(b *testing.B) {
		Benchmarks(b, func() Queue {
			return NewMPSCnw_fl()
		})
	})
}
