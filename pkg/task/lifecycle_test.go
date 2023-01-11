package task

import (
	"context"
	"testing"
	"time"
)

func TestLifecycle_CurrentState(t *testing.T) {
	wait := 10 * time.Millisecond
	t.Run("create", func(t *testing.T) {
		lt := newLifecycleTester(t)
		lt.assertCurrentState(StatusInactive, 0)
	})
	t.Run("start", func(t *testing.T) {
		lt := newLifecycleTester(t)
		lt.startWithin(wait)
		lt.assertCurrentState(StatusActive, wait)
	})
	t.Run("stop", func(t *testing.T) {
		lt := newLifecycleTester(t)
		lt.stopWithin(wait)
		lt.assertCurrentState(StatusInactive, wait)
	})
	t.Run("start,stop", func(t *testing.T) {
		lt := newLifecycleTester(t)
		lt.startWithin(wait)
		lt.assertCurrentState(StatusActive, wait)
		lt.stopWithin(wait)
		lt.assertCurrentState(StatusInactive, wait)
	})
	t.Run("start,configure", func(t *testing.T) {
		lt := newLifecycleTester(t)
		lt.startWithin(wait)
		lt.assertCurrentState(StatusActive, wait)
		lt.applyConfigSleep(time.Millisecond)
		lt.configureWithin("foo", wait)
		lt.assertCurrentState(StatusLoading, wait)
		lt.assertCurrentState(StatusActive, wait)
	})
}

type lifecycleTester struct {
	*testing.T
	*Lifecycle[string]

	applyConfigSetup []applyConfigSetup
	applyConfigCalls []ctxConfig
}

type ctxConfig struct {
	ctx    context.Context
	config string
}

type applyConfigSetup struct {
	sleep time.Duration
}

func newLifecycleTester(t *testing.T) *lifecycleTester {
	lt := &lifecycleTester{T: t}
	lt.Lifecycle = NewLifecycle(lt.applyConfig)
	lt.ReadConfig = func(bytes []byte) (string, error) {
		return string(bytes), nil
	}
	return lt
}

func (lt *lifecycleTester) prepareApplyConfig(setup applyConfigSetup) {
	lt.applyConfigSetup = append(lt.applyConfigSetup, setup)
}

func (lt *lifecycleTester) applyConfigSleep(sleep time.Duration) {
	lt.prepareApplyConfig(applyConfigSetup{sleep: sleep})
}

func (lt *lifecycleTester) applyConfig(ctx context.Context, config string) error {
	lt.applyConfigCalls = append(lt.applyConfigCalls, ctxConfig{ctx: ctx, config: config})
	if len(lt.applyConfigSetup) > 0 {
		setup := lt.applyConfigSetup[0]
		lt.applyConfigSetup = lt.applyConfigSetup[1:]
		if setup.sleep > 0 {
			select {
			case <-time.After(setup.sleep):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return nil
}

func (lt *lifecycleTester) startWithin(wait time.Duration) {
	lt.Helper()

	ctx, stop := context.WithTimeout(context.Background(), wait)
	defer stop()
	if err := lt.Start(ctx); err != nil {
		lt.Fatalf("Start timeout after %s", wait)
	}
}

func (lt *lifecycleTester) configureWithin(config string, wait time.Duration) {
	lt.Helper()

	done := make(chan struct{})
	var err error

	go func() {
		defer close(done)
		err = lt.Configure([]byte(config))
	}()

	select {
	case <-done:
		if err != nil {
			lt.Fatalf("Configure err %v", err)
		}
		return // success
	case <-time.After(wait):
		lt.Fatalf("Configure timeout after %s", wait)
	}
}

func (lt *lifecycleTester) stopWithin(wait time.Duration) {
	lt.Helper()

	stopped := make(chan struct{})
	var stopErr error

	go func() {
		defer close(stopped)
		stopErr = lt.Stop()
	}()

	select {
	case <-stopped:
		if stopErr != nil {
			lt.Fatalf("Stop err %v", stopErr)
		}
		return // success
	case <-time.After(wait):
		lt.Fatalf("Stop timeout after %s", wait)
	}
}

func (lt *lifecycleTester) assertCurrentState(want Status, wait time.Duration) {
	lt.Helper()

	ctx, stop := context.WithTimeout(context.Background(), wait)
	defer stop()
	got := lt.Lifecycle.CurrentState()
	for got != want {
		if err := lt.Lifecycle.WaitForStateChange(ctx, got); err != nil {
			lt.Fatalf("CurrentState want %s, got timeout waiting %s", want, wait)
		}
		got = lt.Lifecycle.CurrentState()
	}

	if got != want {
		lt.Fatalf("CurrentState want %s, got %s", want, got)
	}
}

func (lt *lifecycleTester) assertNextConfig(want string) {
	if len(lt.applyConfigCalls) == 0 {
		lt.Fatalf("Expecting 1 config call, got 0")
	}
	call := lt.applyConfigCalls[0]
	lt.applyConfigCalls = lt.applyConfigCalls[1:]
	if call.config != want {
		lt.Fatalf("Config want %s, got %s", want, call.config)
	}
}
