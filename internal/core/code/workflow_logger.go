// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
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

package code

import (
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
)

type (
	RepoIOWorkflowLogger struct {
		repo   *defs.Repo
		kind   string
		branch string
		action string
		logger log.Logger
	}

	LogWriter func(msg string, keyvals ...any)
)

func NewRepoIOWorkflowLogger(ctx workflow.Context, repo *defs.Repo, kind, branch, action string) *RepoIOWorkflowLogger {
	logger := workflow.GetLogger(ctx)

	return &RepoIOWorkflowLogger{repo, kind, branch, action, logger}
}

func (r *RepoIOWorkflowLogger) Info(msg string, keyvals ...any) {
	r.write(r.logger.Info, msg, keyvals...)
}

func (r *RepoIOWorkflowLogger) Warn(msg string, keyvals ...any) {
	r.write(r.logger.Warn, msg, keyvals...)
}

func (r *RepoIOWorkflowLogger) Error(msg string, keyvals ...any) {
	r.write(r.logger.Error, msg, keyvals...)
}

func (r *RepoIOWorkflowLogger) Debug(msg string, keyvals ...any) {
	r.write(r.logger.Debug, msg, keyvals...)
}

func (r *RepoIOWorkflowLogger) prefix() string {
	prefix := r.kind

	if r.branch != "" {
		prefix += "/" + r.branch
	}

	if r.action != "" {
		prefix += "/" + r.action
	}

	return prefix + ": "
}

func (r *RepoIOWorkflowLogger) write(writer LogWriter, msg string, keyvals ...any) {
	keyvals = append(keyvals, "repo_id", r.repo.ID.String())
	keyvals = append(keyvals, "provider", r.repo.Provider.String())
	keyvals = append(keyvals, "provider_id", r.repo.ProviderID)
	keyvals = append(keyvals, "default_branch", r.repo.DefaultBranch)

	if r.branch != "" {
		keyvals = append(keyvals, "branch", r.branch)
	}

	if r.action != "" {
		keyvals = append(keyvals, "scope", r.action)
	}

	writer(r.prefix()+msg, keyvals...)
}
