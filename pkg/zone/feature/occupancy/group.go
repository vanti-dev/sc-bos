package occupancy

import (
	"context"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/cmp"
	"github.com/vanti-dev/sc-bos/internal/util/pull"
)

type Group struct {
	traits.UnimplementedOccupancySensorApiServer
	client traits.OccupancySensorApiClient
	names  []string

	logger *zap.Logger
}

func (g *Group) GetOccupancy(ctx context.Context, request *traits.GetOccupancyRequest) (*traits.Occupancy, error) {
	var allErrs []error
	var allRes []*traits.Occupancy
	for _, name := range g.names {
		request.Name = name
		res, err := g.client.GetOccupancy(ctx, request)
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
			g.logger.Warn("some occupancy sensors failed", zap.Errors("errors", allErrs))
		}
	}
	return mergeOccupancy(allRes)
}

func (g *Group) PullOccupancy(request *traits.PullOccupancyRequest, server traits.OccupancySensorApi_PullOccupancyServer) error {
	if len(g.names) == 0 {
		return status.Error(codes.FailedPrecondition, "zone has no occupancy sensor names")
	}

	type c struct {
		name string
		val  *traits.Occupancy
	}
	changes := make(chan c)
	defer close(changes)

	group, ctx := errgroup.WithContext(server.Context())
	for _, name := range g.names {
		request := proto.Clone(request).(*traits.PullOccupancyRequest)
		request.Name = name
		group.Go(func() error {
			return pull.Changes(ctx, pull.NewFetcher(
				func(ctx context.Context, changes chan<- c) error {
					stream, err := g.client.PullOccupancy(ctx, request)
					if err != nil {
						return err
					}
					for {
						res, err := stream.Recv()
						if err != nil {
							return err
						}
						for _, change := range res.Changes {
							changes <- c{name: request.Name, val: change.Occupancy}
						}
					}
				},
				func(ctx context.Context, changes chan<- c) error {
					res, err := g.client.GetOccupancy(ctx, &traits.GetOccupancyRequest{Name: name, ReadMask: request.ReadMask})
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
		values := make([]*traits.Occupancy, len(g.names))

		var last *traits.Occupancy
		eq := cmp.Equal(cmp.FloatValueApprox(0, 0.001))

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case change := <-changes:
				values[indexes[change.name]] = change.val
				r, err := mergeOccupancy(values)
				if err != nil {
					return err
				}

				// don't send duplicates
				if eq(last, r) {
					continue
				}
				last = r

				err = server.Send(&traits.PullOccupancyResponse{Changes: []*traits.PullOccupancyResponse_Change{{
					Name:       request.Name,
					ChangeTime: timestamppb.Now(),
					Occupancy:  r,
				}}})
				if err != nil {
					return err
				}
			}
		}
	})

	return group.Wait()
}

func mergeOccupancy(all []*traits.Occupancy) (*traits.Occupancy, error) {
	switch len(all) {
	case 0:
		return nil, status.Error(codes.FailedPrecondition, "zone has no occupancy sensor names")
	case 1:
		return all[0], nil
	default:
		out := &traits.Occupancy{}
		nilCount := 0
		occupiedCount := 0
		for _, occupancy := range all {
			if occupancy == nil {
				nilCount++
				continue
			}

			out.PeopleCount += occupancy.PeopleCount

			// I don't think this logic is correct.
			// What it does is report the last change from sensors as the last change of the group.
			// What we want to do is report the last change from the sensor that caused out to change state as the last change of the group.
			// For example if sensor 1 reports occupied at 2, and sensor 2 reports unoccupied at 3, the state change time should be 2 not 3.
			if out.StateChangeTime == nil {
				out.StateChangeTime = occupancy.StateChangeTime
			} else if occupancy.StateChangeTime != nil {
				if out.StateChangeTime.AsTime().Before(occupancy.StateChangeTime.AsTime()) {
					out.StateChangeTime = occupancy.StateChangeTime
				}
			}

			switch occupancy.State {
			case traits.Occupancy_OCCUPIED:
				occupiedCount++
			}
		}
		if occupiedCount > 0 {
			out.State = traits.Occupancy_OCCUPIED
		} else {
			if len(all) > nilCount {
				out.State = traits.Occupancy_UNOCCUPIED
				out.Confidence = float64(nilCount) / float64(len(all))
			}
		}
		return out, nil
	}
}
