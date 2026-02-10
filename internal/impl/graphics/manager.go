package graphics

import (
	"log"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/sanjayJ369/raft-consensus/internal/core/types"
)

// Manager handles all the raylib parts of the simulation
type Manager struct {
	sync.Mutex
	lgr      types.Logger
	Nodes    map[types.NodeID]*NodeGUI
	Messages map[string]*Message
}

func NewManager(lgr types.Logger) *Manager {
	return &Manager{
		lgr:      lgr,
		Nodes:    map[types.NodeID]*NodeGUI{},
		Messages: map[string]*Message{},
	}
}

func (m *Manager) StartNodes() {
	for _, node := range m.Nodes {
		node.Node.EnterFollower()
	}
}

func (m *Manager) Update() {

	for _, node := range m.Nodes {
		node.Update()
	}

	var toDeleteMsg []*Message
	for _, message := range m.Messages {
		if time.Since(message.CreationTime) > 5*time.Second {
			toDeleteMsg = append(toDeleteMsg, message)
		}
		message.Update()
	}

	// handle collisions
	for _, node := range m.Nodes {
		for _, message := range m.Messages {
			if message.DestID == node.Node.ID {
				if rl.CheckCollisionCircleRec(node.Pos, node.Radius, message.Rect) {
					// message has reached the node...
					// handle message
					go node.HandleMessage(message)

					// delete message
					toDeleteMsg = append(toDeleteMsg, message)
				}
			}
		}
	}

	for _, msg := range toDeleteMsg {
		m.DeleteMessage(msg)
	}
}

func (m *Manager) Render() {
	for _, node := range m.Nodes {
		node.Render()
	}

	for _, message := range m.Messages {
		message.Render()
	}
}

func (m *Manager) AddMessage(message *Message) {
	m.Lock()
	defer m.Unlock()
	m.Messages[message.ID] = message
}

func (m *Manager) DeleteMessage(message *Message) {
	m.Lock()
	defer m.Unlock()
	delete(m.Messages, message.ID)
}

func (m *Manager) RegisterNode(node *NodeGUI) {
	m.Nodes[node.Node.ID] = node
	mt, ok := node.Node.Transport.(*MessageTransaport)
	if ok {
		mt.manager = m
	} else {
		log.Fatalln("must use MessageTransport to render messages")
	}
}

func (m *Manager) SendVoteRequest(from, to types.NodeID, req types.VoteRequest) {
	NewMessage(m, from, to, VOTE_REQUEST_MESSAGE, req)
	m.lgr.Logf("sent vote request from: %d, to: %d, resp: %v", from, to, req)
}

func (m *Manager) SendVoteResponse(from, to types.NodeID, resp types.VoteResponse) {
	// todo: create a new message with response payaload and add it to manager
	if resp.VoteGranted {
		NewMessage(m, from, to, VOTE_RESPONSE_YES_MESSAGE, resp)
	} else {
		NewMessage(m, from, to, VOTE_RESPONSE_NO_MESSAGE, resp)
	}
	m.lgr.Logf("sent vote response from: %d, to: %d, resp: %v", from, to, resp)
}

func (m *Manager) SendAppendEntriesRequest(from, to types.NodeID, req types.AppendEntriesRequest) {
	if len(req.Entries) == 0 {
		NewMessage(m, from, to, HEARTBEAT_MESSAGE, req) // heartbeat
		m.lgr.Logf("sent heartbeat message from %d, to %d", from, to)
	} else {
		NewMessage(m, from, to, APPEND_ENTRIES_REQUEST_MESSAGE, req)
		m.lgr.Logf("sent append entry message from %d, to %d, entries:%v", from, to, req.Entries)
	}
}

func (m *Manager) SendAppendEntriesResponse(from, to types.NodeID, resp types.AppendEntriesResponse) {
	NewMessage(m, from, to, APPEND_ENTRIES_RESPONSE_MESSAGE, resp)
	m.lgr.Logf("sent append entry response message from:%d, to:%d", from, to)
}
