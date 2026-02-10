package node

import "github.com/sanjayJ369/raft-consensus/internal/core/types"

// StartLeader changes the state from candidate to leader
func (n *Node) EnterLeader() {
	n.State = Leader
	n.Timer.Stop()
	n.lgr.Logf("entering leader state")
	n.lgr.Logf("sucessfully elected a new leader, %d", n.ID)
	n.SendHeartBeatMessages()
	n.StartHeartBeatTimer()
}

func (n *Node) StartHeartBeatTimer() {
	n.Timer.Start(n.Config.HeartBeatTimeout, func() {
		n.SendHeartBeatMessages()
		if n.State == Leader {
			n.StartHeartBeatTimer()
		}
	})
}

func (n *Node) SendHeartBeatMessages() {
	message := n.getHeartBeatRequest()
	for _, peerId := range n.PeerIDs {
		n.Transport.SendAppendEntriesRequest(peerId, message)
	}
}

func (n *Node) getHeartBeatRequest() types.AppendEntriesRequest {
	// todo include all other appropriate fields
	return types.AppendEntriesRequest{
		LeaderTerm: n.Term,
	}
}
