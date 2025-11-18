package types

type Logger interface {
	Logf(format string, args ...any)
}
