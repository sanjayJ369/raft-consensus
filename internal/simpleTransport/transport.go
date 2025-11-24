// Package simpletransport is an implementation of the Transport interface.
//
// It simply calls the methods from nodes without
// making any actual communication.
package simpletransport

import (
	"github.com/sanjayJ369/raft-consensus/internal/node"
	"github.com/sanjayJ369/raft-consensus/internal/types"
)

type Transport struct {
	Node  *node.Node                  //  node it belongs to
	Nodes map[types.NodeId]*node.Node // pointers to ther nodes
	lgr   types.Logger                // just a logger ;)
}

func NewSimpleTransport(node *node.Node,
	nodes map[types.NodeId]*node.Node,
	lgr types.Logger) *Transport {
	return &Transport{
		Node:  node,
		Nodes: nodes,
		lgr:   lgr,
	}
}

func (t *Transport) SendVoteRequest(id types.NodeId, req types.VoteRequest) {
	t.lgr.Logf("Sending Vote Request, From: %v, \t To: %v", t.Node.Id, id)
	node, ok := t.Nodes[id]
	if !ok {
		// there is no such node
		t.lgr.Logf("There is no such Node, id:%v", id)
		return
	}

	resp := node.HandleVoteRequest(req)
	t.Node.HandleVoteResponse(resp)
}
