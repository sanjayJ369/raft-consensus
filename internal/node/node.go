package node

import (
	"github.com/sanjayJ369/raft-consensus/internal/log"
	statemachine "github.com/sanjayJ369/raft-consensus/internal/stateMachine"
	"github.com/sanjayJ369/raft-consensus/internal/types"
)

// each node contains  several things
// replicated log
// state machine
// consensus module
// and few variables to hold the state of the machine

// Node represents the state of each server node
// in a general distributed systems setting
type Node struct {
	smachine statemachine.KVStore // state machine which is simple key value store
	log      []log.LogEntry       // replicated log store
	nodes    []*Node              // other nodes in the cluseter
	lgr      types.Logger         // logger to log....
}
