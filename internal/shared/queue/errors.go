package queue

import (
	"errors"
	"fmt"
)

var (
	ErrParentNil = errors.New("parent workflow context is nil")
)

type (
	duplicateIDPropError struct {
		prop string
	}
)

func (e *duplicateIDPropError) Error() string {
	return fmt.Sprintf("duplicate %s", e.prop)
}

func NewDuplicateIDPropError(prop string) error {
	return &duplicateIDPropError{prop: prop}
}
