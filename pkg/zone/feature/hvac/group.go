package hvac

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/cmp"
	"github.com/vanti-dev/sc-bos/internal/util/pull"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Group struct {
	traits.UnimplementedAirTemperatureApiServer
	client traits.AirTemperatureApiClient
	names  []string

	logger *zap.Logger
}

func (g *Group) GetAirTemperature(ctx context.Context, request *traits.GetAirTemperatureRequest) (*traits.AirTemperature, error) {
	var allErrs []error
	var allRes []*traits.AirTemperature
	for _, name := range g.names {
		request.Name = name
		res, err := g.client.GetAirTemperature(ctx, request)
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
			g.logger.Warn("some hvacs failed", zap.Errors("errors", allErrs))
		}
	}
	return mergeAirTemperature(allRes)
}

func (g *Group) UpdateAirTemperature(ctx context.Context, request *traits.UpdateAirTemperatureRequest) (*traits.AirTemperature, error) {
	var allErrs []error
	var allRes []*traits.AirTemperature
	for _, name := range g.names {
		request.Name = name
		res, err := g.client.UpdateAirTemperature(ctx, request)
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
			g.logger.Warn("some hvacs failed", zap.Errors("errors", allErrs))
		}
	}
	return mergeAirTemperature(allRes)
}

func (g *Group) PullAirTemperature(request *traits.PullAirTemperatureRequest, server traits.AirTemperatureApi_PullAirTemperatureServer) error {
	if len(g.names) == 0 {
		return status.Error(codes.FailedPrecondition, "zone has no hvac names")
	}

	type c struct {
		name string
		val  *traits.AirTemperature
	}
	changes := make(chan c)
	defer close(changes)

	group, ctx := errgroup.WithContext(server.Context())
	for _, name := range g.names {
		request := proto.Clone(request).(*traits.PullAirTemperatureRequest)
		request.Name = name
		group.Go(func() error {
			return pull.Changes(ctx, pull.NewFetcher(
				func(ctx context.Context, changes chan<- c) error {
					stream, err := g.client.PullAirTemperature(ctx, request)
					if err != nil {
						return err
					}
					for {
						res, err := stream.Recv()
						if err != nil {
							return err
						}
						for _, change := range res.Changes {
							changes <- c{name: request.Name, val: change.AirTemperature}
						}
					}
				},
				func(ctx context.Context, changes chan<- c) error {
					res, err := g.client.GetAirTemperature(ctx, &traits.GetAirTemperatureRequest{Name: name, ReadMask: request.ReadMask})
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
		values := make([]*traits.AirTemperature, len(g.names))

		var last *traits.AirTemperature
		eq := cmp.Equal(cmp.FloatValueApprox(0, 0.001))

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case change := <-changes:
				values[indexes[change.name]] = change.val
				r, err := mergeAirTemperature(values)
				if err != nil {
					return err
				}

				// don't send duplicates
				if eq(last, r) {
					continue
				}
				last = r

				err = server.Send(&traits.PullAirTemperatureResponse{Changes: []*traits.PullAirTemperatureResponse_Change{{
					Name:           request.Name,
					ChangeTime:     timestamppb.Now(),
					AirTemperature: r,
				}}})
				if err != nil {
					return err
				}
			}
		}
	})

	return group.Wait()
}

func mergeAirTemperature(all []*traits.AirTemperature) (*traits.AirTemperature, error) {
	switch len(all) {
	case 0:
		return nil, status.Error(codes.FailedPrecondition, "zone has no hvac names")
	case 1:
		return all[0], nil
	default:
		// todo: actually merge the air temperature data
		return all[0], nil
	}
}
