package kernel

import (
	"context"

	"github.com/google/uuid"
)

type (
	Chat interface {
		SendMessage(ctx context.Context, to uuid.UUID, message string) error
	}
)
