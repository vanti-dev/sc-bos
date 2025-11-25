package udmi

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smart-core-os/sc-bos/pkg/gen"
)

type udmiExportMessagePuller struct {
	client gen.UdmiServiceClient
	name   string
}

func (m *udmiExportMessagePuller) Pull(ctx context.Context, changes chan<- *gen.PullExportMessagesResponse) error {
	stream, err := m.client.PullExportMessages(ctx, &gen.PullExportMessagesRequest{Name: m.name})
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
		case changes <- change:
		}
	}
}

func (m *udmiExportMessagePuller) Poll(_ context.Context, _ chan<- *gen.PullExportMessagesResponse) error {
	return status.Error(codes.Unimplemented, "not supported")
}

type udmiControlTopicsPuller struct {
	client gen.UdmiServiceClient
	name   string
}

func (m *udmiControlTopicsPuller) Pull(ctx context.Context, changes chan<- *gen.PullControlTopicsResponse) error {
	stream, err := m.client.PullControlTopics(ctx, &gen.PullControlTopicsRequest{Name: m.name})
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
		case changes <- change:
		}
	}
}

func (m *udmiControlTopicsPuller) Poll(_ context.Context, _ chan<- *gen.PullControlTopicsResponse) error {
	return status.Error(codes.Unimplemented, "not supported")
}
