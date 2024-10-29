package errors

import (
	"errors"
)

var (
	ErrNoEventToParse               = errors.New("no event specified to parse")
	ErrInvalidHTTPMethod            = errors.New("invalid HTTP Method")
	ErrMissingHeaderGithubEvent     = errors.New("missing X-GitHub-Event Header")
	ErrMissingHeaderGithubSignature = errors.New("missing X-Hub-Signature Header")
	ErrInvalidEvent                 = errors.New("event not defined to be parsed")
	ErrPayloadParser                = errors.New("error parsing payload")
	ErrVerifySignature              = errors.New("HMAC verification failed")
)
