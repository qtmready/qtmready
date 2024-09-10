// Copyright © 2024, Breu, Inc. <info@breu.io>
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package timers

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	interval struct {
		running  bool
		duration time.Duration
		until    time.Time
		channel  workflow.Channel
	}

	// Interval provides helpers to manage a recurring interval in the context of a temporal workflow.
	Interval interface {
		// Next blocks until the the end of the interval. After that, it prepares the interval for the next iteration.
		Next(ctx workflow.Context)

		// NextWith blocks until the the end of the interval. After that, it prepares the interval for the next iteration
		// with the specified duration.
		NextWith(ctx workflow.Context, duration time.Duration)

		// Restart restarts the current interval.
		Restart(ctx workflow.Context)

		// RestartWith restarts the current interval with the specified duration.
		RestartWith(ctx workflow.Context, duration time.Duration)

		// Cancel stops the current interval.
		Cancel(ctx workflow.Context)
	}
)

// NextWith blocks until the the end of the interval. After that, it prepares the interval for the next iteration
// with the specified duration.
func (t *interval) NextWith(ctx workflow.Context, duration time.Duration) {
	t.running = true
	t.wait(ctx)
	t.update(ctx, duration)
	t.running = false
}

// RestartWith stops the current interval and starts a new one with the specified duration.
func (t *interval) RestartWith(ctx workflow.Context, duration time.Duration) {
	if t.running {
		t.channel.Send(ctx, duration)
	} else {
		t.update(ctx, duration)
	}
}

// Next blocks until the the end of the interval. After that, it prepares the interval for the next iteration.
func (t *interval) Next(ctx workflow.Context) {
	t.NextWith(ctx, t.duration)
}

// Restart restarts the current interval.
func (t *interval) Restart(ctx workflow.Context) {
	t.RestartWith(ctx, t.duration)
}

// Cancel stops the current interval by sending a 0 duration to the channel, which will cause the interval to stop.
func (t *interval) Cancel(ctx workflow.Context) {
	t.channel.Send(ctx, time.Duration(0))
}

// wait blocks until the end of the interval or until a new duration is received on the channel.
// If a new duration is received, it updates the interval's duration and the time at which the interval should stop.
// If a 0 duration is received, it sets the done flag to true to stop the interval.
func (t *interval) wait(ctx workflow.Context) {
	done := false

	for !done && ctx.Err() == nil {
		_ctx, cancel := workflow.WithCancel(ctx)
		duration := time.Duration(0)
		timer := workflow.NewTimer(_ctx, t.duration)
		selector := workflow.NewSelector(_ctx)

		// when the channel receives a message
		selector.AddReceive(t.channel, func(channel workflow.ReceiveChannel, more bool) {
			channel.Receive(_ctx, &duration)
			cancel()

			if duration == 0 {
				done = true // we recieved a 0 duration, so we should stop the interval.
			} else {
				t.update(_ctx, t.duration)
			}
		})

		// when the timer finishes
		selector.AddFuture(timer, func(future workflow.Future) {
			if err := future.Get(_ctx, nil); err == nil {
				done = true
			}
		})

		selector.Select(ctx)
	}
}

// update updates the interval's duration and the time at which the interval should stop.
func (t *interval) update(ctx workflow.Context, duration time.Duration) {
	t.duration = duration
	t.until = Now(ctx).Add(duration)
}

func NewInterval(ctx workflow.Context, duration time.Duration) Interval {
	return &interval{
		duration: duration,
		until:    Now(ctx).Add(duration),
		channel:  workflow.NewChannel(ctx),
	}
}
