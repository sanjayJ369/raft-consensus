package types

// Transport handles the communication between nodes
type Transport interface {
	SendVoteRequest(NodeID, VoteRequest)
	SendVoteResponse(NodeID, VoteResponse)
}
