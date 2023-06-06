package queue

import (
	"errors"
)

var (
	ErrNoParentNoQueue = errors.New("both parent and queue are nil")
)
