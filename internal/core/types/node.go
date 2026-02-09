package types

import "sync/atomic"

type Term int64   // represents election term
type Index int64  // represetns log index
type NodeID int64 // represents node Id

var globalNodeID int64

func NewNodeID() NodeID {
	return NodeID(atomic.AddInt64(&globalNodeID, 1))
}

// ******************** MESSAGES ******************** //

// VoteRequest
type VoteRequest struct {
	CandidateId  NodeID
	FollowerId   NodeID
	PrevLogTerm  Term  // term of the last log entry
	PrevLogIndex Index // index of the last log entry
	Term         Term  // current election term
}

// VoteResponse
type VoteResponse struct {
	To          NodeID // candidate node id
	From        NodeID // follower node id
	Term        Term   // current election term
	VoteGranted bool   // is the vote granted or not
}
