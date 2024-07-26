// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package core

import (
	"sort"
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/timers"
	"go.breu.io/quantm/internal/shared"
)

// NOTE - It is POC.
type (
	Signal struct {
		merge_queue_signal *shared.MergeQueueSignal
		created_at         time.Time       // created_at is the time when the branch was created
		interval           timers.Interval // interval is the interval at which the branch is checked for staleness
		mutex              workflow.Mutex  // mutex is the mutex for the state
		counter            int             // counter is the number of steps taken by the branch
		priority           float64
	}

	MergeQueue []*Signal

	MergeQueueWorkflows struct {
		MergeQueue MergeQueue
	}
)

// set_created_at sets the created_at timestamp for the merge queue signal.
// This method is thread-safe and locks the merge queue signal's mutex before updating the created_at field.
func (signal *Signal) set_created_at(ctx workflow.Context, t time.Time) {
	_ = signal.mutex.Lock(ctx)
	defer signal.mutex.Unlock()

	signal.created_at = t
}

// Workflow method for managing the merge queue.
func (w *MergeQueueWorkflows) MergeQueueWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Merge Queue Workflow started.")

	// Listen for signals to add tasks to the queue
	// Understaind the logic
	for {
		var signal Signal

		workflow.GetSignalChannel(ctx, "merge_queue_signal").Receive(ctx, &signal)

		// Add the signal to the queue
		signal.priority = w.priority(signal)
		w.MergeQueue = append(w.MergeQueue, &signal)

		w.sort(ctx)
	}
}

// TEST Workflow may not needed.
func (w *MergeQueueWorkflows) priority(signal Signal) float64 {
	age := time.Since(signal.created_at).Seconds()
	return 1.0 / (1.0 + age) // Example: simple inverse age
}

// TEST Workflow may not needed.
func (w *MergeQueueWorkflows) sort(ctx workflow.Context) {
	logger := workflow.GetLogger(ctx)

	// SliceStable sorts the slice x using the provided less function, keeping equal elements in their original order.
	// It panics if x is not a slice.
	sort.SliceStable(w.MergeQueue, func(i, j int) bool {
		return w.MergeQueue[i].priority < w.MergeQueue[j].priority
	})

	mqa := MergeQueueActivities{}
	opts := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
	ctx = workflow.WithActivityOptions(ctx, opts)

	// Process each task in the queue
	for len(w.MergeQueue) > 0 {
		task := w.MergeQueue[0]
		w.MergeQueue = w.MergeQueue[1:]

		err := workflow.ExecuteActivity(ctx, mqa.ProcessSignalActivity, task).Get(ctx, nil)
		if err != nil {
			logger.Error("Error executing activity", "error", err)
		}
	}
}

// func (mq MergeQueue) Len() int { return len(mq) }

// func (mq MergeQueue) Less(i, j int) bool {
// 	return mq[i].priority < mq[j].priority
// }

// func (mq MergeQueue) Swap(i, j int) {
// 	mq[i], mq[j] = mq[j], mq[i]
// }

// func (mq *MergeQueue) Push(x *Singal) {
// 	*mq = append(*mq, x)
// }

// func (mq *MergeQueue) Pop() *Singal {
// 	old := *mq
// 	n := len(old)
// 	item := old[n-1]
// 	*mq = old[0 : n-1]

// 	return item
// }
