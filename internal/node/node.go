package node

import (
	"time"

	"github.com/sanjayJ369/raft-consensus/internal/log"
	statemachine "github.com/sanjayJ369/raft-consensus/internal/stateMachine"
	"github.com/sanjayJ369/raft-consensus/internal/types"
)

// each node contains  several things
// replicated log
// state machine
// consensus module
// and few variables to hold the state of the machine

type Config struct {
	HeartBeatTimeout   time.Duration
	ElectionTimeoutMin time.Duration
	ElectionTimeoutMax time.Duration
}

// represents election term
type Term int64

// Node represents the state of each server node
// in a general distributed systems setting
type Node struct {
	Id            types.NodeId         // Node Id
	smachine      statemachine.KVStore // state machine which is simple key value store
	log           []log.LogEntry       // replicated log store
	peers         []*Node              // other nodes in the cluseter
	lgr           types.Logger         // logger to log....
	electionTimer types.Timer          // election timer
	config        Config               // stores all the config fiels
	term          types.Term           // current election term
	votes         []types.Vote         // votes of the current term
	votedFor      types.NodeId         // to whom did the node vote for in the current term
}

// StartNewElection starts a new election
// increment it's term
// ask for votes to all other nodes

// there are three main things to implement
// 1. Leader Election
// 2. Log Replication
// 3. safety mechanism

// main RPCs mentioned in the paper (http://nil.csail.mit.edu/6.824/2020/papers/raft-extended.pdf)
// RequestVote
// AppendEntriesRPC
// InstallSnapshot RPC
