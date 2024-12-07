package cast

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func RefToProto(hook *defs.WebhookRef) eventsv1.GitRef {
	return eventsv1.GitRef{
		Ref:  hook.GetRef(),
		Kind: hook.GetRefType(),
	}
}

func PushToProto(push *defs.Push) eventsv1.Push {
	return eventsv1.Push{
		Ref:        push.GetRef(),
		Before:     push.GetBefore(),
		After:      push.GetAfter(),
		Repository: push.GetRepositoryName(),
		SenderId:   push.GetSenderID(),
		Timestamp:  timestamppb.New(time.Now()),
		Commits:    CommitsToProto(push.GetCommits()),
	}
}

func CommitsToProto(commits []defs.Commit) []*eventsv1.Commit {
	result := make([]*eventsv1.Commit, len(commits))
	for i, commit := range commits {
		result[i] = &eventsv1.Commit{
			Sha:       commit.GetID(),
			Message:   commit.GetMessage(),
			Url:       commit.GetURL(),
			Timestamp: timestamppb.New(commit.GetTimestamp()),
			Added:     commit.GetAdded(),
			Removed:   commit.GetRemoved(),
			Modified:  commit.GetModified(),
		}
	}

	return result
}

func PrToProto(pr *defs.PR) eventsv1.PullRequest {
	return eventsv1.PullRequest{
		Number:     pr.GetNumber(),
		Title:      pr.GetTitle(),
		Body:       pr.GetBody(),
		Author:     pr.GetAuthor(),
		HeadBranch: pr.GetHeadBranch(),
		BaseBranch: pr.GetBaseBranch(),
		Timestamp:  timestamppb.New(pr.GetTimestamp()),
	}
}
