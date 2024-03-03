package events

import (
	"log/slog"
	"time"

	"github.com/gocql/gocql"
)

type (
	// Provider can be core.RepoProvider or core.CloudProvider.
	Provider interface {
		String() string
	}

	MetaData map[string]string

	Event struct {
		ID           gocql.UUID `cql:"id"`
		TeamID       gocql.UUID `cql:"team_id"`
		Level        slog.Level `cql:"level"`
		Name         string     `cql:"name"`
		Type         string     `cql:"type"`
		Provider     Provider   `cql:"provider"`
		ProviderType string     `cql:"type"`
		MetaData     MetaData   `cql:"metadata"`
		CreatedAt    time.Time  `cql:"created_at"`
		UpdatedAt    time.Time  `cql:"updated_at"`
	}

	Option func(*Event)
)

func (e *Event) Save() error {
	return nil
}

func WithName(name string) Option {
	return func(e *Event) {
		e.Name = name
	}
}

func WithType(t string) Option {
	return func(e *Event) {
		e.Type = t
	}
}

func WithProvider(p Provider) Option {
	return func(e *Event) {
		e.Provider = p
	}
}

func WithProviderType(t string) Option {
	return func(e *Event) {
		e.ProviderType = t
	}
}

func WithMetaData(key, value string) Option {
	return func(e *Event) {
		e.MetaData[key] = value
	}
}

// New creates a new event with the given options.
//
//	e := events.New(
//	  events.WithName("push-event"),
//	  events.WithType("provider"),
//	  events.WithProvider(core.RepoProvider),
//	  events.WithProviderType("repo"),
//	  events.WithMetaData("id", "12345"),
//	  events.WithMetaData("repo", ""),
//	)
func New(opts ...Option) *Event {
	e := &Event{
		Level:     slog.LevelInfo,
		MetaData:  MetaData{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}
