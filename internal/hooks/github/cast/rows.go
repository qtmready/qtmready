package githubcast

import (
	"encoding/json"
	"fmt"

	reposcast "go.breu.io/quantm/internal/core/repos/cast"
	reposdefs "go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
)

func RowToFullRepo(row entities.GetRepoRow) (*reposdefs.FullRepo, error) {
	// Unmarshal the Messaging field
	messaging := entities.Messaging{}
	if len(row.Messaging) > 0 {
		if err := json.Unmarshal(row.Messaging, &messaging); err != nil {
			return nil, fmt.Errorf("failed to unmarshal messaging: %w", err)
		}
	}

	// Unmarshal the Org field
	org := entities.Org{}
	if len(row.Org) > 0 {
		if err := json.Unmarshal(row.Org, &org); err != nil {
			return nil, fmt.Errorf("failed to unmarshal org: %w", err)
		}
	}

	core := &reposdefs.FullRepo{
		ID:            row.ID,
		OrgID:         row.OrgID,
		Name:          row.Name,
		Hook:          reposcast.HookToProto(row.Hook),
		HookID:        row.HookID,
		DefaultBranch: row.DefaultBranch,
		IsMonorepo:    row.IsMonorepo,
		Threshold:     row.Threshold,
		StaleDuration: db.IntervalToDuration(row.StaleDuration),
		Url:           row.Url,
		IsActive:      row.IsActive,
		User:          &messaging,
		Org:           &org,
	}

	return core, nil
}
