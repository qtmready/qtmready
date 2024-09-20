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

package main

import (
	"fmt"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

func main() {
	shared.Service().SetName("orm")

	defer db.DB().Session.Close()

	team, _ := db.NewUUID()
	provider, _ := db.NewUUID()

	repo := &defs.Repo{
		DefaultBranch:       "main",
		TeamID:              team,
		Name:                "orm",
		IsMonorepo:          false,
		MessageProvider:     "github",
		MessageProviderData: defs.MessageProviderData{},
		Provider:            "github",
		ProviderID:          provider.String(),
		Threshold:           100,
	}

	if err := db.Save(repo); err != nil {
		shared.Logger().Error("Error saving repo", "error", err)
	}

	repo.IsMonorepo = true

	if err := db.Save(repo); err != nil {
		shared.Logger().Error("Error saving repo", "error", err)
	}

	repos := make([]defs.Repo, 0)

	if err := db.Filter(&defs.Repo{}, &repos, db.QueryParams{"is_monorepo": "true"}); err != nil {
		shared.Logger().Error("Error filter repos", "error", err)
	}

	for idx := range repos {
		repo := repos[idx]
		repo.Name = fmt.Sprintf("repo-%d", idx)

		if err := db.Save(&repo); err != nil {
			shared.Logger().Error("Error saving repo", "error", err)
		}
	}
}
