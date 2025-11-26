// Package simpletransport is an implementation of the Transport interface.
//
// It simply calls the methods from nodes without
// making any actual communication.
package simpletransport

import (
	"github.com/sanjayJ369/raft-consensus/internal/node"
	"github.com/sanjayJ369/raft-consensus/internal/types"
)

type Network map[types.NodeId]*node.Node

func (n *Network) Register(nodeId types.NodeId, node *node.Node) {
	(*n)[node.Id] = node
}

type Transport struct {
	NodeId types.NodeId
	Hub    *Network
	lgr    types.Logger
}

func NewSimpleTransport(nodeId types.NodeId,
	hub *Network,
	lgr types.Logger) *Transport {
	transport := &Transport{
		NodeId: nodeId,
		Hub:    hub,
		lgr:    lgr,
	}

	return transport
}

func (t *Transport) SendVoteRequest(id types.NodeId, req types.VoteRequest) {
	t.lgr.Logf("Sending Vote Request, From: %v, \t To: %v", t.NodeId, id)

	peer, ok := (*t.Hub)[id]
	if !ok {
		// there is no such node
		t.lgr.Logf("There is no such Node, id:%v", id)
		return
	}

	resp := peer.HandleVoteRequest(req)

	self, ok := (*t.Hub)[t.NodeId]
	if !ok || self == nil {
		t.lgr.Logf("Origin node not found in hub, id:%v", t.NodeId)
		return
	}
	self.HandleVoteResponse(resp)
}
