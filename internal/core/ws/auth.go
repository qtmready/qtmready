package ws

import (
	"context"
)

type (
	AuthFn func(context.Context, string) (string, error)
)

// noop is a default auth function that returns the token as is.
func noop(_ context.Context, token string) (string, error) {
	return token, nil
}
