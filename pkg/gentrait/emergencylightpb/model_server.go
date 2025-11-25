package emergencylightpb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

type ModelServer struct {
	gen.UnimplementedEmergencyLightApiServer
	model *Model
}

func NewModelServer(model *Model) *ModelServer {
	return &ModelServer{model: model}
}

func (m *ModelServer) Register(server *grpc.Server) {
	gen.RegisterEmergencyLightApiServer(server, m)
}

func (m *ModelServer) Unwrap() any {
	return m.model
}

func (m *ModelServer) GetTestResultSet(context.Context, *gen.GetTestResultSetRequest) (*gen.TestResultSet, error) {
	return m.model.GetTestResultSet(), nil
}

func (m *ModelServer) StartFunctionTest(context.Context, *gen.StartEmergencyTestRequest) (*gen.StartEmergencyTestResponse, error) {
	m.model.RunFunctionTest()
	return &gen.StartEmergencyTestResponse{}, nil
}

func (m *ModelServer) StartDurationTest(context.Context, *gen.StartEmergencyTestRequest) (*gen.StartEmergencyTestResponse, error) {
	m.model.RunDurationTest()
	return &gen.StartEmergencyTestResponse{}, nil
}

func (m *ModelServer) StopEmergencyTest(context.Context, *gen.StopEmergencyTestsRequest) (*gen.StopEmergencyTestsResponse, error) {
	// No-op for this model, as tests are run immediately
	return &gen.StopEmergencyTestsResponse{}, nil
}

func (m *ModelServer) PullTestResultSets(request *gen.PullTestResultRequest, server grpc.ServerStreamingServer[gen.PullTestResultsResponse]) error {
	for change := range m.model.PullTestResults(server.Context(), resource.WithReadMask(request.ReadMask), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		msg := &gen.PullTestResultsResponse{Changes: []*gen.PullTestResultsResponse_Change{{
			Name:       request.Name,
			ChangeTime: timestamppb.New(change.ChangeTime),
			TestResult: change.Value,
		}}}
		if err := server.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
