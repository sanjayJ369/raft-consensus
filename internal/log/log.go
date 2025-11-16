package log

type LogEntry struct {
	Term  int    //  election term
	Index int    // log index
	Entry string // the command
}
