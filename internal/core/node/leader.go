package node

// StartLeader changes the state from candidate to leader
func (n *Node) EnterLeader() {
	n.State = Leader
	n.ElectionTimer.Stop()
	n.lgr.Logf("entering leader state")
	n.lgr.Logf("sucessfully elected a new leader")
}
