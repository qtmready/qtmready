package main

import (
	"fmt"

	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

func main() {
	shared.Service().SetName("orm")

	defer db.DB().Session.Close()

	team, _ := db.NewUUID()
	provider, _ := db.NewUUID()

	repo := &core.Repo{
		DefaultBranch:       "main",
		TeamID:              team,
		Name:                "orm",
		IsMonorepo:          false,
		MessageProvider:     "github",
		MessageProviderData: core.MessageProviderData{},
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

	repos := make([]core.Repo, 0)

	if err := db.Filter(&core.Repo{}, &repos, db.QueryParams{"is_monorepo": "true"}); err != nil {
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
