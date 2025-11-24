package node

import (
	"time"

	"github.com/sanjayJ369/raft-consensus/internal/log"
	simpletimer "github.com/sanjayJ369/raft-consensus/internal/simpleTimer"
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

type NodeState int

const (
	Follower NodeState = iota
	Candidate
	Leader
)

// Node represents the state of each server node
// in a general distributed systems setting
type Node struct {
	// node and cluster info
	state          NodeState
	Id             types.NodeId         // Node Id
	smachine       statemachine.KVStore // state machine which is simple key value store
	peerIDs        []types.NodeId       // other nodes in the cluster
	nodesInCluster int                  // number of peers + 1
	config         Config               // stores all the config fiels

	transport types.Transport // way to communicate with othern odes
	lgr       types.Logger    // logger to log....

	// node states
	log           []log.LogEntry // replicated log store
	electionTimer types.Timer    // election timer
	term          types.Term     // current election term
	votes         int            // votes of the current term
	votedFor      *types.NodeId  // to whom did the node vote for in the current term
	lastApplied   types.Index    // last log entry that is being applied to the state machine
	comittedIndex types.Index    // highest log entry that is known to be comitted

	// leader specific states
	nextIndex  map[types.NodeId]types.Index // index to be shared with the follower
	matchIndex map[types.NodeId]types.Index
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

func NewNode(Id types.NodeId,
	stateMachine statemachine.KVStore,
	peerIDs []types.NodeId,
	nodesInCluster int,
	config Config,
	transport types.Transport,
	lgr types.Logger) *Node {
	return &Node{
		state:          Follower, // every node starts as follower
		Id:             Id,
		smachine:       stateMachine,
		peerIDs:        peerIDs,
		nodesInCluster: nodesInCluster,
		config:         config,
		transport:      transport,
		lgr:            lgr,

		log:           make([]log.LogEntry, 100), // initally reserve like 100 log entries
		electionTimer: simpletimer.NewSimpleTimer(),
		term:          0, // start from term zero
		votes:         0,
		votedFor:      nil, // not yet voted
		lastApplied:   -1,
		comittedIndex: -1,

		// leader specific states
		nextIndex:  make(map[types.NodeId]types.Index),
		matchIndex: make(map[types.NodeId]types.Index),
	}
}
