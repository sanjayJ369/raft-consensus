package types

type Term int64   // represents election term
type Index int64  // represetns log index
type NodeId int64 // represents node Id

type Vote struct {
	Term        Term   // current election term
	VoteGranted bool   // is the vote granted or not
	From        NodeId // follower node id
	To          NodeId // candidate node id
}

// VoteRequest
type VoteRequest struct {
	CanidateId   NodeId
	FollowerId   NodeId
	PrevLogTerm  Term  // term of the last log entry
	PrevLogIndex Index // index of the last log entry
	Term         Term  // current election term
}

// VoteResponse
type VoteResponse Vote
