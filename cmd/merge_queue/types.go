package main

import (
	"time"
)

type (
	MergeQueueSignal struct {
		PullRequestID  int64
		InstallationID int64
		RepoOwner      string
		RepoName       string
		Branch         string
		RepoProvider   string
	}

	RepoActivities struct {
		// Define RepoActivities fields as required
	}

	Interval struct {
		// Define Interval fields as required
	}

	Signal struct {
		merge_queue_signal *MergeQueueSignal
		activities         *RepoActivities
		created_at         time.Time
		interval           Interval
		counter            int
		priority           float64
	}

	MergeQueue []*Signal

	MergeQueueWorkflows struct {
		MergeQueue MergeQueue
	}
)
