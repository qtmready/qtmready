package routes

import "errors"

var (
	ErrorPasswordMismatch   = errors.New("password mismatch")
	ErrorEmailAlreadyExists = errors.New("email already exists")
)
