package auth

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
)

func NomadAuthenticator() connect.UnaryInterceptorFunc {
	intercept := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			header := req.Header().Get("Authorization")

			if strings.HasPrefix(header, "Bearer ") {
				token := strings.TrimPrefix(header, "Bearer ")

				cliams, err := DecodeJWE(Secret(), token)
				if err != nil {
					return nil, connect.NewError(connect.CodeUnauthenticated, err)
				}

				if cliams != nil {
					ctx = context.WithValue(ctx, AuthContextUser, cliams.UserID)
					ctx = context.WithValue(ctx, AuthContextOrg, cliams.OrgID)
				}
			} else {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing bearer token"))
			}

			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(intercept)
}
