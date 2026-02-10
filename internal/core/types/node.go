package types

import "sync/atomic"

type Term int64   // represents election term
type Index int64  // represetns log index
type NodeID int64 // represents node Id

var globalNodeID int64

func NewNodeID() NodeID {
	return NodeID(atomic.AddInt64(&globalNodeID, 1))
}
