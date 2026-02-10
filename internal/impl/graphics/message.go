package graphics

import (
	"log"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	"github.com/sanjayJ369/raft-consensus/internal/core/types"
)

type MessageType int

const (
	HEARTBEAT_MESSAGE MessageType = iota
	VOTE_REQUEST_MESSAGE
	VOTE_RESPONSE_YES_MESSAGE
	VOTE_RESPONSE_NO_MESSAGE
	APPEND_ENTRIES_REQUEST_MESSAGE
	APPEND_ENTRIES_RESPONSE_MESSAGE
)

const (
	SIDEL_LEN_MESSAGE = 20
	SPEED             = 500

	INCREMENT   = 0
	DECEREMENT  = 1
	TIMER_RESET = 2
)

var (
	HEARTBEAT_COLOR           = rl.SkyBlue
	VOTE_REQUEST_COLOR        = rl.Orange
	VOTE_RESPONSE_YES_COLOR   = rl.Green
	VOTE_RESPONSE_NO_COLOR    = rl.Red
	APPEND_ENTRIES_REQ_COLOR  = rl.Purple
	APPEND_ENTRIES_RESP_COLOR = rl.Magenta
)

type Message struct {
	ID      string
	SrcID   types.NodeID // source node id
	DestID  types.NodeID // dest node id
	Type    MessageType
	Payload any

	Pos          rl.Vector2
	Dest         rl.Vector2
	SideLen      float32
	Color        rl.Color
	Rect         rl.Rectangle
	CreationTime time.Time
}

func (m *Message) Update() {
	delta := rl.Vector2Subtract(m.Dest, m.Pos)
	if rl.Vector2Length(delta) > 0 {
		dir := rl.Vector2Normalize(delta)
		vel := rl.Vector2Scale(dir, SPEED*rl.GetFrameTime())
		m.Pos = rl.Vector2Add(m.Pos, vel)
	}

	m.Rect = rl.NewRectangle(m.Pos.X, m.Pos.Y, m.SideLen, m.SideLen)
}

func NewMessage(
	m *Manager,
	from, to types.NodeID,
	msgType MessageType,
	payload any,
) *Message {

	srcNode, ok := m.Nodes[from]
	if !ok {
		log.Fatalf("source node %v not registered", from)
	}

	destNode, ok := m.Nodes[to]
	if !ok {
		log.Fatalf("destination node %v not registered", to)
	}

	pos := srcNode.Pos
	dest := destNode.Pos

	msg := &Message{
		ID:           uuid.NewString(),
		SrcID:        from,
		DestID:       to,
		Type:         msgType,
		Payload:      payload,
		Pos:          pos,
		Dest:         dest,
		SideLen:      SIDEL_LEN_MESSAGE,
		Color:        colorForMessageType(msgType),
		Rect:         rl.NewRectangle(pos.X, pos.Y, SIDEL_LEN_MESSAGE, SIDEL_LEN_MESSAGE),
		CreationTime: time.Now(),
	}

	m.AddMessage(msg)

	return msg
}

func (m *Message) Render() {
	// ---------------------------------------------------------
	// 1. Setup & Dimensions
	// ---------------------------------------------------------
	// Use the SideLen to determine size, but make it rectangular (packet shape)
	width := m.SideLen * 1.5 // Slightly wider than tall
	height := m.SideLen

	// Center the drawing point
	x := m.Pos.X - width/2
	y := m.Pos.Y - height/2

	rec := rl.NewRectangle(x, y, width, height)

	// ---------------------------------------------------------
	// 2. Draw "Motion Trail" (The Tail)
	// ---------------------------------------------------------
	// Calculate direction vector to draw the tail behind the message
	// Direction = Normalized(Dest - Pos)
	dx := m.Dest.X - m.Pos.X
	dy := m.Dest.Y - m.Pos.Y
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if length > 0 {
		// Normalize direction
		dirX, dirY := dx/length, dy/length

		// Draw 3 fading circles behind the packet to simulate speed
		for i := 1; i <= 3; i++ {
			dist := float32(i) * 8.0 // spacing
			trailPos := rl.NewVector2(
				m.Pos.X-(dirX*dist),
				m.Pos.Y-(dirY*dist),
			)
			// Fade out further back
			alpha := 0.4 - (float32(i) * 0.1)
			rl.DrawCircleV(trailPos, 3, rl.Fade(m.Color, alpha))
		}
	}

	// ---------------------------------------------------------
	// 3. Draw The Packet (Body)
	// ---------------------------------------------------------
	// Glow effect
	rl.DrawRectangleRounded(
		rl.NewRectangle(x-2, y-2, width+4, height+4),
		0.5, 10,
		rl.Fade(m.Color, 0.3),
	)

	// Solid Core
	rl.DrawRectangleRounded(rec, 0.4, 10, rl.NewColor(20, 20, 20, 255))

	// Border
	rl.DrawRectangleRoundedLines(rec, 0.4, 10, m.Color)

	// ---------------------------------------------------------
	// 4. Text Info (Type Abbreviation)
	// ---------------------------------------------------------
	// Instead of full text "RequestVote", we map it to short codes
	// to keep the packet small and readable.
	shortText := getShortMsgType(m.Type) // Helper function below
	fontSize := int32(10)

	textW := rl.MeasureText(shortText, fontSize)

	rl.DrawText(
		shortText,
		int32(m.Pos.X)-textW/2,
		int32(m.Pos.Y)-fontSize/2,
		fontSize,
		rl.White,
	)
}

// Helper to shorten message types for the UI
// You should add this function or integrate it into your logic
func getShortMsgType(t MessageType) string {
	switch t {
	case HEARTBEAT_MESSAGE:
		return "HB"
	case VOTE_REQUEST_MESSAGE:
		return "RV"
	case VOTE_RESPONSE_NO_MESSAGE:
		return "VRN"
	case VOTE_RESPONSE_YES_MESSAGE:
		return "VRY"
	case APPEND_ENTRIES_REQUEST_MESSAGE:
		return "AE"
	case APPEND_ENTRIES_RESPONSE_MESSAGE:
		return "AER"
	default:
		return "?"
	}
}

func colorForMessageType(t MessageType) rl.Color {
	switch t {
	case HEARTBEAT_MESSAGE:
		return HEARTBEAT_COLOR
	case VOTE_REQUEST_MESSAGE:
		return VOTE_REQUEST_COLOR
	case VOTE_RESPONSE_YES_MESSAGE:
		return VOTE_RESPONSE_YES_COLOR
	case VOTE_RESPONSE_NO_MESSAGE:
		return VOTE_RESPONSE_NO_COLOR
	case APPEND_ENTRIES_REQUEST_MESSAGE:
		return APPEND_ENTRIES_REQ_COLOR
	case APPEND_ENTRIES_RESPONSE_MESSAGE:
		return APPEND_ENTRIES_RESP_COLOR
	default:
		return rl.Gray
	}
}
