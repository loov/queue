package testsuite

import (
	"strconv"
	"testing"
)

// Descs is a list of Queue implementation descriptions
type Descs []*Desc

// Desc describes a Queue implementation
type Desc struct {
	Name   string
	Param  CreateParam
	Create func(batchSize int, size int) Queue
}

func (desc *Desc) HasSizeParam() bool {
	return desc.Param == ParamSize || desc.Param == ParamBatchSizeAndSize
}
func (desc *Desc) HasBatchSizeParam() bool {
	return desc.Param == ParamBatchSize || desc.Param == ParamBatchSizeAndSize
}

type CreateParam int

const (
	ParamNone = CreateParam(iota)
	ParamSize
	ParamBatchSize
	ParamBatchSizeAndSize
)

func (descs *Descs) Append(list ...*Desc) {
	*descs = append(*descs, list...)
}

func (descs Descs) TestDefault(t *testing.T) {
	t.Helper()
	descs.Test(t, Tests)
}

func (descs Descs) BenchmarkDefault(b *testing.B) {
	b.Helper()
	descs.Benchmark(b, Benchmarks)
}

func (descs Descs) Test(t *testing.T, test func(t *testing.T, create func() Queue)) {
	t.Helper()
	for _, desc := range descs {
		batchSizes := BatchSizes
		if !desc.HasBatchSizeParam() {
			batchSizes = []int{0}
		}

		testSizes := TestSizes
		if desc.HasSizeParam() {
			testSizes = []int{0}
		}

		t.Run(desc.Name, func(t *testing.T) {
			t.Helper()
			for _, batchSize := range batchSizes {
				for _, size := range testSizes {
					if size != 0 && size <= batchSize {
						continue
					}

					name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
					t.Run(name, func(t *testing.T) {
						t.Helper()
						test(t, func() Queue {
							return desc.Create(batchSize, size)
						})
					})
				}
			}
		})
	}
}

func (descs Descs) Benchmark(b *testing.B, bench func(b *testing.B, create func() Queue)) {
	b.Helper()
	for _, desc := range descs {
		batchSizes := BenchBatchSizes
		if !desc.HasBatchSizeParam() {
			batchSizes = []int{0}
		}

		benchSizes := BenchSizes
		if desc.HasSizeParam() {
			benchSizes = []int{0}
		}

		b.Run(desc.Name, func(b *testing.B) {
			b.Helper()
			for _, batchSize := range batchSizes {
				for _, size := range benchSizes {
					if size != 0 && size <= batchSize {
						continue
					}

					name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
					b.Run(name, func(b *testing.B) {
						b.Helper()
						bench(b, func() Queue {
							return desc.Create(batchSize, size)
						})
					})
				}
			}
		})
	}
}
