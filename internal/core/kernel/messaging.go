package kernel

import (
	"context"

	"github.com/google/uuid"
)

type (
	Messaging interface {
		SendMessage(ctx context.Context, to uuid.UUID, message string) error
	}
)
