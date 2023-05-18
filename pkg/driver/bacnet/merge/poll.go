package merge

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/task"
)

func startPoll(init context.Context, name string, pollDelay, pollTimeout time.Duration, logger *zap.Logger, pollPeer func(ctx context.Context) error) (task.StopFn, error) {
	runUntil, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(pollDelay)
	go func() {
		cleanup := func() {}
		defer func() { cleanup() }()
		for {
			cleanup()
			ctx, stop := context.WithTimeout(runUntil, pollTimeout)
			cleanup = stop
			err := pollPeer(ctx)
			LogPollError(logger, fmt.Sprintf("%s poll error", name), err)
			select {
			case <-ticker.C:
			case <-runUntil.Done():
				return
			}
		}
	}()
	return cancel, nil
}
