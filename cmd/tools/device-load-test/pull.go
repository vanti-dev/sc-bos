package main

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/accesspb"
)

type pullResult struct {
	deviceTrait
	Time time.Time
	Data proto.Message
	Err  error
}

func pull(ctx context.Context, dst chan<- pullResult, conn grpc.ClientConnInterface, source deviceTrait) {
	switch source.Trait {
	case trait.AirTemperature:
		pullAirTemperature(ctx, dst, conn, source)
	case trait.Light:
		pullLight(ctx, dst, conn, source)
	case trait.OpenClose:
		pullOpenClose(ctx, dst, conn, source)
	case accesspb.TraitName:
		pullAccess(ctx, dst, conn, source)
	default:
		log.Printf("unknown trait type: %s", source.Trait)
	}
}

func pullAirTemperature(ctx context.Context, dst chan<- pullResult, conn grpc.ClientConnInterface, source deviceTrait) {
	client := traits.NewAirTemperatureApiClient(conn)
	stream, err := client.PullAirTemperature(ctx, &traits.PullAirTemperatureRequest{
		Name: source.Name,
	})
	if err != nil {
		failPull(ctx, dst, source, err)
		return
	}

	pullStream(ctx, dst, func() (proto.Message, error) {
		return stream.Recv()
	}, source)
}

func pullLight(ctx context.Context, dst chan<- pullResult, conn grpc.ClientConnInterface, source deviceTrait) {
	client := traits.NewLightApiClient(conn)
	stream, err := client.PullBrightness(ctx, &traits.PullBrightnessRequest{
		Name: source.Name,
	})
	if err != nil {
		failPull(ctx, dst, source, err)
		return
	}
	pullStream(ctx, dst, func() (proto.Message, error) {
		return stream.Recv()
	}, source)
}

func pullOpenClose(ctx context.Context, dst chan<- pullResult, conn grpc.ClientConnInterface, source deviceTrait) {
	client := traits.NewOpenCloseApiClient(conn)
	stream, err := client.PullPositions(ctx, &traits.PullOpenClosePositionsRequest{
		Name: source.Name,
	})
	if err != nil {
		failPull(ctx, dst, source, err)
		return
	}

	pullStream(ctx, dst, func() (proto.Message, error) {
		return stream.Recv()
	}, source)
}

func pullAccess(ctx context.Context, dst chan<- pullResult, conn grpc.ClientConnInterface, source deviceTrait) {
	client := gen.NewAccessApiClient(conn)
	stream, err := client.PullAccessAttempts(ctx, &gen.PullAccessAttemptsRequest{
		Name: source.Name,
	})
	if err != nil {
		failPull(ctx, dst, source, err)
		return
	}
	pullStream(ctx, dst, func() (proto.Message, error) {
		return stream.Recv()
	}, source)
}

func pullStream(ctx context.Context, dst chan<- pullResult, recv func() (proto.Message, error), source deviceTrait) {
	for {
		data, err := recv()
		if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		} else if err != nil {
			msg := pullResult{
				deviceTrait: source,
				Err:         err,
				Time:        time.Now(),
			}
			select {
			case dst <- msg:
			case <-ctx.Done():
			}
			return
		}

		msg := pullResult{
			deviceTrait: source,
			Data:        data,
			Time:        time.Now(),
		}
		select {
		case dst <- msg:
		case <-ctx.Done():
			return
		}
	}
}

func failPull(ctx context.Context, dst chan<- pullResult, source deviceTrait, err error) {
	msg := pullResult{
		deviceTrait: source,
		Err:         err,
		Time:        time.Now(),
	}
	select {
	case dst <- msg:
	case <-ctx.Done():
	}
}
