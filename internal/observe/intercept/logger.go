package intercept

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
)

// RequestLogger returns a unary interceptor that logs request and response information only once,
// using the procedure name as the log message and logging errors at ERROR level and successes at INFO level.
func RequestLogger() connect.UnaryInterceptorFunc {
	intercept := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			peer := req.Peer()
			procedure := req.Spec().Procedure
			fields := []any{
				"peer_address", peer.Addr,
				"protocol", peer.Protocol,
			}

			resp, err := next(ctx, req)

			elapsed := time.Since(start)
			fields = append(fields, "elapsed", elapsed)

			if err != nil {
				fields = append(fields, "error", err.Error())
				if connectErr, ok := err.(*connect.Error); ok {
					fields = append(fields, "connect_code", connectErr.Code())
					fields = append(fields, "connect_details", connectErr.Details())
				}

				slog.Error(procedure, fields...)
			} else {
				slog.Info(procedure, fields...)
			}

			return resp, err
		})
	}

	return connect.UnaryInterceptorFunc(intercept)
}
