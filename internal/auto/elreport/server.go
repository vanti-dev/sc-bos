package elreport

import (
	"bytes"
	"context"
	"errors"
	"strconv"

	"github.com/timshannon/bolthold"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const maxPageSize = 100

var (
	errLightNotFound = status.Error(codes.NotFound, "EmergencyLight not found")
	errDatabase      = status.Error(codes.Unavailable, "database query failed")
)

type server struct {
	gen.UnimplementedEmergencyLightingApiServer

	db     *bolthold.Store
	logger *zap.Logger
}

func (s *server) register(registrar grpc.ServiceRegistrar) {
	gen.RegisterEmergencyLightingApiServer(registrar, s)
}

func (s *server) GetEmergencyLight(ctx context.Context, request *gen.GetEmergencyLightRequest) (*gen.EmergencyLight, error) {
	var record LatestStatusRecord
	err := s.db.Get(request.Name, &record)
	if errors.Is(err, bolthold.ErrNotFound) {
		return nil, errLightNotFound
	} else if err != nil {
		s.logger.Error("GetEmergencyLight: database query for LatestStatusRecord failed", zap.Error(err),
			zap.String("deviceName", request.Name))
		return nil, errDatabase
	}

	return translateStatusRecord(record), nil
}

func (s *server) ListEmergencyLights(ctx context.Context, request *gen.ListEmergencyLightsRequest) (*gen.ListEmergencyLightsResponse, error) {
	page, nextToken, err := findLatestStatusPaged(s.db, request.PageToken, int(request.PageSize))
	if err != nil {
		s.logger.Error("find LatestStatusRecord failure", zap.Error(err))
		return nil, errDatabase
	}

	res := &gen.ListEmergencyLightsResponse{
		NextPageToken: nextToken,
	}
	for _, record := range page {
		res.EmergencyLights = append(res.EmergencyLights, translateStatusRecord(record))
	}
	return res, nil
}

func (s *server) ListEmergencyLightEvents(ctx context.Context, request *gen.ListEmergencyLightEventsRequest) (*gen.ListEmergencyLightEventsResponse, error) {
	page, nextToken, err := findEventsPaged(s.db, request.PageToken, int(request.PageSize))
	if err != nil {
		s.logger.Error("find EventRecord failure", zap.Error(err))
		return nil, errDatabase
	}

	res := &gen.ListEmergencyLightEventsResponse{
		NextPageToken: nextToken,
	}
	for _, record := range page {
		res.Events = append(res.Events, translateEventRecord(record))
	}
	return res, nil
}

func (s *server) GetReportCSV(ctx context.Context, request *gen.GetReportCSVRequest) (*gen.ReportCSV, error) {
	report, err := GenerateReport(s.db)
	if err != nil {
		s.logger.Error("failed to generate report", zap.Error(err))
		return nil, errDatabase
	}

	var buf bytes.Buffer
	err = WriteReportCSV(&buf, report, request.IncludeHeader)
	if err != nil {
		s.logger.Error("failed to write report as a CSV", zap.Error(err), zap.Int("count", len(report)))
		return nil, status.Error(codes.Internal, "failed to convert report to CSV")
	}

	return &gen.ReportCSV{
		Csv: buf.Bytes(),
	}, nil
}

func translateStatusRecord(record LatestStatusRecord) *gen.EmergencyLight {
	return &gen.EmergencyLight{
		Name:       record.Name,
		UpdateTime: timestamppb.New(record.LastUpdate),
		Faults:     record.Faults,
	}
}

func translateEventRecord(record EventRecord) *gen.EmergencyLightingEvent {
	event := &gen.EmergencyLightingEvent{
		Name:      record.Name,
		Id:        strconv.FormatUint(record.ID, 10),
		Timestamp: timestamppb.New(record.Timestamp),
	}
	switch record.Kind {
	case StatusReportEvent:
		event.Event = &gen.EmergencyLightingEvent_StatusReport_{
			StatusReport: record.StatusReport,
		}
	case FunctionTestPassEvent:
		event.Event = &gen.EmergencyLightingEvent_FunctionTestPass_{
			FunctionTestPass: &gen.EmergencyLightingEvent_FunctionTestPass{},
		}
	case DurationTestPassEvent:
		event.Event = &gen.EmergencyLightingEvent_DurationTestPass_{
			DurationTestPass: record.DurationTestPass,
		}
	}
	return event
}
