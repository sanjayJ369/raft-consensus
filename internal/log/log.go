package log

import (
	"github.com/sanjayJ369/raft-consensus/internal/types"
)

type LogEntry struct {
	Term  types.Term   //  election term
	Index int          // log index
	Entry string       // the command
	lgr   types.Logger // logger
}
