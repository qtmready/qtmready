package intercepts

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"

	"go.breu.io/quantm/internal/erratic"
)

// RequestLogger returns a unary interceptor that logs request and response information.
//
// It logs the peer address, protocol, HTTP method, and latency.
// If the request returns an error, it logs the error details, including the error ID, code, message, hints, and
// internal error (if any).
// It converts any error to a connect.Error using erratic.QuantmError.ToConnectError.
// The procedure name is used as the log message. Errors are logged at ERROR level, and successes are logged at
// INFO level.
func RequestLogger() connect.UnaryInterceptorFunc {
	intercept := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			peer := req.Peer()
			procedure := req.Spec().Procedure
			fields := []any{
				"peer_address", peer.Addr,
				"protocol", peer.Protocol,
				"method", req.HTTPMethod(),
			}

			resp, err := next(ctx, req)

			elapsed := time.Since(start)
			fields = append(fields, "latency", elapsed)

			if err != nil {
				qerr, ok := err.(*erratic.QuantmError)
				if !ok {
					qerr = erratic.NewUnknownError(erratic.CommonModule).Wrap(err)
				}

				fields = append(fields, "error_id", qerr.ID)
				fields = append(fields, "error_code", qerr.Code)
				fields = append(fields, "error", qerr.Error())

				for k, v := range qerr.Hints {
					fields = append(fields, k, v)
				}

				if qerr.Unwrap() != nil {
					fields = append(fields, "internal_error", qerr.Unwrap().Error())
				}

				slog.Warn(procedure, fields...)

				err = qerr.ToConnectError()
			} else {
				slog.Info(procedure, fields...)
			}

			return resp, err
		})
	}

	return connect.UnaryInterceptorFunc(intercept)
}
