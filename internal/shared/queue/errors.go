package queue

import (
	"errors"
	"fmt"
)

var (
	ErrParentNil = errors.New("parent workflow context is nil")
)

type (
	duplicateIdPropError struct {
		prop string
	}
)

func (e *duplicateIdPropError) Error() string {
	return fmt.Sprintf("duplicate %s", e.prop)
}

func NewDuplicateIdPropError(prop string) error {
	return &duplicateIdPropError{prop: prop}
}
