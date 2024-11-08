package erratic

import (
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/grpc/codes"
)

func CodeToProto(code int) codes.Code {
	switch code {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusInternalServerError:
		return codes.Internal
	default:
		return codes.Unknown
	}
}

func CodeToConnect(code int) connect.Code {
	switch code {
	case http.StatusBadRequest:
		return connect.CodeInvalidArgument
	case http.StatusUnauthorized:
		return connect.CodeUnauthenticated
	case http.StatusForbidden:
		return connect.CodePermissionDenied
	case http.StatusNotFound:
		return connect.CodeNotFound
	case http.StatusInternalServerError:
		return connect.CodeInternal
	default:
		return connect.CodeUnknown
	}
}
