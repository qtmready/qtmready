package github

import "errors"

var (
	ErrorNoEventToParse               = errors.New("no event specified to parse")
	ErrorInvalidHttpMethod            = errors.New("invalid HTTP Method")
	ErrorMissingHeaderGithubEvent     = errors.New("missing X-GitHub-Event Header")
	ErrorMissingHeaderGithubSignature = errors.New("missing X-Hub-Signature Header")
	ErrorInvalidEvent                 = errors.New("event not defined to be parsed")
	ErrorPayloadParser                = errors.New("error parsing payload")
	ErrorVerifySignature              = errors.New("HMAC verification failed")
)
