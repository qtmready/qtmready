// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

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
