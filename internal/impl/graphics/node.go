package graphics

import (
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/sanjayJ369/raft-consensus/internal/core/node"
	"github.com/sanjayJ369/raft-consensus/internal/core/types"
	"github.com/sanjayJ369/raft-consensus/internal/impl/simple"
	statemachine "github.com/sanjayJ369/raft-consensus/internal/impl/stateMachine"
)

// is a wrapper around node with some additional fields
// needed for rendering
type NodeGUI struct {
	Node        *node.Node
	Pos         rl.Vector2
	Radius      float32
	timerLoader *CircularLoader
	Manager     *Manager
	Color       rl.Color
	Move        bool
	idWidth     int32
}

func NewNodeGUI(width, height int32, lgr types.Logger) *NodeGUI {
	lgr.Logf("creating new node")

	id := types.NewNodeID()
	sm := statemachine.NewKVStore(lgr)
	timer := simple.NewSimpleTimer()

	config := node.Config{
		HeartBeatTimeout:   3000 * time.Millisecond,
		ElectionTimeoutMin: 10000 * time.Millisecond,
		ElectionTimeoutMax: 15000 * time.Millisecond,
	}

	transport := NewMessageTransport(id)
	coreNode := node.NewNode(id, sm, config, timer, transport, lgr)

	radius := float32(80)

	x := rand.Float32()*(float32(width)-2*radius) + radius
	y := rand.Float32()*(float32(height)-2*radius) + radius

	n := &NodeGUI{
		Node:        coreNode,
		Pos:         rl.NewVector2(x, y),
		Radius:      radius,
		timerLoader: NewCircularLoader(x, y, radius),
	}

	idText := fmt.Sprintf("ID: %v", n.Node.ID)
	n.idWidth = rl.MeasureText(idText, 10)

	return n
}

func (n *NodeGUI) HandleMessage(message *Message) {
	switch message.Type {
	case VOTE_REQUEST_MESSAGE:
		req, ok := message.Payload.(types.VoteRequest)
		if !ok {
			log.Fatalln("invalid message type")
		}
		n.Node.HandleVoteRequest(req)
	case VOTE_RESPONSE_NO_MESSAGE, VOTE_RESPONSE_YES_MESSAGE:
		resp, ok := message.Payload.(types.VoteResponse)
		if !ok {
			log.Fatalln("invalid message type")
		}
		n.Node.HandleVoteResponse(resp)
	case HEARTBEAT_MESSAGE, APPEND_ENTRIES_REQUEST_MESSAGE:
		req, ok := message.Payload.(types.AppendEntriesRequest)
		if !ok {
			log.Fatalln("invalid message type")
		}
		n.Node.HandleAppendEntriesRequest(req)
	case APPEND_ENTRIES_RESPONSE_MESSAGE:
		resp, ok := message.Payload.(types.AppendEntriesResponse)
		if !ok {
			log.Fatalln("invalid message type")
		}
		n.Node.HandleAppendEntriesResponse(resp)
	}
}

func (n *NodeGUI) Update() {
	angle := n.Node.ElectionProgress() * 360
	n.timerLoader.Update(float32(angle))

	// on mouse click down move
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) &&
		rl.CheckCollisionPointCircle(rl.GetMousePosition(), n.Pos, n.Radius) {
		n.Move = true
		n.Pos = rl.GetMousePosition()
		n.timerLoader.Pos = n.Pos
	}

	if n.Move && rl.IsMouseButtonUp(rl.MouseButtonLeft) &&
		rl.CheckCollisionPointCircle(rl.GetMousePosition(), n.Pos, n.Radius) {
		n.Move = false
	}
}

func (n *NodeGUI) Render() {
	// ---------------------------------------------------------
	// 1. Setup & Styling
	// ---------------------------------------------------------
	nodeBg := rl.NewColor(30, 30, 30, 255) // Dark Grey
	statsColor := rl.NewColor(150, 150, 150, 255)

	x := int32(n.Pos.X)
	y := int32(n.Pos.Y)
	r := n.Radius

	// ---------------------------------------------------------
	// 2. Logic to Text Conversion
	// ---------------------------------------------------------
	var stateText string
	var stateColor rl.Color

	// Assuming 0=Follower, 1=Candidate, 2=Leader
	switch n.Node.State {
	case node.Follower: // Follower
		stateText = "FOLLOWER"
		stateColor = rl.NewColor(100, 100, 100, 255) // Grey
	case node.Candidate: // Candidate
		stateText = "CANDIDATE"
		stateColor = rl.NewColor(0, 121, 241, 255) // Blue
	case node.Leader: // Leader
		stateText = "LEADER"
		stateColor = rl.Gold
	default:
		stateText = "UNKNOWN"
		stateColor = rl.Red
	}

	// ---------------------------------------------------------
	// 3. Drawing Body & Rings (Optimized)
	// ---------------------------------------------------------
	// Glow (State Color)
	rl.DrawCircleGradient(x, y, r*1.2, rl.Fade(stateColor, 0.3), rl.Fade(rl.Black, 0.0))

	// Solid Body
	rl.DrawCircleV(n.Pos, r, nodeBg)

	// Selection/Timer Ring
	// OPTIMIZATION: Reduced segments from 60 -> 30 to save CPU
	rl.DrawRing(n.Pos, r-3, r, 0, 360, 30, rl.Fade(rl.White, 0.1))

	// Timer Loader
	var loaderColor rl.Color = rl.SkyBlue
	if n.Node.State == node.Leader {
		loaderColor = rl.Red // heartbeat timeout
	}
	n.timerLoader.Render(loaderColor)

	// ---------------------------------------------------------
	// 4. Typography
	// ---------------------------------------------------------

	// -- A. ID (Top) --
	idText := fmt.Sprintf("ID: %v", n.Node.ID)
	idSize := int32(10)

	// Fix: We calculate width here (or use n.idWidth if you cached it in the struct)
	// The previous error was passing 'n.idWidth' as the text string.
	idWidth := rl.MeasureText(idText, idSize)
	rl.DrawText(idText, x-idWidth/2, y-(int32(r)/2)-5, idSize, statsColor)

	// -- B. State (Center) --
	stateSize := int32(16)
	stateWidth := rl.MeasureText(stateText, stateSize)
	rl.DrawText(stateText, x-stateWidth/2, y-8, stateSize, stateColor)

	// -- C. Divider --
	lineY := y + 12
	rl.DrawLine(x-int32(r/2), lineY, x+int32(r/2), lineY, rl.Fade(rl.White, 0.2))

	// -- D. Stats (Bottom) --
	statSize := int32(10)

	stats1 := fmt.Sprintf("TERM: %d  VOTES: %d", n.Node.Term, n.Node.Votes)
	stats1Width := rl.MeasureText(stats1, statSize)
	rl.DrawText(stats1, x-stats1Width/2, lineY+6, statSize, statsColor)

	stats2 := fmt.Sprintf("LOG: %d   CMT: %d", len(n.Node.Log), n.Node.ComittedIndex)
	stats2Width := rl.MeasureText(stats2, statSize)
	rl.DrawText(stats2, x-stats2Width/2, lineY+18, statSize, statsColor)

	// ---------------------------------------------------------
	// 5. Drag Feedback
	// ---------------------------------------------------------
	if n.Move {
		// OPTIMIZATION: Reduced segments 60 -> 30
		rl.DrawRing(n.Pos, r+1, r+3, 0, 360, 30, rl.White)
	}
}
