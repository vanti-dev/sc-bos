package udmi

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type udmiMessagePuller struct {
	client gen.UdmiServiceClient
	name   string
}

func (m *udmiMessagePuller) Pull(ctx context.Context, changes chan<- *gen.PullExportMessagesResponse) error {
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

func (m *udmiMessagePuller) Poll(_ context.Context, _ chan<- *gen.PullExportMessagesResponse) error {
	return status.Error(codes.Unimplemented, "not supported")
}
