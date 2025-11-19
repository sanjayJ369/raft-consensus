package node

import (
	"time"

	"github.com/sanjayJ369/raft-consensus/utils"
)

// EnterFollower sets the state of the node to follower
func (n *Node) EnterFollower() {
	// every node starts as a follwer
	// things the follower must do

	//  start election timeout timer
	n.StartElectionTimer()

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
