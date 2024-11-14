package githubcast

import (
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	reposdefs "go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func PushToProto(payload *githubdefs.Push) *eventsv1.Push {
	return &eventsv1.Push{
		Ref:        payload.Ref,
		Before:     payload.Before,
		After:      payload.After,
		Repository: payload.Repository.Name,
		SenderId:   payload.SenderID(),
		Timestamp:  timestamppb.New(time.Now()),
		Commits:    CommitsToProto(payload.Commits),
	}
}

func CommitsToProto(commits []githubdefs.Commit) []*eventsv1.Commit {
	comts := make([]*eventsv1.Commit, len(commits))
	for i, commit := range commits {
		comts[i] = &eventsv1.Commit{
			Sha:       commit.ID,
			Message:   commit.Message,
			Url:       commit.URL,
			Timestamp: timestamppb.New(commit.Timestamp.Time()),
			Added:     commit.Added,
			Removed:   commit.Removed,
			Modified:  commit.Modified,
		}
	}

	return comts
}

func ConvertGetRepoRowToCoreRepo(row entities.GetRepoRow) (*reposdefs.CoreRepo, error) {
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

	core := &reposdefs.CoreRepo{
		ID:            row.ID,
		OrgID:         row.OrgID,
		Name:          row.Name,
		Hook:          row.Hook,
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
