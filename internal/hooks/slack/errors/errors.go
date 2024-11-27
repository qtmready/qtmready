package errors

import (
	"errors"
)

var (
	ErrCodeEmpty   = errors.New("code is empty")
	ErrCipherText  = errors.New("ciphertext too short")
	ErrRecordExist = errors.New("record exist already")
)
