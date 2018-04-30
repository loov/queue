package extqueue

import (
	"math/bits"
	"unsafe"
)

// Value is the data being sent
type Value = int64

// Node for using intrusive implementations
type Node struct {
	next  unsafe.Pointer
	Value Value
}

type seqValue struct {
	sequence int64
	value    Value
}
type seqPaddedValue struct {
	sequence int64
	value    Value
	_        [8 - 2]int64
}

type seqValue32 struct {
	sequence uint32
	value    Value
}

type seqPaddedValue32 struct {
	sequence uint32
	value    Value
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
