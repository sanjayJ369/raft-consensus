package node

import "github.com/sanjayJ369/raft-consensus/internal/core/types"

func (n *Node) HandleAppendEntriesRequest(req types.AppendEntriesRequest) {
	// just reset timer for now
	n.Timer.Restart()

	// todo handle logic properly
}

func (n *Node) HandleAppendEntriesResponse(resp types.AppendEntriesResponse) {
	// todo handle append entery response
}
