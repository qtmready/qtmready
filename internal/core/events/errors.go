package events

import (
	"fmt"
)

type (
	// ValidationError represents a validation error.
	ValidationError struct {
		Field   string
		Message string
	}
)

// Error implements the error interface.
func (ve ValidationError) Error() string {
	return fmt.Sprintf("validation error: field %s - %s", ve.Field, ve.Message)
}

// NewValidationError creates a new ValidationError.
func NewValidationError(field string, val any) error {
	return ValidationError{
		Field:   field,
		Message: fmt.Sprintf("invalid %s: %v", field, val),
	}
}

// NewVersionError creates a new ValidationError for the Version field.
func NewVersionError(val any) error {
	return NewValidationError("Version", val)
}

// NewIDError creates a new ValidationError for the ID field.
func NewIDError(val any) error {
	return NewValidationError("ID", val)
}

// NewActionError creates a new ValidationError for the Context.Action field.
func NewActionError(val any) error {
	return NewValidationError("Context.Action", val)
}

// NewScopeError creates a new ValidationError for the Context.Scope field.
func NewScopeError(val any) error {
	return NewValidationError("Context.Scope", val)
}

// NewSourceError creates a new ValidationError for the Context.Source field.
func NewSourceError(val any) error {
	return NewValidationError("Context.Source", val)
}

// NewProviderError creates a new ValidationError for the Context.Provider field.
func NewProviderError(val any) error {
	return NewValidationError("Context.Provider", val)
}

// NewSubjectIDError creates a new ValidationError for the Subject.ID field.
func NewSubjectIDError(val any) error {
	return NewValidationError("Subject.ID", val)
}

// NewTeamIDError creates a new ValidationError for the Subject.TeamID field.
func NewTeamIDError(val any) error {
	return NewValidationError("Subject.TeamID", val)
}

// NewTimestampError creates a new ValidationError for the Context.Timestamp field.
func NewTimestampError(val any) error {
	return NewValidationError("Context.Timestamp", val)
}

// NewSubjectNameError creates a new ValidationError for the Subject.Name field.
func NewSubjectNameError(val any) error {
	return NewValidationError("Subject.Name", val)
}
