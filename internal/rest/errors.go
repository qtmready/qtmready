package rest

type (
	APIError struct {
		Code        int              `json:"status"`
		Message     string           `json:"message"`
		Information ErrorInformation `json:"information"`
	}

	ErrorInformation map[string]string

	BadRequestError   = APIError
	UnauthorizedError = APIError
	NotFoundError     = APIError
)

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(code int, message string, args ...string) *APIError {
	extra := false
	if len(args)%2 != 0 {
		extra = true
	}

	info := make(ErrorInformation)

	for i := 0; i < len(args); i += 2 {
		info[args[i]] = args[i+1]
	}

	if extra {
		info["unknown"] = args[len(args)-1]
	}

	return &APIError{
		Code:        code,
		Message:     message,
		Information: info,
	}
}

func NewBadRequestError(args ...string) *APIError {
	return NewAPIError(400, "Bad Request", args...)
}

func NewUnauthorizedError(args ...string) *APIError {
	return NewAPIError(401, "Unauthorized", args...)
}

func NewNotFoundError(args ...string) *APIError {
	return NewAPIError(404, "Not Found", args...)
}

func NewInternalServerError(args ...string) *APIError {
	return NewAPIError(500, "Internal Server Error", args...)
}
