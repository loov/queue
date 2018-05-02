// Package extqueue contains many different concurrent queue algorithms, tests and benchmarks
//
// All names follow  a convention: "[SM]P[SM]C[sw]?i?p?[rnacq]<variant>""
//
//     [SM]P:
//        supports either single `S` or multiple `M` concurrent producers
//
//     [SM]C:
//        supports either single `S` or multiple `M` concurrent consumers
//
//     [rnacq]: buffer implementation
//        `r` dynamically sized ring buffer
//        `n` node based,
//        `a` fixed size array based,
//        `c` channel based,
//        `q` dynamically sized ring buffer with sequence number.
//
//     [sw]?: waiting behavior
//        `` : when it is a waiting implementation (no CPU burn)
//        `s`: when it is a spinning implementation,
//        `w`: when it is partially spinning and partially waiting
//
//     p?: memory usage
//        `` : usually one value per bounded size (sometimes with one uint64 or uint32)
//        `p`: value padded to a cacheline
//
//     i?: intrusiveness
//        `` : doesn't have intrusiveness
//        `i`: has an intrusive API, where producer can provide a node to store the value
//
//     <variant>:
//        special variant identifier for a particular implementation,
//        which indicates either base implementation author / paper / code.
//
// Guideline for selecting an implementation:
//
// 1. Select the minimal producers and consumers that you need.
// For example, if you always have single consumer, but unknown number of producers
// you would need an MPSC queue.
//
// 2. Select which waiting behavior you need. Spinning implementations tend to be faster however
// they do burn CPU while they spin, this can start affecting everything else in the system.
// If you do not have real-time requirements then using a non-spinning or partially spinning
// implementation, is probably better, because the queue is a better citizen.
//
// 3. When you need to copy large values, then you using an intrusive implementation
// can be helpful. It allows to allocate and fill the sent node by letting the producer
// create the appropriate node. Of course this means that the API is somewhat more annoying to use.
//
// 4. When you are concerned with memory usage you should try to avoid node based and
// padded implementations. They can use 2x to 10x more memory to store the queue.
//
// 5. With unbounded queues (don't have the size parameter) you should be very careful.
// Unbounded queues do not create back-pressure in case the consumer isn't able to handle all
// the incoming values. This can lead to out-of-memory situations.
//
// 6. Some queues here support batching. They tend to be faster, however care must be taken
// to properly flush the batches with FlushSend and FlushRecv, otherwise the queue can deadlock.
//
// However, the most reliable way is to write a realistic benchmark for your situtation and
// see what works the best. This package contains a convenient way to implement them.
//
//    func Benchmark(b *testing.B) {
//    	// iterate over all implementations
//    	for _, desc := range extqueue.Descs {
//    		// detect which properties do we need
//    		if !desc.MultipleProducers() || desc.Unbounded() {
//    			continue
//    		}
//
//    		// try different batch sizes
//    		batchSizes := testsuite.BenchBatchSizes
//    		if !desc.BatchSize() {
//    			batchSizes = []int{0}
//    		}
//
//    		// run a sub-benchmark for this implementation
//    		b.Run(desc.Name, func(b *testing.B) {
//    			b.Helper()
//    			// iterate over all combinations of batchSizes and testSizes
//    			for _, batchSize := range batchSizes {
//    				for _, size := range testsuite.BenchSizes {
//    					if size <= batchSize {
//    						continue
//    					}
//
//    					// use a naming scheme to identify between
//    					name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
//    					b.Run(name, func(b *testing.B) {
//    						b.Helper()
//    						RunBenchmarks(b, func() testsuite.Queue {
//    							return desc.Create(batchSize, size)
//    						})
//    					})
//    				}
//    			}
//    		})
//    	}
//    }
//
//    // write a benchmark for the queues, this should be outside of this function
//    func RunBenchmarks(b *testing.B, create func() testsuite.Queue) {
//    	// use sub benchmarks to test a particular aspect
//    	b.Run("Basic", func(b *testing.B) {
//    		q := create().(testsuite.MPSC)
//    		testsuite.ProducerConsumerBenchmark(b,
//    			4, 1, // 4 producers and 1 consumer
//    			func(int) {
//    				for i := 0; i < b.N; i++ {
//    					q.Send(Value(i))
//    				}
//    			}, func(int) {
//    				for i := 0; i < 4*b.N; i++ {
//    					var v Value
//    					q.Recv(&v)
//    				}
//    			})
//    	})
//    }
package extqueue

import (
	"github.com/loov/queue/testsuite"
)

////go:generate go run all_gen.go -out all_test.go

type Desc struct {
	Name   string
	Flags  DescFlag
	Create func(batchSize int, size int) testsuite.Queue
	q      testsuite.Queue
}

func (desc *Desc) BatchSize() bool { return desc.Flags&Batched == Batched }
func (desc *Desc) Unbounded() bool { return desc.Flags&Unbounded == Unbounded }

func (desc *Desc) MultipleProducers() bool {
	_, ok := desc.q.(interface{ MultipleProducers() })
	return ok
}
func (desc *Desc) MultipleConsumers() bool {
	_, ok := desc.q.(interface{ MultipleConsumers() })
	return ok
}

type DescFlag int

const (
	Default = DescFlag(1 << iota)
	Batched
	Unbounded
)

var Descs = []Desc{
	{"MPMCcGo", Default, func(bs, s int) testsuite.Queue { return NewMPMCcGo(s) }, nil},
	{"MPMCqGo", Default, func(bs, s int) testsuite.Queue { return NewMPMCqGo(s) }, nil},
	{"MPMCqpGo", Default, func(bs, s int) testsuite.Queue { return NewMPMCqpGo(s) }, nil},

	// {"SPSCrMC", Batched, func(bs, s int) testsuite.Queue { return NewSPSCrMC(bs, s) }, nil},
	{"SPSCrsMC", Batched, func(bs, s int) testsuite.Queue { return NewSPSCrsMC(bs, s) }, nil},
	// {"MPSCrMC", Batched, func(bs, s int) testsuite.Queue { return NewMPSCrMC(bs, s) }, nil},
	{"MPSCrsMC", Batched, func(bs, s int) testsuite.Queue { return NewMPSCrsMC(bs, s) }, nil},

	{"SPSCnsDV", Unbounded, func(bs, s int) testsuite.Queue { return NewSPSCnsDV() }, nil},
	{"MPSCnsDV", Unbounded, func(bs, s int) testsuite.Queue { return NewMPSCnsDV() }, nil},
	{"MPSCnsiDV", Unbounded, func(bs, s int) testsuite.Queue { return NewMPSCnsiDV() }, nil},

	{"MPMCqsDV", Default, func(bs, s int) testsuite.Queue { return NewMPMCqsDV(s) }, nil},
	{"MPMCqspDV", Default, func(bs, s int) testsuite.Queue { return NewMPMCqspDV(s) }, nil},
	{"SPMCqsDV", Default, func(bs, s int) testsuite.Queue { return NewSPMCqsDV(s) }, nil},
	{"SPMCqspDV", Default, func(bs, s int) testsuite.Queue { return NewSPMCqspDV(s) }, nil},
	{"SPSCqsDV", Default, func(bs, s int) testsuite.Queue { return NewSPSCqsDV(s) }, nil},
	{"SPSCqspDV", Default, func(bs, s int) testsuite.Queue { return NewSPSCqspDV(s) }, nil},
}

func init() {
	for i := range Descs {
		Descs[i].q = Descs[i].Create(2, 2)
	}
}
