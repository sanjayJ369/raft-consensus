package node

import "os"

// StartLeader changes the state from candidate to leader
func (n *Node) EnterLeader() {
	n.state = Leader
	n.electionTimer.Stop()
	n.lgr.Logf("entering leader state")
	n.lgr.Logf("sucessfully elected a new leader")
	os.Exit(1)
}
