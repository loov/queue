package extqueue

import (
	"testing"

	"github.com/loov/queue/testsuite"
)

func Test(t *testing.T)      { All.Test(t, testsuite.Tests) }
func Benchmark(b *testing.B) { All.Benchmark(b, testsuite.Benchmarks) }
