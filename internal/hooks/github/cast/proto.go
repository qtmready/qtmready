package cast

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func RefToProto(payload *defs.WebhookRef) eventsv1.GitRef {
	return eventsv1.GitRef{
		Ref:  payload.Ref,
		Kind: payload.RefType,
	}
}

func PushToProto(payload *defs.Push) eventsv1.Push {
	return eventsv1.Push{
		Ref:        payload.Ref,
		Before:     payload.Before,
		After:      payload.After,
		Repository: payload.Repository.Name,
		SenderId:   payload.Sender.ID,
		Timestamp:  timestamppb.New(time.Now()),
		Commits:    CommitsToProto(payload.Commits),
	}
}

func CommitsToProto(commits []defs.Commit) []*eventsv1.Commit {
	result := make([]*eventsv1.Commit, len(commits))
	for i, commit := range commits {
		result[i] = &eventsv1.Commit{
			Sha:       commit.ID,
			Message:   commit.Message,
			Url:       commit.URL,
			Timestamp: timestamppb.New(commit.Timestamp.Time()),
			Added:     commit.Added,
			Removed:   commit.Removed,
			Modified:  commit.Modified,
		}
	}

	return result
}
