package electric

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
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/merge"
)

type Group struct {
	traits.UnimplementedElectricApiServer
	client traits.ElectricApiClient
	names  []string

	logger *zap.Logger
}

func (g *Group) GetDemand(ctx context.Context, request *traits.GetDemandRequest) (*traits.ElectricDemand, error) {
	var allErrs []error
	var allRes []*traits.ElectricDemand
	for _, name := range g.names {
		request.Name = name
		res, err := g.client.GetDemand(ctx, request)
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
			g.logger.Warn("some electrics failed", zap.Errors("errors", allErrs))
		}
	}
	return mergeDemand(allRes)
}

func (g *Group) PullDemand(request *traits.PullDemandRequest, server traits.ElectricApi_PullDemandServer) error {
	if len(g.names) == 0 {
		return status.Error(codes.FailedPrecondition, "zone has no electric names")
	}

	type c struct {
		name string
		val  *traits.ElectricDemand
	}
	changes := make(chan c)
	defer close(changes)

	group, ctx := errgroup.WithContext(server.Context())
	for _, name := range g.names {
		request := proto.Clone(request).(*traits.PullDemandRequest)
		request.Name = name
		group.Go(func() error {
			return pull.Changes(ctx, pull.NewFetcher(
				func(ctx context.Context, changes chan<- c) error {
					stream, err := g.client.PullDemand(ctx, request)
					if err != nil {
						return err
					}
					for {
						res, err := stream.Recv()
						if err != nil {
							return err
						}
						for _, change := range res.Changes {
							changes <- c{name: request.Name, val: change.Demand}
						}
					}
				},
				func(ctx context.Context, changes chan<- c) error {
					res, err := g.client.GetDemand(ctx, &traits.GetDemandRequest{Name: name, ReadMask: request.ReadMask})
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
		values := make([]*traits.ElectricDemand, len(g.names))

		var last *traits.ElectricDemand
		eq := cmp.Equal(cmp.FloatValueApprox(0, 0.001))

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case change := <-changes:
				values[indexes[change.name]] = change.val
				r, err := mergeDemand(values)
				if err != nil {
					return err
				}

				// don't send duplicates
				if eq(last, r) {
					continue
				}
				last = r

				err = server.Send(&traits.PullDemandResponse{Changes: []*traits.PullDemandResponse_Change{{
					Name:       request.Name,
					ChangeTime: timestamppb.Now(),
					Demand:     r,
				}}})
				if err != nil {
					return err
				}
			}
		}
	})

	return group.Wait()
}

func mergeDemand(all []*traits.ElectricDemand) (*traits.ElectricDemand, error) {
	switch len(all) {
	case 0:
		return nil, status.Error(codes.FailedPrecondition, "zone has no electric names")
	case 1:
		return all[0], nil
	default:
		out := &traits.ElectricDemand{}
		out.Current, _ = merge.Sum(all, func(e *traits.ElectricDemand) (float32, bool) {
			if e == nil {
				return 0, false
			}
			return e.Current, true
		})
		out.Rating, _ = merge.Sum(all, func(e *traits.ElectricDemand) (float32, bool) {
			if e == nil {
				return 0, false
			}
			return e.Rating, true
		})
		// Either all the voltages are the same or we can't set out.Voltage
		for _, e := range all {
			if e == nil || e.Voltage == nil {
				continue
			}
			if out.Voltage == nil {
				out.Voltage = e.Voltage
				continue
			}
			if *out.Voltage != *e.Voltage {
				// not all voltages are equal, so we can't set
				out.Voltage = nil
				break
			}
		}
		return out, nil
	}
}
