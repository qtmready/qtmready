package shared

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrValidation          = errors.New("validation error")
)

type (
	ErrorMap map[string]string
)

// NewAPIError replaces echo.NewHTTPError.
func NewAPIError(code int, message error) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// SetInternal sets the internal error.
func (e *APIError) SetInternal(err error) {
	e.Internal = err
}

func (e *APIError) WithInternal(err error) *APIError {
	return &APIError{
		Code:     e.Code,
		Message:  e.Message,
		Internal: err,
	}
}

func (e *APIError) Normalize() *APIError {
	e.Message = e.Message.(error).Error()
	if e.Internal != nil && e.Errors == nil {
		errs := ErrorMap{}
		errs["internal"] = e.Internal.Error()
		e.Errors = &errs
	}
	return e
}

func (e *APIError) Unwrap() error {
	return e.Internal
}

// APIErrorHandler adds syntax sugar to the default echo error handler.
func APIErrorHandler(err error, ctx echo.Context) {
	var apierr *APIError
	if ctx.Response().Committed {
		return
	}

	// We create an APIError from the error if it is not already one.
	apierr, aok := err.(*APIError)
	if !aok {
		apierr = NewAPIError(http.StatusInternalServerError, ErrInternalServerError)
		apierr.WithInternal(err)
	}

	// Now we check if the internal error is a valiator error.
	ve, vok := apierr.Message.(validator.ValidationErrors)
	if vok {
		errs := ErrorMap{}
		for _, fe := range ve {
			errs[fe.Field()] = TagMessage(fe.Tag())
		}
		// We set the error map to the APIError and set the error to ErrValidation.
		apierr = NewAPIError(apierr.Code, ErrValidation)
		apierr.WithInternal(ve)
		apierr.Errors = &errs
	}

	// We set the status code and return the error.
	if ctx.Request().Method == http.MethodHead {
		// Work around for echo issue #1327
		err = ctx.NoContent(apierr.Code)
	} else {
		err = ctx.JSON(apierr.Code, apierr.Normalize())
	}

	if err != nil {
		ctx.Logger().Error(err)
	}
}

func TagMessage(tag string) string {
	switch tag {
	case "required":
		return "required"
	case "db_unique":
		return "already exists"
	default:
		return fmt.Sprintf("%s, validation error", tag)
	}
}
