package erratic

const (
	StatusValidationError = 100001
)

func NewValidationError(args ...string) *QuantmError {
	return New(StatusValidationError, "Validation Error.", args...)
}
