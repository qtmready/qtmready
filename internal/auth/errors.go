package auth

import (
	"errors"
)

var (
	ErrInvalidAPIKey         = errors.New("invalid API key")
	ErrInvalidAuthHeader     = errors.New("invalid authorization header")
	ErrInvalidOrExpiredToken = errors.New("invalid or expired token")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrMalformedAPIKey       = errors.New("malformed API key")
	ErrMissingAuthHeader     = errors.New("no authorization header provided")
	ErrCrypto                = errors.New("crypto error")
	ErrTokenExpired          = errors.New("token has expired")
)
