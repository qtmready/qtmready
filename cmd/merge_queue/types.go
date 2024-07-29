package main

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core"
)

type (
	Queue struct {
		pull_request_id string
		repo            *core.Repo
		branches        []*string
		activties       *core.RepoActivities
		created_at      time.Time      // created_at is the time when the branch was created
		mutex           workflow.Mutex // mutex is the mutex for the state
		priority        float64
	}

	MergeQueue []*Queue

	MergeQueueWorkflows struct {
		MergeQueue MergeQueue
	}
)

func (q *Queue) next() bool {
	if len(q.branches) == 0 {
		return false
	}

	return true
}

func (q *Queue) pop(ctx workflow.Context) *string {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	if q.next() {
		branch := q.branches[0]
		q.branches = q.branches[1:]

		return branch
	}

	return nil
}

func (q *Queue) push(ctx workflow.Context, branch string) {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	q.branches = append(q.branches, &branch)
}
