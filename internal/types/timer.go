package types

import "time"

type Timer interface {
	Stop()
	Start(time.Duration, func())
}
