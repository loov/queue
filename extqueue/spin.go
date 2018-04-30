package extqueue

import "runtime"

type spinT256 struct{ count int }

func (s *spinT256) Spin() bool {
	s.count++
	if s.count > 256 {
		runtime.Gosched()
		s.count = 0
	}
	return true
}
