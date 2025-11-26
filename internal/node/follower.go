package node

import (
	"time"

	"github.com/sanjayJ369/raft-consensus/internal/types"
	"github.com/sanjayJ369/raft-consensus/utils"
)

// EnterFollower sets the state of the node to follower
func (n *Node) EnterFollower() {
	n.lgr.Logf("Entered Follower State: %v", n.Id)
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
	randDuration := utils.RandomRangeInt64(
		int64(n.config.ElectionTimeoutMin),
		int64(n.config.ElectionTimeoutMax))

	duration := time.Duration(randDuration) * time.Nanosecond
	n.lgr.Logf("started election timeout timer duration: %d", duration)
	go n.electionTimer.Start(duration, n.EnterCandidate)
}

func (n *Node) ResetElectionTimer() {
	n.lgr.Logf("Reset Election Timer NodeId: %v", n.Id)
	n.electionTimer.Stop()
	n.StartElectionTimer()
}

// HandleVoteRequest handles the vote request received from the candidate
// this is usually invoked by Transport.SendRequestVote
func (n *Node) HandleVoteRequest(req types.VoteRequest) types.VoteResponse {
	n.Lock()
	defer n.Unlock()

	n.lgr.Logf("Received Vote Request From: %v, \t Request: %v", req.CanidateId, req)

	// if already voted in this term don't vote again
	if n.term >= req.Term {
		// don't grant vote
		return donotgrantVote(n, req, "candidate has lower or equal term")
	}

	lastLogIndex := len(n.log) - 1
	// no log entries yet.. grant an vote
	if lastLogIndex < 0 {
		return grantVote(n, req, "there is no last log index")
	}

	lastLogEntry := n.log[lastLogIndex]

	// grant an vote
	// if the candidate has newer logs vote for them
	// (safety mechanism)
	// vote only if the follower has not
	// previously voted in current term
	if req.PrevLogTerm >= lastLogEntry.Term {
		// last log entry has greater term
		if req.PrevLogTerm > lastLogEntry.Term {
			return grantVote(n, req, "candidate has newer term")
		} else if req.Term == lastLogEntry.Term &&
			lastLogEntry.Index >= int(req.PrevLogIndex) {
			// same term then longer logs are considered newer
			return grantVote(n, req, "candidate has longer logs")
		} else {
			return donotgrantVote(n, req, "candidate has shorter logs")
		}
	}

	return donotgrantVote(n, req, "candidate has older term")
}

func donotgrantVote(n *Node, req types.VoteRequest, reason string) types.VoteResponse {
	vote := types.Vote{
		Term:        req.Term,
		VoteGranted: false,
		From:        n.Id,
		To:          req.CanidateId,
	}
	n.lgr.Logf("%d: rejecting vote to candidate: %d, reason: %s", n.Id, req.CanidateId, reason)
	return types.VoteResponse(vote)
}

func grantVote(n *Node, req types.VoteRequest, reason string) types.VoteResponse {
	vote := types.Vote{
		Term:        req.Term,
		VoteGranted: true,
		From:        n.Id,
		To:          req.CanidateId,
	}

	n.votedFor = &req.CanidateId
	n.term = req.Term
	n.lgr.Logf("%d: granting vote to candidate : %d, reason: %s", n.Id, req.CanidateId, reason)
	n.ResetElectionTimer()

	return types.VoteResponse(vote)
}
