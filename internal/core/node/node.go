package node

import (
	"sync"
	"time"

	"github.com/sanjayJ369/raft-consensus/internal/core/types"
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
	sync.Mutex // handle concurrent vote requests

	// node and cluster info
	State          NodeState
	ID             types.NodeID   // Node Id
	StateMachine   types.DB       // state machine which is simple key value store
	PeerIDs        []types.NodeID // other nodes in the cluster
	NodesInCluster int            // number of peers + 1
	Config         Config         // stores all the config fiels

	Transport types.Transport // way to communicate with othern odes
	lgr       types.Logger    // logger to log....

	// node states
	Log           []types.LogEntry // replicated log store
	Timer         types.Timer      // election timer or heartbeat timer
	Term          types.Term       // current election term
	Votes         int              // votes of the current term
	VotedFor      *types.NodeID    // to whom did the node vote for in the current term
	LastApplied   types.Index      // last log entry that is being applied to the state machine
	ComittedIndex types.Index      // highest log entry that is known to be comitted

	// leader specific states
	nextIndex  map[types.NodeID]types.Index // index to be shared with the follower
	matchIndex map[types.NodeID]types.Index
}

func (n *Node) AddPeer(id types.NodeID) {
	n.PeerIDs = append(n.PeerIDs, id)
	n.NodesInCluster++
}

func (n *Node) ElectionProgress() float64 {
	elapsed := n.Timer.Elapsed()
	duration := n.Timer.Duration()

	if duration == 0 {
		return 0
	}

	return float64(elapsed) / float64(duration)
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

func NewNode(Id types.NodeID,
	stateMachine types.DB,
	config Config,
	timer types.Timer,
	transport types.Transport,
	lgr types.Logger) *Node {
	return &Node{
		State:          Follower, // every node starts as follower
		ID:             Id,
		StateMachine:   stateMachine,
		NodesInCluster: 1,
		PeerIDs:        []types.NodeID{},
		Config:         config,
		lgr:            lgr,
		Transport:      transport,

		Log:           make([]types.LogEntry, 0, 100), // initally reserve like 100 log entries
		Timer:         timer,
		Term:          0, // start from term zero
		Votes:         0,
		VotedFor:      nil, // not yet voted
		LastApplied:   -1,
		ComittedIndex: -1,

		// leader specific states
		nextIndex:  make(map[types.NodeID]types.Index),
		matchIndex: make(map[types.NodeID]types.Index),
	}
}
