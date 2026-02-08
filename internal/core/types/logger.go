package types

type Logger interface {
	Logf(format string, args ...any)

	// if we are logging to file we want to sync data from
	// buffer to file
	Sync()
}
