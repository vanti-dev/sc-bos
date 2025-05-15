package devicemonitor

import (
	"context"
	"errors"
	"sort"
	"testing"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperaturepb"
	"github.com/vanti-dev/sc-bos/pkg/auto/devicemonitor/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

func TestCheckReturnTemperaturesAreOk(t *testing.T) {
	fcuName := "fcu1"
	devices := []*config.Device{{
		Name: fcuName,
	}}
	nodeName := "node1"
	model := airtemperaturepb.NewModel()
	n := node.New(nodeName)

	// todo I think this needs to change to be an alert thing with an in mem storage
	n.Announce(fcuName,
		node.HasTrait(trait.AirTemperature),
		node.HasServer(traits.RegisterAirTemperatureApiServer, traits.AirTemperatureApiServer(airtemperaturepb.NewModelServer(model))),
	)
	ctx, stop := context.WithCancel(context.Background())
	t.Cleanup(stop)

	alertClient := &memalert{}
	alertClient.alerts = make(map[string]*gen.Alert)
	auto := &deviceMonitorAuto{
		alertAdminClient: alertClient,
	}
	auto.Logger, _ = zap.NewDevelopment()

	tempClient := traits.NewAirTemperatureApiClient(n.ClientConn())

	tests := []struct {
		name       string
		okLower    float64
		okUpper    float64
		temp       float64
		isResolved bool
		want       *gen.Alert
	}{
		{"Test with normal temperature", 10, 30, 20, false, nil},
		{"Test with low temperature", 10, 30, 5, false, &gen.Alert{
			Id:          "low",
			Source:      fcuName,
			Severity:    gen.Alert_WARNING,
			Description: "Ambient temperature is abnormally low",
		}},
		{"Test with temperature returning within bounds", 10, 30, 15, true, &gen.Alert{
			Id:          "low",
			Source:      fcuName,
			Severity:    gen.Alert_WARNING,
			Description: "Ambient temperature is abnormally low",
		}},
		{"Test with high temperature", 10, 30, 35, false, &gen.Alert{
			Id:          "high",
			Source:      fcuName,
			Severity:    gen.Alert_WARNING,
			Description: "Ambient temperature is abnormally high",
		}},
		{"Test with temperature returning within bounds", 10, 30, 25, true, &gen.Alert{
			Id:          "high",
			Source:      fcuName,
			Severity:    gen.Alert_WARNING,
			Description: "Ambient temperature is abnormally high",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			model.UpdateAirTemperature(&traits.AirTemperature{
				AmbientTemperature: &types.Temperature{
					ValueCelsius: tt.temp,
				},
			})
			auto.checkReturnTemperaturesAreNormal(ctx, tempClient, &config.AirTempConfig{
				Devices:        devices,
				OkRtLowerBound: proto.Float64(tt.okLower),
				OkRtUpperBound: proto.Float64(tt.okUpper),
			})

			s, err := alertClient.ListAlerts(ctx, &gen.ListAlertsRequest{
				Name: fcuName,
			})
			if err != nil {
				t.Fatal(err)
			}
			if tt.want == nil {
				if len(s.Alerts) != 0 {
					t.Fatalf("expected no problems, got %d", len(s.Alerts))
				}
			} else {
				// check all fields of tt want match s
				if len(s.Alerts) == 0 {
					t.Fatalf("expected problems got 0")
				}
				last := len(s.Alerts) - 1
				if s.Alerts[last].Source != tt.want.Source {
					t.Fatalf("expected problem source %s, got %s", tt.want.Source, s.Alerts[last].Source)
				}
				if s.Alerts[last].Severity != tt.want.Severity {
					t.Fatalf("expected problem severity %s, got %s", tt.want.Severity, s.Alerts[last].Severity)
				}
				if s.Alerts[last].Description != tt.want.Description {
					t.Fatalf("expected problem description %s, got %s", tt.want.Description, s.Alerts[last].Description)
				}
				if tt.isResolved {
					if s.Alerts[last].ResolveTime == nil {
						t.Fatalf("expected problem to be resolved")
					}
				} else {
					if s.Alerts[last].ResolveTime != nil {
						t.Fatalf("expected problem to not be resolved")
					}
				}
			}
		})
	}
}

type memalert struct {
	gen.AlertApiClient

	alerts map[string]*gen.Alert
}

func (m *memalert) CreateAlert(_ context.Context, req *gen.CreateAlertRequest, _ ...grpc.CallOption) (*gen.Alert, error) {

	if req.Alert == nil {
		return nil, errors.New("alert is nil")
	}

	if m.alerts[req.Alert.Id] != nil {
		return nil, errors.New("alert already exists")
	} else {
		m.alerts[req.Alert.Id] = req.Alert
	}
	return req.Alert, nil
}

func (m *memalert) UpdateAlert(_ context.Context, req *gen.UpdateAlertRequest, _ ...grpc.CallOption) (*gen.Alert, error) {
	if m.alerts[req.Alert.Id] == nil {
		return nil, errors.New("alert doesn't exists")
	} else {
		m.alerts[req.Alert.Id] = req.Alert
	}
	return req.Alert, nil
}

func (m *memalert) ResolveAlert(_ context.Context, req *gen.ResolveAlertRequest, _ ...grpc.CallOption) (*gen.Alert, error) {
	if m.alerts[req.Alert.Id] == nil {
		return nil, errors.New("alert doesn't exists")
	} else {
		m.alerts[req.Alert.Id].ResolveTime = req.Alert.ResolveTime
	}
	return req.Alert, nil
}

func (m *memalert) DeleteAlert(_ context.Context, req *gen.DeleteAlertRequest, _ ...grpc.CallOption) (*gen.DeleteAlertResponse, error) {

	if m.alerts[req.Id] == nil {
		return nil, errors.New("alert doesn't exist")
	}
	delete(m.alerts, req.Id)
	return &gen.DeleteAlertResponse{}, nil
}

func (m *memalert) ListAlerts(ctx context.Context, in *gen.ListAlertsRequest, _ ...grpc.CallOption) (*gen.ListAlertsResponse, error) {
	// sort the map and then return
	var alerts []*gen.Alert
	for _, a := range m.alerts {
		alerts = append(alerts, a)
	}

	sort.Slice(alerts, func(i, j int) bool {
		return alerts[i].CreateTime.AsTime().Before(alerts[j].CreateTime.AsTime())
	})
	return &gen.ListAlertsResponse{
		Alerts: alerts,
	}, nil
}
