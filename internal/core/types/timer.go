package types

import "time"

type Timer interface {
	Stop()
	Start(time.Duration, func()) // on timeout call the given function
	Restart()
}
