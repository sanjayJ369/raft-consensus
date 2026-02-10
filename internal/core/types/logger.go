package types

type LogEntry struct {
	Term  Term   //  election term
	Index int    // log index
	Entry string // the command
	lgr   Logger // logger
}

type Logger interface {
	Logf(format string, args ...any)

	// if we are logging to file we want to sync data from
	// buffer to file
	Sync()
}
