//	implementation of the Transport interface.
//
// It simply calls the methods from nodes without
// making any actual communication.
package simple

import (
	"github.com/sanjayJ369/raft-consensus/internal/core/node"
	"github.com/sanjayJ369/raft-consensus/internal/core/types"
)

type Network map[types.NodeID]*node.Node

func (n *Network) Register(nodeId types.NodeID, node *node.Node) {
	(*n)[node.ID] = node
}

type Transport struct {
	NodeId types.NodeID
	Hub    *Network
	lgr    types.Logger
}

func NewSimpleTransport(nodeId types.NodeID,
	hub *Network,
	lgr types.Logger) *Transport {
	transport := &Transport{
		NodeId: nodeId,
		Hub:    hub,
		lgr:    lgr,
	}

	return transport
}

func (t *Transport) SendVoteRequest(id types.NodeID, req types.VoteRequest) {
	t.lgr.Logf("Sending Vote Request, From: %v, \t To: %v", t.NodeId, id)

	peer, ok := (*t.Hub)[id]
	if !ok {
		// there is no such node
		t.lgr.Logf("There is no such Node, id:%v", id)
		return
	}

	// directly calling handler of peer node
	peer.HandleVoteRequest(req)

}

func (t *Transport) SendVoteResponse(id types.NodeID, res types.VoteResponse) {
	peer, ok := (*t.Hub)[id]
	if !ok || peer == nil {
		t.lgr.Logf("Origin node not found in hub, id:%v", t.NodeId)
		return
	}
	peer.HandleVoteResponse(res)
}
