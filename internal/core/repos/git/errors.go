package git

import (
	"errors"
)

var (
	ErrRepoAlreadyExists = errors.New("repository already exists")
	ErrInvalidBranch     = errors.New("invalid branch")
	ErrTokenization      = errors.New("tokenization error")
	ErrClone             = errors.New("clone error")
	ErrOpen              = errors.New("open error")
	ErrNoCommonAncestor  = errors.New("no common ancestor found")
)
