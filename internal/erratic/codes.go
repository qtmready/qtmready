package erratic

const (
	CodeBadRequest         int = 100
	CodeCancelled          int = 101
	CodeFailedPrecondition int = 102
	CodeExhausted          int = 103
	CodeNotFound           int = 200
	CodeExists             int = 201
	CodeCorrupted          int = 202
	CodeConflict           int = 203
	CodeAuthentication     int = 300
	CodePermissionDenied   int = 301
	CodeSystem             int = 400
	CodeConfig             int = 401
	CodeDatabase           int = 402
	CodeNetwork            int = 403
	CodeUnavailable        int = 404
	CodeNotImplemented     int = 405
	CodeUnknown            int = 999
)

// compose combines the module and error type codes into a single error code.
func compose(module, kind int) int {
	return module*10000 + kind
}

// Decompose splits a full error code into its module and error type components.
func Decompose(decompose int) (int, int) {
	return decompose / 10000, decompose % 10000
}

// NewBadRequestError creates a new BadRequest error.
func NewBadRequestError(module int, fields ...string) *QuantmError {
	return New(module, CodeBadRequest, "bad request", fields...)
}

// NewCancelledError creates a new Cancelled error.
func NewCancelledError(module int, fields ...string) *QuantmError {
	return New(module, CodeCancelled, "cancelled", fields...)
}

// NewFailedPreconditionError creates a new FailedPrecondition error.
func NewFailedPreconditionError(module int, fields ...string) *QuantmError {
	return New(module, CodeFailedPrecondition, "failed precondition", fields...)
}

// NewExhaustedError creates a new Exhausted error.
func NewExhaustedError(module int, fields ...string) *QuantmError {
	return New(module, CodeExhausted, "exhausted", fields...)
}

// NewNotFoundError creates a new NotFound error.
func NewNotFoundError(module int, fields ...string) *QuantmError {
	return New(module, CodeNotFound, "not found", fields...)
}

// NewExistsError creates a new Exists error.
func NewExistsError(module int, fields ...string) *QuantmError {
	return New(module, CodeExists, "exists", fields...)
}

// NewCorruptedError creates a new Corrupted error.
func NewCorruptedError(module int, fields ...string) *QuantmError {
	return New(module, CodeCorrupted, "corrupted", fields...)
}

// NewConflictError creates a new Conflict error.
func NewConflictError(module int, fields ...string) *QuantmError {
	return New(module, CodeConflict, "conflict", fields...)
}

// NewAuthnError creates a new Authn error.
func NewAuthnError(module int, fields ...string) *QuantmError {
	return New(module, CodeAuthentication, "authentication failed", fields...)
}

// NewAuthzError creates a new Authz error.
func NewAuthzError(module int, fields ...string) *QuantmError {
	return New(module, CodePermissionDenied, "authorization failed", fields...)
}

// NewSystemError creates a new System error.
func NewSystemError(module int, fields ...string) *QuantmError {
	return New(module, CodeSystem, "system error", fields...)
}

// NewConfigError creates a new Config error.
func NewConfigError(module int, fields ...string) *QuantmError {
	return New(module, CodeConfig, "config error", fields...)
}

// NewDatabaseError creates a new Database error.
func NewDatabaseError(module int, fields ...string) *QuantmError {
	return New(module, CodeDatabase, "database error", fields...)
}

// NewNetworkError creates a new Network error.
func NewNetworkError(module int, fields ...string) *QuantmError {
	return New(module, CodeNetwork, "network error", fields...)
}

// NewUnavailableError creates a new Unavailable error.
func NewUnavailableError(module int, fields ...string) *QuantmError {
	return New(module, CodeUnavailable, "unavailable", fields...)
}

// NewNotImplementedError creates a new NotImplemented error.
func NewNotImplementedError(module int, fields ...string) *QuantmError {
	return New(module, CodeNotImplemented, "not implemented", fields...)
}

// NewUnknownError creates a new Unknown error.
func NewUnknownError(module int, fields ...string) *QuantmError {
	return New(module, CodeUnknown, "unknown error", fields...)
}
