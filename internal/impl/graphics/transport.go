package graphics

import "github.com/sanjayJ369/raft-consensus/internal/core/types"

// type Transport interface {
// 	SendVoteRequest(NodeId, VoteRequest) sends vote to nodeId, with voterequest body
// 	SendVoteResponse(NodeId, VoteResponse) sends vote response to the node NodeId
// }

// the request and responses should be in the form of packets

type MessageTransaport struct {
	ID      types.NodeID
	manager *Manager
}

func NewMessageTransport(id types.NodeID) *MessageTransaport {
	return &MessageTransaport{ID: id}
}

func (mt MessageTransaport) SendVoteRequest(dest types.NodeID, req types.VoteRequest) {
	mt.manager.SendVoteRequest(mt.ID, dest, req)
}

func (mt MessageTransaport) SendVoteResponse(dest types.NodeID, resp types.VoteResponse) {
	mt.manager.SendVoteResponse(mt.ID, dest, resp)
}

func (mt MessageTransaport) SendAppendEntriesRequest(dest types.NodeID, req types.AppendEntriesRequest) {
	mt.manager.SendAppendEntriesRequest(mt.ID, dest, req)
}

func (mt MessageTransaport) SendAppendEntriesResponse(dest types.NodeID, resp types.AppendEntriesResponse) {
	mt.manager.SendAppendEntriesResponse(mt.ID, dest, resp)
}
