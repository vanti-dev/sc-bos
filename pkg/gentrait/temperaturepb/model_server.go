package temperaturepb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type ModelServer struct {
	gen.UnimplementedTemperatureApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{
		model: model,
	}
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterTemperatureApiServer(server, m)
}

func (m *ModelServer) GetTemperature(_ context.Context, request *gen.GetTemperatureRequest) (*gen.Temperature, error) {
	return m.model.Temperature(resource.WithReadMask(request.ReadMask))
}

func (m *ModelServer) PullTemperature(request *gen.PullTemperatureRequest, server gen.TemperatureApi_PullTemperatureServer) error {
	for change := range m.model.PullTemperature(server.Context(), resource.WithReadMask(request.ReadMask), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		c := &gen.PullTemperatureResponse_Change{
			Name:        request.Name,
			ChangeTime:  timestamppb.New(change.ChangeTime),
			Temperature: change.Value,
		}
		if err := server.Send(&gen.PullTemperatureResponse{Changes: []*gen.PullTemperatureResponse_Change{c}}); err != nil {
			return err
		}
	}
	return nil
}

func (m *ModelServer) UpdateTemperature(_ context.Context, request *gen.UpdateTemperatureRequest) (*gen.Temperature, error) {
	return m.model.SetTemperature(
		request.Temperature,
		resource.WithUpdateMask(request.UpdateMask),
		resource.InterceptBefore(func(old, new proto.Message) {
			ot, nt := old.(*gen.Temperature), new.(*gen.Temperature)
			if request.Delta {
				if nt.SetPoint != nil {
					nt.SetPoint.ValueCelsius += ot.SetPoint.ValueCelsius
				}
				if nt.Measured != nil {
					nt.Measured.ValueCelsius += ot.Measured.ValueCelsius
				}
			}
		}),
	)
}
