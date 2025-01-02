package access

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
	"github.com/smart-core-os/sc-golang/pkg/masks"
	"github.com/vanti-dev/sc-bos/pkg/util/pull"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/run"
)

type Group struct {
	traits.UnimplementedAccessApiServer

	client traits.AccessApiClient
	names  []string

	logger *zap.Logger
}

func (g *Group) GetLastAccessAttempt(ctx context.Context, request *traits.GetLastAccessAttemptRequest) (*traits.AccessAttempt, error) {
	fns := make([]func() (*traits.AccessAttempt, error), len(g.names))
	for i, name := range g.names {
		request := proto.Clone(request).(*traits.GetLastAccessAttemptRequest)
		request.Name = name
		fns[i] = func() (*traits.AccessAttempt, error) {
			return g.client.GetLastAccessAttempt(ctx, request)
		}
	}
	allRes, allErrs := run.Collect(ctx, run.DefaultConcurrency, fns...)

	err := multierr.Combine(allErrs...)
	if len(multierr.Errors(err)) == len(g.names) {
		return nil, err
	}

	if err != nil {
		if g.logger != nil {
			g.logger.Warn("some access implementors failed to get", zap.Errors("errors", multierr.Errors(err)))
		}
	}
	// the last access attempt is the on that happened most recently
	var last *traits.AccessAttempt
	for _, res := range allRes {
		if last == nil {
			last = res
		} else if res.AccessAttemptTime.AsTime().After(last.AccessAttemptTime.AsTime()) {
			last = res
		}
	}
	return last, nil
}

func (g *Group) PullAccessAttempts(request *traits.PullAccessAttemptsRequest, server traits.AccessApi_PullAccessAttemptsServer) error {
	if len(g.names) == 0 {
		return status.Errorf(codes.FailedPrecondition, "zone has no access implementor names")
	}

	type c struct {
		name string
		val  *traits.AccessAttempt
	}
	changes := make(chan c)
	defer close(changes)

	group, ctx := errgroup.WithContext(server.Context())

	for _, name := range g.names {
		request := proto.Clone(request).(*traits.PullAccessAttemptsRequest)
		request.Name = name
		group.Go(func() error {
			return pull.Changes(ctx, pull.NewFetcher(
				func(ctx context.Context, changes chan<- c) error {
					stream, err := g.client.PullAccessAttempts(ctx, request)
					if err != nil {
						return err
					}
					for {
						res, err := stream.Recv()
						if err != nil {
							return err
						}
						for _, change := range res.Changes {
							changes <- c{name: request.Name, val: change.AccessAttempt}
						}
					}
				},
				func(ctx context.Context, changes chan<- c) error {
					res, err := g.client.GetLastAccessAttempt(ctx, &traits.GetLastAccessAttemptRequest{Name: name, ReadMask: request.ReadMask})
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
		indexes := make(map[string]int, len(g.names))
		for i, name := range g.names {
			indexes[name] = i
		}
		values := make([]*traits.AccessAttempt, len(g.names))

		var last []*traits.AccessAttempt
		filter := masks.NewResponseFilter(masks.WithFieldMask(request.ReadMask))

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case change := <-changes:
				values[indexes[change.name]] = change.val
				if len(values) == 0 {
					return status.Errorf(codes.FailedPrecondition, "zone has no access implementor names")
				}

				for _, v := range values {
					filter.Filter(v)
				}

				if equal(last, values) {
					continue
				}

				last = values
				var accessAttemptChanges []*traits.PullAccessAttemptsResponse_Change
				for _, accessAttempt := range values {
					accessAttemptChanges = append(accessAttemptChanges, &traits.PullAccessAttemptsResponse_Change{
						Name:          change.name,
						AccessAttempt: accessAttempt,
						ChangeTime:    timestamppb.Now(),
					})
				}

				err := server.Send(&traits.PullAccessAttemptsResponse{Changes: accessAttemptChanges})
				if err != nil {
					return err
				}
			}
		}
	})
	return group.Wait()
}

func equal(as, bs []*traits.AccessAttempt) bool {

	if len(as) != len(bs) {
		return false
	}

	for i, a := range as {
		b := bs[i]
		if a == nil || b == nil {
			return false
		}

		if !(a.AccessAttemptTime.AsTime().Equal(b.AccessAttemptTime.AsTime()) &&
			a.Grant == b.Grant &&
			a.Actor == b.Actor &&
			a.Reason == b.Reason) {
			return false
		}
	}
	return true
}
