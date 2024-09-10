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

package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/gocql/gocql"

	"go.breu.io/quantm/internal/core/code"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/shared"
)

func uuid() gocql.UUID {
	return gocql.TimeUUID()
}

func main() {
	slog.Info("starting merge queue ...")

	worker := shared.Temporal().Worker(shared.CoreQueue)
	worker.RegisterWorkflow(code.QueueCtrl)

	if err := worker.Start(); err != nil {
		slog.Error("unable to start worker", slog.String("error", err.Error()))
	}

	slog.Info("merge queue worker connected ...")

	// Convert string duration to time.Duration
	duration, err := time.ParseDuration("1m")
	if err != nil {
		slog.Error("error parsing duration", slog.String("error", err.Error()))
		return
	}

	// Wrap time.Duration into shared.Duration
	stale := shared.NewDuration(duration)

	repo := &defs.Repo{
		CtrlID:          uuid(),
		DefaultBranch:   "main",
		ID:              uuid(),
		IsMonorepo:      true,
		MessageProvider: "slack",
		MessageProviderData: defs.MessageProviderData{
			Slack: &defs.MessageProviderSlackData{
				BotToken:      "5Ry5/wFMD6yUenY94DXKO8zJNIUIVXs4O8YtoiPtcgyOtvTBJXTJK5RD+gObrNJm7RJlF0vrrwm+1z4ceXhQ5X06L4afVV4=",
				ChannelID:     "C06M7V3ADHV",
				ChannelName:   "#quantm-test-channel",
				WorkspaceID:   "T1U5BLPRB",
				WorkspaceName: "Breu Inc.",
			},
		},
		Name:          "quantm",
		Provider:      "github",
		ProviderID:    "506113918",
		StaleDuration: stale,
		TeamID:        uuid(),
		Threshold:     100,
	}

	branch := "test-branch"

	slog.Info("branch", "info", branch)

	opts := shared.Temporal().Queue(shared.CoreQueue).WorkflowOptions(
		shared.WithWorkflowBlock("merge_queue"),
		shared.WithWorkflowBlockID(repo.ID.String()),
	)
	ctx := context.Background()

	_, err = shared.Temporal().Client().ExecuteWorkflow(ctx, opts, code.QueueCtrl, repo, branch)
	if err != nil {
		slog.Error("workflow error", "error", err.Error())
	}

	slog.Info("execute workflow successfully")
}
