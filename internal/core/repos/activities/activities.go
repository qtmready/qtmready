package activities

import (
	"context"
)

type (
	Repo struct{}
)

func (r *Repo) SignalBranch(ctx context.Context, branch string, payload any) error { return nil }

func (r *Repo) SignalTrunk(ctx context.Context, payload any) error { return nil }
