package cloudrun

import (
	"sync"
)

type (
	Constructor struct{}
)

var (
	registerOnce sync.Once
)
