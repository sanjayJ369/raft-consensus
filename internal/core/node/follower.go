package node

import (
	"time"

	"github.com/sanjayJ369/raft-consensus/internal/core/types"
	"github.com/sanjayJ369/raft-consensus/utils"
)

// EnterFollower sets the state of the node to follower
func (n *Node) EnterFollower() {
	n.lgr.Logf("Entered Follower State: %v", n.ID)
	// every node starts as a follwer
	// things the follower must do

	//  start election timeout timer
	n.StartElectionTimer()
	n.VotedFor = nil

	// todo: handle appendEntries request
	// todo: handle RequestVote request
	// todo: handle installSnapshot request
}

// Start starts election timeout timer
func (n *Node) StartElectionTimer() {
	randDuration := utils.RandomRangeInt64(
		int64(n.Config.ElectionTimeoutMin),
		int64(n.Config.ElectionTimeoutMax))

	duration := time.Duration(randDuration) * time.Nanosecond
	n.lgr.Logf("started election timeout timer duration: %d", duration)
	n.ElectionTimer.Start(duration, n.EnterCandidate)
}

func (n *Node) ResetElectionTimer() {
	n.lgr.Logf("Reset Election Timer NodeId: %v", n.ID)
	n.ElectionTimer.Stop()
	n.StartElectionTimer()
}

func (n *Node) stepDown(newTerm types.Term) {
	n.Term = newTerm
	n.VotedFor = nil
	n.Votes = 0
	n.EnterFollower()
	n.ResetElectionTimer()
}

// HandleVoteRequest handles the vote request received from the candidate
// this is usually invoked by Transport.SendRequestVote
func (n *Node) HandleVoteRequest(req types.VoteRequest) {
	n.Lock()
	defer n.Unlock()

	n.lgr.Logf("Received Vote Request From: %v, Request: %+v", req.CandidateId, req)

	granted := false
	reason := ""

	// Reject stale term
	if req.Term < n.Term {
		reason = "stale term"
	} else {

		//  Update term if newer
		if req.Term > n.Term {
			n.stepDown(req.Term)
		}

		// Check if already voted
		if n.VotedFor != nil && *n.VotedFor != req.CandidateId {
			reason = "already voted"
		} else {

			// Check log freshness
			lastIndex := len(n.Log) - 1
			var lastTerm types.Term

			if lastIndex >= 0 {
				lastTerm = n.Log[lastIndex].Term
			}

			// check is candidate is more upto date
			upToDate :=
				req.PrevLogTerm > lastTerm ||
					(req.PrevLogTerm == lastTerm &&
						req.PrevLogIndex >= types.Index(lastIndex))

			if upToDate {
				granted = true
				n.VotedFor = &req.CandidateId
				n.ResetElectionTimer()
				reason = "vote granted"
			} else {
				reason = "log not up to date"
			}
		}
	}

	n.Transport.SendVoteResponse(req.CandidateId, types.VoteResponse{
		Term:        n.Term,
		VoteGranted: granted,
		From:        n.ID,
		To:          req.CandidateId,
	})

	n.lgr.Logf("Vote response to %v: granted=%v (%s)", req.CandidateId, granted, reason)
}
