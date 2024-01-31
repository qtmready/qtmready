package logger

import (
	"log/slog"
	"os"
	"sync"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
)

type (
	Level2Severity map[slog.Level]logging.Severity
)

var (
	project     string
	projectOnce sync.Once
)

const (
	LevelDefault   = 0
	LevelDebug     = 100
	LevelInfo      = 200
	LevelNotice    = 300
	LevelWarning   = 400
	LevelError     = 500
	LevelCritical  = 600
	LevelAlert     = 700
	LevelEmergency = 800
)

func SeverityFromLevel(level slog.Level) logging.Severity {
	translator := Level2Severity{
		slog.LevelDebug: logging.Debug,
		slog.LevelInfo:  logging.Info,
		slog.LevelWarn:  logging.Warning,
		slog.LevelError: logging.Error,
	}

	return translator[level]
}

func GoogleCloudProject() string {
	projectOnce.Do(func() {
		if metadata.OnGCE() {
			project, _ = metadata.ProjectID()
		} else {
			project = os.Getenv("GOOGLE_CLOUD_PROJECT") // This is required for local development.
		}
	})

	return project
}
