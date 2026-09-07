package coordinator

import "errors"

var (
	ErrNoActiveNodes     = errors.New("no active nodes available for placement")
	ErrInsufficientSpace = errors.New("not enough space available for chunk placement")
)
