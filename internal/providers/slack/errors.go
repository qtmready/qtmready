package slack

import (
	"errors"
	"fmt"

	"github.com/slack-go/slack/socketmode"
)

type (
	socketEventPayloadError struct {
		event socketmode.EventType
	}
)

var (
	ErrOpenView            = errors.New("failed to open view")
	ErrInvalidCommand      = errors.New("invalid command")
	ErrInvalidAction       = errors.New("invalid action")
	ErrInvalidEvent        = errors.New("invalid event")
	ErrInvalidEventPayload = errors.New("invalid event payload")
)

func (e *socketEventPayloadError) Error() string {
	return fmt.Sprintf("unable to translate payload for: %s", e.event)
}

func NewSocketEventPayloadError(event socketmode.EventType) error {
	return &socketEventPayloadError{event: event}
}
