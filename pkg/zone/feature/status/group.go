package status

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
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
)

type Group struct {
	gen.UnimplementedStatusApiServer
	client gen.StatusApiClient
	names  []string

	logger *zap.Logger
}

func (g *Group) GetCurrentStatus(ctx context.Context, request *gen.GetCurrentStatusRequest) (*gen.StatusLog, error) {
	var allErrs []error
	var allRes []*gen.StatusLog
	for _, name := range g.names {
		request.Name = name
		res, err := g.client.GetCurrentStatus(ctx, request)
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
			// don't bother logging if all the errors are NotFound or Unimplemented
			var ignoreCount int
			for _, err := range allErrs {
				if c := status.Code(err); c != codes.NotFound && c != codes.Unimplemented {
					ignoreCount++
				}
			}
			if ignoreCount < len(allErrs) {
				g.logger.Warn("some status logs failed", zap.Errors("errors", allErrs))
			}
		}
	}
	return mergeStatusLog(allRes)
}

func (g *Group) PullCurrentStatus(request *gen.PullCurrentStatusRequest, server gen.StatusApi_PullCurrentStatusServer) error {
	if len(g.names) == 0 {
		return status.Error(codes.FailedPrecondition, "zone has no status names")
	}

	type c struct {
		name string
		val  *gen.StatusLog
	}
	changes := make(chan c)
	defer close(changes)

	group, ctx := errgroup.WithContext(server.Context())
	for _, name := range g.names {
		request := proto.Clone(request).(*gen.PullCurrentStatusRequest)
		request.Name = name
		group.Go(func() error {
			err := pull.Changes(ctx, pull.NewFetcher(
				func(ctx context.Context, changes chan<- c) error {
					stream, err := g.client.PullCurrentStatus(ctx, request)
					if err != nil {
						return err
					}
					for {
						res, err := stream.Recv()
						if err != nil {
							return err
						}
						for _, change := range res.Changes {
							changes <- c{name: request.Name, val: change.CurrentStatus}
						}
					}
				},
				func(ctx context.Context, changes chan<- c) error {
					res, err := g.client.GetCurrentStatus(ctx, &gen.GetCurrentStatusRequest{Name: name, ReadMask: request.ReadMask})
					if err != nil {
						return err
					}
					changes <- c{name: request.Name, val: res}
					return nil
				}),
				changes,
			)
			if c := status.Code(err); c == codes.NotFound || c == codes.Unimplemented {
				return nil // ignore NotFound and Unimplemented
			}
			return err
		})
	}

	group.Go(func() error {
		// indexes reports which index in values each name name has
		indexes := make(map[string]int, len(g.names))
		for i, name := range g.names {
			indexes[name] = i
		}
		values := make([]*gen.StatusLog, len(g.names))

		var last *gen.StatusLog
		eq := cmp.Equal(cmp.FloatValueApprox(0, 0.001))

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case change := <-changes:
				values[indexes[change.name]] = change.val
				r, err := mergeStatusLog(values)
				if err != nil {
					return err
				}

				// don't send duplicates
				if eq(last, r) {
					continue
				}
				last = r

				err = server.Send(&gen.PullCurrentStatusResponse{Changes: []*gen.PullCurrentStatusResponse_Change{{
					Name:          request.Name,
					ChangeTime:    timestamppb.Now(),
					CurrentStatus: r,
				}}})
				if err != nil {
					return err
				}
			}
		}
	})

	return group.Wait()
}

func mergeStatusLog(all []*gen.StatusLog) (*gen.StatusLog, error) {
	switch len(all) {
	case 0:
		return nil, status.Error(codes.FailedPrecondition, "zone has no statusLogs names")
	case 1:
		return all[0], nil
	default:
		pm := &statuspb.ProblemMerger{}
		for _, sl := range all {
			if sl == nil {
				continue
			}
			pm.AddStatusLog(sl)
		}
		return pm.Build(), nil
	}
}
