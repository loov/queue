package extqueue

import (
	"math/bits"
	"unsafe"
)

// Node for using intrusive implementations
type Node[T any] struct {
	next  unsafe.Pointer
	Value T
}

type seqValue[T any] struct {
	sequence int64
	value    T
}

type seqPaddedValue[T any] struct {
	sequence int64
	value    T
	_        [7]int64
}
type seqValue32[T any] struct {
	sequence uint32
	value    T
}

type seqPaddedValue32[T any] struct {
	sequence uint32
	value    T
	_        [8 - 2]int32
}

func ceil(a, n int) int {
	r := ((a + n - 1) / n) * n
	if r <= n {
		return n * 2
	}
	return r
}

func nextPowerOfTwo(v uint32) uint32 {
	return 1 << uint(32-bits.LeadingZeros32(v))
}
