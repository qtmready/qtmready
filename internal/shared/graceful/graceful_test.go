package graceful_test

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"go.breu.io/quantm/internal/shared/graceful"
)

type (
	Parameterized[T comparable] struct {
		name  string
		value T
		start time.Duration
		end   time.Duration
	}

	Interruptable struct {
		name string
		end  time.Duration
	}

	GracefulSuite struct {
		suite.Suite
		ctx       context.Context
		cancel    context.CancelFunc
		interrupt chan any
		errs      chan error
		terminate chan os.Signal
	}
)

func (p *Parameterized[T]) Start(arg T) error {
	time.Sleep(p.start)

	if arg != p.value {
		return errors.New("incorrect argument for Start")
	}

	return nil
}

func (p *Parameterized[T]) Stop(ctx context.Context) error {
	time.Sleep(p.end)

	return nil
}

func (i *Interruptable) Run(interrupt <-chan any) error {
	<-interrupt
	time.Sleep(i.end)

	return nil
}

func (s *GracefulSuite) SetupTest() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.interrupt = make(chan any)
	s.errs = make(chan error)
	s.terminate = make(chan os.Signal, 1)
	signal.Notify(s.terminate, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
}

func (s *GracefulSuite) TestForError() {
	// Create Parameterized instances
	fn1 := NewParameterized("string", "test1", 1*time.Second, 500*time.Millisecond)
	fn2 := NewParameterized("numbers", 123, 500*time.Millisecond, 1*time.Second)

	// Create Interruptable instances
	fn3 := NewInterruptable("interruptable1", 1*time.Second)
	fn4 := NewInterruptable("interruptable2", 500*time.Millisecond)

	// Start multiple GrabAndGo functions
	graceful.Go(s.ctx, graceful.GrabAndGo(fn1.Start, "test1"), s.errs)
	graceful.Go(s.ctx, graceful.GrabAndGo(fn2.Start, 124), s.errs)

	// Start multiple StopAndDrop functions using common interrupt channel
	graceful.Go(s.ctx, graceful.FreezeAndFizzle(fn3.Run, s.interrupt), s.errs)
	graceful.Go(s.ctx, graceful.FreezeAndFizzle(fn4.Run, s.interrupt), s.errs)

	cleanups := []graceful.Cleanup{
		fn1.Stop,
		fn2.Stop,
	}

	// Wait for error or signal
	select {
	case <-s.errs:
		// Error received, initiate shutdown
		code := graceful.Shutdown(s.ctx, cleanups, s.interrupt, 3*time.Second, 0)
		s.Equal(0, code)
	case <-s.terminate:
		// Signal received, initiate shutdown
		code := graceful.Shutdown(s.ctx, cleanups, s.interrupt, 3*time.Second, 0)
		s.Equal(0, code)
	case <-time.After(5 * time.Second):
		s.Fail("timeout waiting for error or signal")
	}
}
func (s *GracefulSuite) TestSystemInterrupt() {
	// Create Parameterized instances
	fn1 := NewParameterized("string", "test1", 1*time.Second, 500*time.Millisecond)
	fn2 := NewParameterized("numbers", 123, 500*time.Millisecond, 1*time.Second)

	// Create Interruptable instances
	fn3 := NewInterruptable("interruptable1", 1*time.Second)
	fn4 := NewInterruptable("interruptable2", 500*time.Millisecond)

	// Start multiple GrabAndGo functions
	graceful.Go(s.ctx, graceful.GrabAndGo(fn1.Start, "test1"), s.errs)
	graceful.Go(s.ctx, graceful.GrabAndGo(fn2.Start, 123), s.errs)

	// Start multiple StopAndDrop functions using common interrupt channel
	graceful.Go(s.ctx, graceful.FreezeAndFizzle(fn3.Run, s.interrupt), s.errs)
	graceful.Go(s.ctx, graceful.FreezeAndFizzle(fn4.Run, s.interrupt), s.errs)

	// Create a goroutine to send the signal
	go func() {
		s.terminate <- syscall.SIGTERM
	}()

	// Wait for signal or error
	select {
	case <-s.errs:
		// Error received, initiate shutdown
		code := graceful.Shutdown(s.ctx, []graceful.Cleanup{fn1.Stop, fn2.Stop}, s.interrupt, 3*time.Second, 0)
		s.Equal(1, code)
	case <-s.terminate:
		// Signal received, initiate shutdown
		code := graceful.Shutdown(s.ctx, []graceful.Cleanup{fn1.Stop, fn2.Stop}, s.interrupt, 3*time.Second, 0)
		s.Equal(0, code)
	case <-time.After(5 * time.Second):
		s.Fail("timeout waiting for error or signal")
	}
}

func NewParameterized[T comparable](name string, value T, start, end time.Duration) *Parameterized[T] {
	return &Parameterized[T]{name: name, value: value, start: start, end: end}
}

func NewInterruptable(name string, end time.Duration) *Interruptable {
	return &Interruptable{name: name, end: end}
}

func TestGraceful(t *testing.T) {
	suite.Run(t, new(GracefulSuite))
}
