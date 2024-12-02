package erratic

import (
	"connectrpc.com/connect"
	"google.golang.org/grpc/codes"
)

func CodeToProto(code int) codes.Code {
	_, kind := Decompose(code)

	switch kind {
	case CodeBadRequest:
		return codes.InvalidArgument
	case CodeCancelled:
		return codes.Canceled
	case CodeFailedPrecondition:
		return codes.FailedPrecondition
	case CodeExhausted:
		return codes.ResourceExhausted
	case CodeNotFound:
		return codes.NotFound
	case CodeExists:
		return codes.AlreadyExists
	case CodeCorrupted:
		return codes.Aborted
	case CodeConflict:
		return codes.Aborted
	case CodeAuthentication:
		return codes.Unauthenticated
	case CodePermissionDenied:
		return codes.PermissionDenied
	case CodeSystem:
		return codes.Internal
	case CodeConfig:
		return codes.Internal
	case CodeDatabase:
		return codes.Internal
	case CodeNetwork:
		return codes.Internal
	case CodeUnavailable:
		return codes.Unavailable
	default:
		return codes.Unknown
	}
}

func CodeToConnect(code int) connect.Code {
	_, kind := Decompose(code)

	switch kind {
	case CodeBadRequest:
		return connect.CodeInvalidArgument
	case CodeCancelled:
		return connect.CodeCanceled
	case CodeFailedPrecondition:
		return connect.CodeFailedPrecondition
	case CodeExhausted:
		return connect.CodeResourceExhausted
	case CodeNotFound:
		return connect.CodeNotFound
	case CodeExists:
		return connect.CodeAlreadyExists
	case CodeCorrupted:
		return connect.CodeAborted
	case CodeConflict:
		return connect.CodeAborted
	case CodeAuthentication:
		return connect.CodeUnauthenticated
	case CodePermissionDenied:
		return connect.CodePermissionDenied
	case CodeSystem:
		return connect.CodeInternal
	case CodeConfig:
		return connect.CodeInternal
	case CodeDatabase:
		return connect.CodeInternal
	case CodeNetwork:
		return connect.CodeInternal
	case CodeUnavailable:
		return connect.CodeUnavailable
	default:
		return connect.CodeUnknown
	}
}
