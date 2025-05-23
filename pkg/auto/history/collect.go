package history

import (
	"bytes"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto/history/config"
	"github.com/vanti-dev/sc-bos/pkg/util/pull"
)

type collector func(context.Context, config.Source, chan<- []byte)

type getter func(context.Context, chan<- []byte) error

func collectChanges(ctx context.Context, cfg config.Source, pullFn, pollFn getter, changes chan<- []byte, logger *zap.Logger) error {
	if cfg.SparsePollingSchedule != nil {
		t := time.Now()

		payloads := make(chan []byte)
		defer close(payloads)

		var prev []byte

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(time.Until(cfg.SparsePollingSchedule.Schedule.Next(t))):
				t = time.Now()

				go func() {
					err := pollFn(ctx, payloads)
					if err != nil {
						logger.Warn("sparse poll aborted", zap.Error(err))
					}
				}()

				select {
				case <-ctx.Done():
					return nil
				case payload := <-payloads:
					if bytes.Equal(payload, prev) {
						continue
					}

					prev = payload
					changes <- payload
				}
			}
		}
	}

	err := pull.Changes(ctx, pull.NewFetcher[[]byte](pullFn, pollFn), changes, pull.WithLogger(logger))
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}
