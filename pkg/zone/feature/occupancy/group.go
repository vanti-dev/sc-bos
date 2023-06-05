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

	clients []traits.OccupancySensorApiClient // dedicated clients that don't use names for anything

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
	for _, client := range g.clients {
		res, err := client.GetOccupancy(ctx, request)
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
	if len(g.names) == 0 && len(g.clients) == 0 {
		return status.Error(codes.FailedPrecondition, "zone has no occupancy sensors")
	}

	type c struct {
		index int
		val   *traits.Occupancy
	}
	changes := make(chan c)
	defer close(changes)

	group, ctx := errgroup.WithContext(server.Context())

	// get occupancy from each of the named devices
	for i, name := range g.names {
		request := proto.Clone(request).(*traits.PullOccupancyRequest)
		request.Name = name
		index := i
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
							changes <- c{index: index, val: change.Occupancy}
						}
					}
				},
				func(ctx context.Context, changes chan<- c) error {
					res, err := g.client.GetOccupancy(ctx, &traits.GetOccupancyRequest{Name: name, ReadMask: request.ReadMask})
					if err != nil {
						return err
					}
					changes <- c{index: index, val: res}
					return nil
				}),
				changes,
			)
		})
	}

	// get occupancy from each of the dedicated clients
	for i, client := range g.clients {
		client := client
		index := len(g.names) + i
		group.Go(func() error {
			return pull.Changes(ctx, pull.NewFetcher(
				func(ctx context.Context, changes chan<- c) error {
					stream, err := client.PullOccupancy(ctx, request)
					if err != nil {
						return err
					}
					for {
						res, err := stream.Recv()
						if err != nil {
							return err
						}
						for _, change := range res.Changes {
							changes <- c{index: index, val: change.Occupancy}
						}
					}
				},
				func(ctx context.Context, changes chan<- c) error {
					res, err := client.GetOccupancy(ctx, &traits.GetOccupancyRequest{Name: request.Name, ReadMask: request.ReadMask})
					if err != nil {
						return err
					}
					changes <- c{index: index, val: res}
					return nil
				}),
				changes,
			)
		})
	}

	// merge all the changes into one occupancy and send to server
	group.Go(func() error {
		// indexes reports which index in values each name name has
		values := make([]*traits.Occupancy, len(g.names)+len(g.clients))

		var last *traits.Occupancy
		eq := cmp.Equal(cmp.FloatValueApprox(0, 0.001))

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case change := <-changes:
				values[change.index] = change.val
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
