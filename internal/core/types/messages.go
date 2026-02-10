package types

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

// Append Entries Requset
// sent from leader to followers
type AppendEntriesRequest struct {
	LeaderID          NodeID
	LeaderTerm        Term
	PrevLogIndex      Index
	PrevLogTerm       Term
	Entries           []LogEntry
	LeaderCommitIndex Index
}

type AppendEntriesResponse struct {
	Term    Term
	Success bool
}
