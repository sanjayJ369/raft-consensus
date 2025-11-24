package node

import (
	"time"

	"github.com/sanjayJ369/raft-consensus/internal/types"
	"github.com/sanjayJ369/raft-consensus/utils"
)

// EnterFollower sets the state of the node to follower
func (n *Node) EnterFollower() {
	// every node starts as a follwer
	// things the follower must do

	//  start election timeout timer
	n.StartElectionTimer()

	// todo: reset voteFor in current term
	// todo: handle appendEntries request
	// todo: handle RequestVote request
	// todo: handle installSnapshot request
}

// Start starts election timeout timer
func (n *Node) StartElectionTimer() {
	duration := utils.RandomRangeInt64(
		int64(n.config.ElectionTimeoutMin),
		int64(n.config.ElectionTimeoutMax))

	n.lgr.Logf("started election timeout timer duration: %d", duration)
	go n.electionTimer.Start(time.Duration(duration), n.EnterCandidate)
}

func (n *Node) ResetElectionTimer() {
	n.electionTimer.Stop()
	n.StartElectionTimer()
}

// HandleVoteRequest handles the vote request received from the candidate
// this is usually invoked by Transport.SendRequestVote
func (n *Node) HandleVoteRequest(req types.VoteRequest) types.VoteResponse {
	// grant an vote
	// if the candidate has newer logs vote for them
	// (safety mechanism)
	lastLogEntry := n.log[len(n.log)]
	vote := types.Vote{
		Term:        req.Term,
		VoteGranted: false,
		From:        n.Id,
		To:          req.CanidateId,
	}

	// vote only if the follower has not
	// previously voted in current term
	if n.votedFor == nil && req.PrevLogTerm >= lastLogEntry.Term {
		// last log entry has greater term
		if req.PrevLogTerm > lastLogEntry.Term {
			vote.VoteGranted = true
		} else if req.Term == lastLogEntry.Term &&
			lastLogEntry.Index >= int(req.PrevLogIndex) {
			// same term then longer logs are considered newer
			vote.VoteGranted = true
		} else {
			vote.VoteGranted = false
		}
	}

	return types.VoteResponse(vote)
}
