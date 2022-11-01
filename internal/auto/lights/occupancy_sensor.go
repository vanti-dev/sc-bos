package lights

import (
	"context"
	"errors"
	"github.com/smart-core-os/sc-api/go/traits"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// OccupancySensorPatches contributes patches for changing the state based on occupancy sensor readings.
type OccupancySensorPatches struct {
	name   string
	client traits.OccupancySensorApiClient
	logger *zap.Logger
}

func (o *OccupancySensorPatches) Subscribe(ctx context.Context, changes chan<- Patcher) error {
	// remove our signal when we shouldn't be contributing anymore
	defer func() {
		changes <- clearOccupancyTransition(o.name)
	}()

	poll := false
	initialDelay, maxDelay := 100*time.Millisecond, 10*time.Second
	var delay time.Duration
	var errCount int

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if poll {
			return o.poll(ctx, changes)
		} else {
			err := o.pull(ctx, changes)
			if o.shouldReturn(err) {
				return err
			}
			if o.fallBackToPolling(err) {
				poll = true
				delay = 0
				errCount = 0
				continue // skip the wait
			}
			if err != nil {
				errCount++
				if errCount == 5 {
					o.logger.Warn("occupancy subscriptions are failing, will keep retrying", zap.Error(err))
				}
			} else {
				errCount = 0
				delay = 0
			}
		}

		if delay == 0 {
			delay = initialDelay
		} else {
			delay = time.Duration(float64(delay) * 1.2)
			if delay > maxDelay {
				delay = maxDelay
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
}

func (o *OccupancySensorPatches) shouldReturn(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func (o *OccupancySensorPatches) pull(ctx context.Context, changes chan<- Patcher) error {
	stream, err := o.client.PullOccupancy(ctx, &traits.PullOccupancyRequest{Name: o.name})
	if err != nil {
		return err
	}

	for {
		change, err := stream.Recv()
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case changes <- (*pullOccupancyTransition)(change):
		}
	}
}

func (o *OccupancySensorPatches) poll(ctx context.Context, changes chan<- Patcher) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		res, err := o.client.GetOccupancy(ctx, &traits.GetOccupancyRequest{Name: o.name})
		if err != nil {
			// todo: log
		} else {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case changes <- getOccupancyPatcher{o.name, res}:
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (o *OccupancySensorPatches) fallBackToPolling(err error) bool {
	if grpcError, ok := status.FromError(err); ok {
		if grpcError.Code() == codes.Unimplemented {
			return true
		}
	}
	return false
}

type pullOccupancyTransition traits.PullOccupancyResponse

func (o *pullOccupancyTransition) Patch(s *ReadState) {
	r := (*traits.PullOccupancyResponse)(o)

	for _, change := range r.Changes {
		s.Occupancy[change.Name] = change.Occupancy
	}
}

type getOccupancyPatcher struct {
	name string
	res  *traits.Occupancy
}

func (g getOccupancyPatcher) Patch(s *ReadState) {
	s.Occupancy[g.name] = g.res
}

type clearOccupancyTransition string

func (c clearOccupancyTransition) Patch(s *ReadState) {
	delete(s.Occupancy, string(c))
}
