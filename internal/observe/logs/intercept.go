package logs

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
)

// NomadRequestLogger returns a unary interceptor that logs request information.
//
// TODO: log the request after the response is formulated, this will allow us to set the log level
// based on the response status code. Also, instead of manually returning ConnectError from function,
// we return the QuantmError and let the ConnectError be created by the logging interceptor.
func NomadRequestLogger() connect.UnaryInterceptorFunc {
	intercept := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			slog.Info(req.Spec().Procedure, "method", req.HTTPMethod())

			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(intercept)
}
