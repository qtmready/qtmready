package repoacts

import (
	"context"
)

type (
	Repo struct{}
)

func (r *Repo) SignalBranch(ctx context.Context, branch string) error { return nil }
