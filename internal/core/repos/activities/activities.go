package activities

import (
	"context"
)

type (
	Activity struct{}
)

func (a *Activity) SignalBranch(ctx context.Context, branch string, payload any) error { return nil }

func (a *Activity) SignalTrunk(ctx context.Context, payload any) error { return nil }
