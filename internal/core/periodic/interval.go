// Package periodic provides tools for managing recurring intervals within Temporal workflows.
// It simplifies the process of executing tasks at regular intervals, offering a more convenient
// and expressive way to handle periodic operations compared to using raw timers.
//
// Example:
//
//	// Create a new interval timer with a 5-second duration.
//	timer := periodic.New(ctx, 5*time.Second)
//
//	// Execute a single tick of the timer.
//	timer.Tick(ctx)
//
//	// Adjust the interval to 10 seconds (takes effect after current interval).
//	timer.Adjust(ctx, 10*time.Second)
//
//	// Restart with a 2-second interval (cancels current, starts new immediately).
//	timer.Restart(ctx, 2*time.Second)
//
//	// Stop the timer.
//	timer.Stop(ctx)
package periodic

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

	// Interval provides tools for managing recurring intervals within Temporal workflows.
	// It simplifies the process of executing tasks at regular intervals.
	Interval interface {
		// Tick executes a single iteration of the current interval using the current duration.
		Tick(ctx workflow.Context)

		// Adjust changes the duration of the current interval to a new value.
		// The change takes effect after the current interval completes.
		Adjust(ctx workflow.Context, duration time.Duration)

		// Reset restarts the current interval with its original duration.
		// Any changes to the interval's duration are discarded.
		Reset(ctx workflow.Context)

		// Restart immediately cancels the current interval and starts a new one with the specified duration.
		// This allows you to switch to a new frequency without waiting for the current interval to complete.
		Restart(ctx workflow.Context, duration time.Duration)

		// Stop cancels the interval and prevents any further ticks.
		Stop(ctx workflow.Context)
	}
)

func (t *interval) Adjust(ctx workflow.Context, duration time.Duration) {
	t.running = true
	t.wait(ctx)
	t.update(ctx, duration)
	t.running = false
}

func (t *interval) Restart(ctx workflow.Context, duration time.Duration) {
	if t.running {
		t.channel.Send(ctx, duration)
	} else {
		t.update(ctx, duration)
	}
}

func (t *interval) Tick(ctx workflow.Context) {
	t.Adjust(ctx, t.duration)
}

func (t *interval) Reset(ctx workflow.Context) {
	t.Restart(ctx, t.duration)
}

func (t *interval) Stop(ctx workflow.Context) {
	t.channel.Send(ctx, time.Duration(0))
}

// wait manages the execution loop of the interval, waiting for either the timer to expire or a new duration to be received on the channel.
//
// - If a new duration is received, it updates the interval's duration and resets the time until the next tick.
// - If a 0 duration is received, it stops the loop, effectively canceling the interval.
func (t *interval) wait(ctx workflow.Context) {
	done := false

	for !done && ctx.Err() == nil {
		_ctx, cancel := workflow.WithCancel(ctx)
		duration := time.Duration(0)
		timer := workflow.NewTimer(_ctx, t.duration)
		selector := workflow.NewSelector(_ctx)

		selector.AddReceive(t.channel, func(channel workflow.ReceiveChannel, more bool) {
			channel.Receive(_ctx, &duration)
			cancel()

			if duration == 0 {
				done = true
			} else {
				t.update(_ctx, t.duration)
			}
		})

		selector.AddFuture(timer, func(future workflow.Future) {
			if err := future.Get(_ctx, nil); err == nil {
				done = true
			}
		})

		selector.Select(ctx)
	}
}

func (t *interval) update(ctx workflow.Context, duration time.Duration) {
	t.duration = duration
	t.until = Now(ctx).Add(duration)
}

// Now returns the current time using a side effect.
// This is useful for obtaining the current time within a Temporal workflow.
func Now(ctx workflow.Context) time.Time {
	var now time.Time

	_ = workflow.SideEffect(ctx, func(_ctx workflow.Context) any { return time.Now() }).Get(&now)

	return now
}

// New creates a new Interval with the specified initial duration.
//
// Example:
//
//	timer := periodic.New(ctx, 5 * time.Second) // Create a new interval timer with a 5-second duration
func New(ctx workflow.Context, duration time.Duration) Interval {
	return &interval{
		duration: duration,
		until:    Now(ctx).Add(duration),
		channel:  workflow.NewChannel(ctx),
	}
}
