package timers

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	// interval represents a recurring timer interval.
	// The duration field specifies the length of the interval.
	// The until field specifies the time at which the interval should stop.
	interval struct {
		duration time.Duration
		until    time.Time
		channel  workflow.Channel
	}
)

// Next blocks until the the end of the interval. After that, it prepares the interval for the next iteration.
func (t *interval) Next(ctx workflow.Context) {
	t.wait(ctx)
	t.update(ctx, t.duration)
}

// NextWith blocks until the the end of the interval. After that, it prepares the interval for the next iteration
// with the specified duration.
func (t *interval) NextWith(ctx workflow.Context, duration time.Duration) {
	t.wait(ctx)
	t.update(ctx, duration)
}

// ForceUpdate stops the current interval and starts a new one with the specified duration.
func (t *interval) ForceUpdate(ctx workflow.Context, duration time.Duration) {
	t.channel.Send(ctx, duration)
}

func (t *interval) Cancel(ctx workflow.Context) {
	t.channel.Send(ctx, time.Duration(0))
}

// wait blocks until the timer expires or a message is received on the channel. The timer is cancelled if the duration is 0,
// otherwise it is reset.
func (t *interval) wait(ctx workflow.Context) {
	fired := false

	for !fired && ctx.Err() == nil {
		_ctx, cancel := workflow.WithCancel(ctx)
		duration := time.Duration(0)
		timer := workflow.NewTimer(_ctx, t.duration)
		selector := workflow.NewSelector(_ctx)

		// when the channel receives a message
		selector.AddReceive(t.channel, func(channel workflow.ReceiveChannel, more bool) {
			channel.Receive(_ctx, &duration)
			cancel()

			if duration == 0 {
				fired = true
			} else {
				t.update(_ctx, t.duration)
			}
		})

		// when the timer finishes
		selector.AddFuture(timer, func(future workflow.Future) {
			if err := future.Get(_ctx, nil); err == nil {
				fired = true
			}
		})

		selector.Select(ctx)
	}
}

// update updates the interval's duration and the time at which the interval should stop.
// The duration parameter specifies the new interval duration.
func (t *interval) update(_ workflow.Context, duration time.Duration) {
	t.duration = duration
	t.until = time.Now().Add(duration)
}

func NewInterval(ctx workflow.Context, duration time.Duration) *interval {
	return &interval{
		duration: duration,
		until:    time.Now().Add(duration),
		channel:  workflow.NewChannel(ctx),
	}
}
