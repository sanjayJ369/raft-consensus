package node

import "github.com/sanjayJ369/raft-consensus/internal/types"

// StartLeader sets the state of the node to a candidate
func (n *Node) EnterCandidate() {
	// once the node enters a canidate state
	// start a new election term
	n.StartNewElectionTerm()

}

func (n *Node) StartNewElectionTerm() {
	// it should increment it's election term
	n.term += 1
	n.votedFor = n.Id // vote itself

	// todo: ask for the votes from other nodes

}

// VoteRequest
type VoteRequest struct {
	canidateId   types.NodeId
	followerId   types.NodeId
	prevLogTerm  types.Term
	prevLogIndex types.Index
}

// VoteResponse
type VoteResponse types.Vote

func (n *Node) HandleRequestVote(req VoteRequest) VoteResponse {
	// grant an vote
	return types.Vote{}
}
