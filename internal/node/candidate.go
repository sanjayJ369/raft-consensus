package node

import (
	"math"

	"github.com/sanjayJ369/raft-consensus/internal/types"
)

// StartLeader sets the state of the node to a candidate
func (n *Node) EnterCandidate() {
	// once the node enters a canidate state
	// start a new election term
	n.StartNewElectionTerm()

}

func (n *Node) StartNewElectionTerm() {
	// it should increment it's election term
	n.term += 1
	n.votedFor = &n.Id // vote itself
	n.votes = 1

	// todo: ask for the votes from other nodes
	prevLog := n.log[len(n.log)-1] // get the most recent log
	for _, nodeId := range n.peerIDs {
		go n.transport.SendVoteRequest(nodeId, types.VoteRequest{
			CanidateId:   n.Id,
			FollowerId:   nodeId,
			Term:         n.term,
			PrevLogTerm:  prevLog.Term,
			PrevLogIndex: types.Index(prevLog.Index),
		})
	}
}

// HandleVoteResponse processes a VoteResponse received from a follower.
// It is typically invoked by the transport layer when a follower replies
// to this node's vote request. The caller must ensure any required
// synchronization (for example, holding the Node's lock) is in place
// before calling this method.
func (n *Node) HandleVoteResponse(res types.VoteResponse) {
	if n.state != Candidate {
		return
	}

	// if recevied majority of the votes
	// become leader
	majoryReq := math.Ceil(float64(n.nodesInCluster) / 2)
	if res.VoteGranted {
		n.votes++
	}

	if n.votes > int(majoryReq) {
		n.EnterLeader()
	}
}
