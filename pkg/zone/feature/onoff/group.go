package onoff

import (
	"context"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/zone/feature/run"
)

type Group struct {
	traits.UnimplementedOnOffApiServer
	client   traits.OnOffApiClient
	names    []string
	readOnly bool

	logger *zap.Logger
}

func (g *Group) GetOnOff(ctx context.Context, request *traits.GetOnOffRequest) (*traits.OnOff, error) {
	fns := make([]func() (*traits.OnOff, error), len(g.names))
	for i, name := range g.names {
		request := proto.Clone(request).(*traits.GetOnOffRequest)
		request.Name = name
		fns[i] = func() (*traits.OnOff, error) {
			return g.client.GetOnOff(ctx, request)
		}
	}
	allRes, allErrs := run.Collect(ctx, run.DefaultConcurrency, fns...)

	err := multierr.Combine(allErrs...)
	if len(multierr.Errors(err)) == len(g.names) {
		return nil, err
	}

	if err != nil {
		if g.logger != nil {
			g.logger.Warn("some hvacs failed to get", zap.Errors("errors", multierr.Errors(err)))
		}
	}
	return mergeOnOff(allRes)
}

func (g *Group) UpdateOnOff(ctx context.Context, request *traits.UpdateOnOffRequest) (*traits.OnOff, error) {
	if g.readOnly {
		return nil, status.Errorf(codes.FailedPrecondition, "read-only")
	}
	fns := make([]func() (*traits.OnOff, error), len(g.names))
	for i, name := range g.names {
		request := proto.Clone(request).(*traits.UpdateOnOffRequest)
		request.Name = name
		fns[i] = func() (*traits.OnOff, error) {
			return g.client.UpdateOnOff(ctx, request)
		}
	}
	allRes, allErrs := run.Collect(ctx, run.DefaultConcurrency, fns...)

	err := multierr.Combine(allErrs...)
	if len(multierr.Errors(err)) == len(g.names) {
		return nil, err
	}

	if err != nil {
		if g.logger != nil {
			g.logger.Warn("some hvacs failed to get", zap.Errors("errors", multierr.Errors(err)))
		}
	}
	return mergeOnOff(allRes)
}

func mergeOnOff(all []*traits.OnOff) (*traits.OnOff, error) {
	switch len(all) {
	case 0:
		return nil, status.Errorf(codes.FailedPrecondition, "no onoff devices in group")
	case 1:
		return all[0], nil
	default:
		state := all[0].State
		var err error
		// either they are all off, all on or it is an error
		for _, v := range all[1:] {
			if v.State != state {
				state = traits.OnOff_STATE_UNSPECIFIED
				err = multierr.Append(err, status.Errorf(codes.FailedPrecondition, "not all onoff devices have the same state"))
				break
			}
		}
		return &traits.OnOff{
			State: state,
		}, err
	}
}
