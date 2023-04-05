package meter

import (
	"context"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-golang/pkg/cmp"
	"github.com/vanti-dev/sc-bos/internal/util/pull"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type Group struct {
	gen.UnimplementedMeterApiServer
	client   gen.MeterApiClient
	names    []string
	readOnly bool

	logger *zap.Logger
}

func (g *Group) GetMeterReading(ctx context.Context, request *gen.GetMeterReadingRequest) (*gen.MeterReading, error) {
	var allErrs []error
	var allRes []*gen.MeterReading
	for _, name := range g.names {
		request.Name = name
		res, err := g.client.GetMeterReading(ctx, request)
		if err != nil {
			allErrs = append(allErrs, err)
			continue
		}
		allRes = append(allRes, res)
	}

	if len(allErrs) == len(g.names) {
		return nil, multierr.Combine(allErrs...)
	}

	if allErrs != nil {
		if g.logger != nil {
			g.logger.Warn("some meters failed", zap.Errors("errors", allErrs))
		}
	}
	return mergeMeterReading(allRes)
}

func (g *Group) PullMeterReadings(request *gen.PullMeterReadingsRequest, server gen.MeterApi_PullMeterReadingsServer) error {
	if len(g.names) == 0 {
		return status.Error(codes.FailedPrecondition, "zone has no meter names")
	}

	type c struct {
		name string
		val  *gen.MeterReading
	}
	changes := make(chan c)
	defer close(changes)

	group, ctx := errgroup.WithContext(server.Context())
	for _, name := range g.names {
		request := proto.Clone(request).(*gen.PullMeterReadingsRequest)
		request.Name = name
		group.Go(func() error {
			return pull.Changes(ctx, pull.NewFetcher(
				func(ctx context.Context, changes chan<- c) error {
					stream, err := g.client.PullMeterReadings(ctx, request)
					if err != nil {
						return err
					}
					for {
						res, err := stream.Recv()
						if err != nil {
							return err
						}
						for _, change := range res.Changes {
							changes <- c{name: request.Name, val: change.MeterReading}
						}
					}
				},
				func(ctx context.Context, changes chan<- c) error {
					res, err := g.client.GetMeterReading(ctx, &gen.GetMeterReadingRequest{Name: name, ReadMask: request.ReadMask})
					if err != nil {
						return err
					}
					changes <- c{name: request.Name, val: res}
					return nil
				}),
				changes,
			)
		})
	}

	group.Go(func() error {
		// indexes reports which index in values each name name has
		indexes := make(map[string]int, len(g.names))
		for i, name := range g.names {
			indexes[name] = i
		}
		values := make([]*gen.MeterReading, len(g.names))

		var last *gen.MeterReading
		eq := cmp.Equal(cmp.FloatValueApprox(0, 0.001))

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case change := <-changes:
				values[indexes[change.name]] = change.val
				r, err := mergeMeterReading(values)
				if err != nil {
					return err
				}

				// don't send duplicates
				if eq(last, r) {
					continue
				}
				last = r

				err = server.Send(&gen.PullMeterReadingsResponse{Changes: []*gen.PullMeterReadingsResponse_Change{{
					Name:         request.Name,
					ChangeTime:   timestamppb.Now(),
					MeterReading: r,
				}}})
				if err != nil {
					return err
				}
			}
		}
	})

	return group.Wait()
}

func mergeMeterReading(all []*gen.MeterReading) (*gen.MeterReading, error) {
	switch len(all) {
	case 0:
		return nil, status.Error(codes.FailedPrecondition, "zone has no meter names")
	case 1:
		return all[0], nil
	default:
		out := &gen.MeterReading{}
		// This algorithm assumes that readings before the start time are zero
		for _, reading := range all {
			// usage is summed up
			out.Usage += reading.Usage
			// use the earliest start time
			if reading.StartTime != nil {
				if out.StartTime == nil || reading.StartTime.AsTime().Before(out.StartTime.AsTime()) {
					out.StartTime = reading.StartTime
				}
			}
			// use the latest end time
			if reading.EndTime != nil {
				if out.EndTime == nil || reading.EndTime.AsTime().After(out.EndTime.AsTime()) {
					out.EndTime = reading.EndTime
				}
			}
		}
		return out, nil
	}
}
