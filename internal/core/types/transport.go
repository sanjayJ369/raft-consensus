package types

// Transport handles the communication between nodes
type Transport interface {
	SendVoteRequest(NodeId, VoteRequest)
}
