package core

import (
	"fmt"
)

type (
	providerNotFoundError struct {
		name string
	}
)

func (e *providerNotFoundError) Error() string {
	return fmt.Sprintf("provider %s not found. plese register your providers first.", e.name)
}

func ErrProviderNotFound(name string) error {
	return &providerNotFoundError{name}
}
