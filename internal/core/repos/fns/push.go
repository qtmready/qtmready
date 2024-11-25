package fns

import (
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func GetLatestCommit(push *eventsv1.Push) *eventsv1.Commit {
	if push == nil || len(push.Commits) == 0 {
		return nil
	}

	latest := push.Commits[0]
	timestamp := latest.Timestamp.AsTime()

	for _, commit := range push.Commits {
		_ts := commit.Timestamp.AsTime()
		if _ts.After(timestamp) {
			timestamp = _ts
			latest = commit
		}
	}

	return latest
}
