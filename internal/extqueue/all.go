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
//    All.Benchmark(b, func(b *testing.B, create func() testsuite.Queue) {
//    	q := create()
//    	if _, ok := q.(interface{ MultipleProducers() }); !ok {
//    		b.Skip("does not support multiple producers")
//    	}
//    	b.Run("Basic", func(b *testing.B) {
//    		q := create().(testsuite.MPSC)
//    		testsuite.ProducerConsumerBenchmark(b,
//    			4, 1, // 4 producers and 1 consumer
//    			func(int) {
//    				for i := 0; i < b.N; i++ {
//    					q.Send(T(i))
//    				}
//    			}, func(int) {
//    				for i := 0; i < 4*b.N; i++ {
//    					var v T
//    					q.Recv(&v)
//    				}
//    			})
//    	})
//    })
//
package extqueue

import (
	"loov.dev/queue/internal/testsuite"
)

////go:generate go run all_gen.go -out all_test.go

var All = testsuite.Descs{
	{
		Name:   "MPMCcGo",
		Param:  testsuite.ParamSize,
		Create: func(bs, s int) testsuite.Queue { return NewMPMCcGo[testsuite.Value](s) }},
	{
		Name:   "MPMCqGo",
		Param:  testsuite.ParamSize,
		Create: func(bs, s int) testsuite.Queue { return NewMPMCqGo[testsuite.Value](s) }},
	{
		Name:   "MPMCqpGo",
		Param:  testsuite.ParamSize,
		Create: func(bs, s int) testsuite.Queue { return NewMPMCqpGo[testsuite.Value](s) }},

	{
		Name:   "SPSCrMC",
		Param:  testsuite.ParamBatchSizeAndSize,
		Create: func(bs, s int) testsuite.Queue { return NewSPSCrMC[testsuite.Value](bs, s) }},
	{
		Name:   "SPSCrsMC",
		Param:  testsuite.ParamBatchSizeAndSize,
		Create: func(bs, s int) testsuite.Queue { return NewSPSCrsMC[testsuite.Value](bs, s) }},
	{
		Name:   "MPSCrMC",
		Param:  testsuite.ParamBatchSizeAndSize,
		Create: func(bs, s int) testsuite.Queue { return NewMPSCrMC[testsuite.Value](bs, s) }},
	{
		Name:   "MPSCrsMC",
		Param:  testsuite.ParamBatchSizeAndSize,
		Create: func(bs, s int) testsuite.Queue { return NewMPSCrsMC[testsuite.Value](bs, s) }},

	{
		Name:   "SPSCnsDV",
		Param:  testsuite.ParamNone,
		Create: func(bs, s int) testsuite.Queue { return NewSPSCnsDV[testsuite.Value]() }},
	{
		Name:   "MPSCnsDV",
		Param:  testsuite.ParamNone,
		Create: func(bs, s int) testsuite.Queue { return NewMPSCnsDV[testsuite.Value]() }},
	{
		Name:   "MPSCnsiDV",
		Param:  testsuite.ParamNone,
		Create: func(bs, s int) testsuite.Queue { return NewMPSCnsiDV[testsuite.Value]() }},

	{
		Name:   "MPMCqsDV",
		Param:  testsuite.ParamSize,
		Create: func(bs, s int) testsuite.Queue { return NewMPMCqsDV[testsuite.Value](s) }},
	{
		Name:   "MPMCqspDV",
		Param:  testsuite.ParamSize,
		Create: func(bs, s int) testsuite.Queue { return NewMPMCqspDV[testsuite.Value](s) }},
	{
		Name:   "SPMCqsDV",
		Param:  testsuite.ParamSize,
		Create: func(bs, s int) testsuite.Queue { return NewSPMCqsDV[testsuite.Value](s) }},
	{
		Name:   "SPMCqspDV",
		Param:  testsuite.ParamSize,
		Create: func(bs, s int) testsuite.Queue { return NewSPMCqspDV[testsuite.Value](s) }},
	{
		Name:   "MPSCqsDV",
		Param:  testsuite.ParamSize,
		Create: func(bs, s int) testsuite.Queue { return NewMPSCqsDV[testsuite.Value](s) }},
	{
		Name:   "MPSCqspDV",
		Param:  testsuite.ParamSize,
		Create: func(bs, s int) testsuite.Queue { return NewMPSCqspDV[testsuite.Value](s) }},
	{
		Name:   "SPSCqsDV",
		Param:  testsuite.ParamSize,
		Create: func(bs, s int) testsuite.Queue { return NewSPSCqsDV[testsuite.Value](s) }},
	{
		Name:   "SPSCqspDV",
		Param:  testsuite.ParamSize,
		Create: func(bs, s int) testsuite.Queue { return NewSPSCqspDV[testsuite.Value](s) }},
}
