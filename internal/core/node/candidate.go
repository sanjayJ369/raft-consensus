package node

import (
	"math"

	"github.com/sanjayJ369/raft-consensus/internal/core/types"
	"github.com/sanjayJ369/raft-consensus/internal/impl/log"
)

// StartLeader sets the state of the node to a candidate
func (n *Node) EnterCandidate() {
	n.State = Candidate
	// stop election timeout timer as currently in candidate state
	n.lgr.Logf("Entering Canidate State: nodeId: %v", n.ID)
	n.ResetElectionTimer()

	// once the node enters a canidate state
	// start a new election term
	n.StartNewElectionTerm()
}

func (n *Node) StartNewElectionTerm() {
	n.Lock()
	defer n.Unlock()

	n.lgr.Logf("Starting New Election Term: nodeId: %v, oldterm: %d", n.ID, n.Term)
	// it should increment it's election term
	n.Term += 1
	n.VotedFor = &n.ID // vote itself
	n.Votes = 1

	// todo: ask for the votes from other nodes
	var prevLog log.LogEntry
	prevlogIdx := len(n.Log) - 1
	if prevlogIdx >= 0 {
		prevLog = n.Log[prevlogIdx] // get the most recent log
	}

	for _, nodeId := range n.PeerIDs {
		go n.Transport.SendVoteRequest(nodeId, types.VoteRequest{
			CandidateId:  n.ID,
			FollowerId:   nodeId,
			Term:         n.Term,
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
	n.Lock()

	n.lgr.Logf("Recevied Response from Node: %v", res.From)

	// receives a response from a legit leader
	if res.Term > n.Term {
		n.Unlock()
		n.EnterFollower()
		return
	}

	// if recevied majority of the votes
	// become leader
	majoryReq := math.Ceil(float64(n.NodesInCluster) / 2)
	if res.VoteGranted {
		n.Votes++
	}

	if n.Votes >= int(majoryReq) {
		n.lgr.Logf("!!!... received majority of vote.... becoming leader")
		n.Unlock()
		n.EnterLeader()
		return
	}

	n.Unlock()
}
