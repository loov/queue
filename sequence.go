package queue

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
