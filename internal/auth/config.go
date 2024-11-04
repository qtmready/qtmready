package auth

import (
	"sync/atomic"
)

const (
	_secret string = "set me"
)

var (
	_defaultsecret atomic.Value
)

func init()                { _defaultsecret.Store(_secret) }
func Secret() string       { return _defaultsecret.Load().(string) }
func SetSecret(val string) { _defaultsecret.Store(val) }
