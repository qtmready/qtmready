// Copyright © 2023, Breu, Inc. <info@breu.io>
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

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
