package githubcast

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

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
			Timestamp: timestamppb.New(commit.Timestamp),
			Added:     commit.Added,
			Removed:   commit.Removed,
			Modified:  commit.Modified,
		}
	}

	return comts
}
