package queue

import "unsafe"

// Node for using intrusive implementations
type Node struct {
	next  unsafe.Pointer
	Value Value
}
