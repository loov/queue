package queue

import "unsafe"

// Exposed node for intrusive implementations
type Node struct {
	next  unsafe.Pointer
	Value Value
}
