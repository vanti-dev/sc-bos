package merge

import (
	"context"
	"encoding/json"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/cmp"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperaturepb"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/comm"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/status"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

type airTemperatureConfig struct {
	config.Trait
	SetPoint           *config.ValueSource `json:"setPoint,omitempty"`
	AmbientTemperature *config.ValueSource `json:"ambientTemperature,omitempty"`
	AmbientHumidity    *config.ValueSource `json:"ambientHumidity,omitempty"`
	SetPointLow        *config.ValueSource `json:"setPointLow,omitempty"`
	SetPointHigh       *config.ValueSource `json:"setPointHigh,omitempty"`
	// SetPointDeadBand should be defined when SetPointLow & SetPointHigh are, defaults to 1
	SetPointDeadBand *float32 `json:"deadBand,omitempty,omitzero"`
}

func readAirTemperatureConfig(raw []byte) (cfg airTemperatureConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	if err == nil {
		if cfg.SetPointDeadBand == nil || *cfg.SetPointDeadBand == 0 {
			cfg.SetPointDeadBand = new(float32)
			*cfg.SetPointDeadBand = 1
		}
	}
	return
}

type airTemperature struct {
	client   *gobacnet.Client
	known    known.Context
	statuses *statuspb.Map
	logger   *zap.Logger

	model *airtemperaturepb.Model
	*airtemperaturepb.ModelServer
	config   airTemperatureConfig
	pollTask *task.Intermittent
}

func newAirTemperature(client *gobacnet.Client, devices known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*airTemperature, error) {
	cfg, err := readAirTemperatureConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	model := airtemperaturepb.NewModel(resource.WithMessageEquivalence(cmp.Equal(
		cmp.FloatValueApprox(0, 0.1), // report temperature changes of 0.1C or more
	)))
	t := &airTemperature{
		client:      client,
		known:       devices,
		statuses:    statuses,
		logger:      logger,
		model:       model,
		ModelServer: airtemperaturepb.NewModelServer(model),
		config:      cfg,
	}
	t.pollTask = task.NewIntermittent(t.startPoll)
	initTraitStatus(statuses, cfg.Name, "AirTemperature")
	return t, nil
}

func (t *airTemperature) startPoll(init context.Context) (stop task.StopFn, err error) {
	return startPoll(init, "air temperature", t.config.PollPeriodDuration(), t.config.PollTimeoutDuration(), t.logger, func(ctx context.Context) error {
		_, err := t.pollPeer(ctx)
		return err
	})
}

func (t *airTemperature) AnnounceSelf(a node.Announcer) node.Undo {
	return a.Announce(t.config.Name, node.HasTrait(trait.AirTemperature, node.WithClients(airtemperaturepb.WrapApi(t))))
}

func (t *airTemperature) GetAirTemperature(ctx context.Context, request *traits.GetAirTemperatureRequest) (*traits.AirTemperature, error) {
	_, err := t.pollPeer(ctx)
	if err != nil {
		return nil, err
	}
	return t.ModelServer.GetAirTemperature(ctx, request)
}

func (t *airTemperature) UpdateAirTemperature(ctx context.Context, request *traits.UpdateAirTemperatureRequest) (*traits.AirTemperature, error) {
	if request.GetState().GetTemperatureSetPoint() == nil {
		return t.GetAirTemperature(ctx, &traits.GetAirTemperatureRequest{Name: request.Name})
	}
	newSetPoint := float32(request.GetState().GetTemperatureSetPoint().GetValueCelsius())

	if t.config.SetPointLow != nil && t.config.SetPointHigh != nil {
		deadBand := *t.config.SetPointDeadBand
		err := comm.WriteProperty(ctx, t.client, t.known, *t.config.SetPointLow, newSetPoint-deadBand, 0)
		if err != nil {
			t.logger.Error("WriteProperty SetPointLow", zap.Error(err))
			return nil, err
		}

		err = comm.WriteProperty(ctx, t.client, t.known, *t.config.SetPointHigh, newSetPoint+deadBand, 0)
		if err != nil {
			t.logger.Error("WriteProperty SetPointHigh", zap.Error(err))
			return nil, err
		}
	}

	if t.config.SetPoint != nil {
		err := comm.WriteProperty(ctx, t.client, t.known, *t.config.SetPoint, newSetPoint, 0)
		if err != nil {
			return nil, err
		}
	}

	// todo: not strictly correct as we're not paying attention to the require customisation properties that ModelServer would give us
	return pollUntil(ctx, t.config.DefaultRWConsistencyTimeoutDuration(), t.pollPeer, func(temperature *traits.AirTemperature) bool {
		return temperature.GetTemperatureSetPoint().ValueCelsius == float64(newSetPoint)
	})
}

func (t *airTemperature) PullAirTemperature(request *traits.PullAirTemperatureRequest, server traits.AirTemperatureApi_PullAirTemperatureServer) error {
	_ = t.pollTask.Attach(server.Context())
	return t.ModelServer.PullAirTemperature(request, server)
}

// pollPeer fetches data from the peer device and saves the data locally.
func (t *airTemperature) pollPeer(ctx context.Context) (*traits.AirTemperature, error) {
	data := &traits.AirTemperature{}
	var resProcessors []func(response any) error
	var readValues []config.ValueSource
	var requestNames []string

	if t.config.SetPoint != nil {
		requestNames = append(requestNames, "setPoint")
		readValues = append(readValues, *t.config.SetPoint)
		resProcessors = append(resProcessors, func(response any) error {
			setPoint, err := comm.Float64Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "setPoint", Cause: err}
			}
			data.TemperatureGoal = &traits.AirTemperature_TemperatureSetPoint{
				TemperatureSetPoint: &types.Temperature{ValueCelsius: setPoint},
			}
			return nil
		})
	}

	if t.config.SetPointHigh != nil {
		requestNames = append(requestNames, "setPointHigh")
		readValues = append(readValues, *t.config.SetPointHigh)
		resProcessors = append(resProcessors, func(response any) error {
			setPointHigh, err := comm.Float64Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "setPointHigh", Cause: err}
			}
			setPoint := setPointHigh - float64(*t.config.SetPointDeadBand)
			data.TemperatureGoal = &traits.AirTemperature_TemperatureSetPoint{
				TemperatureSetPoint: &types.Temperature{ValueCelsius: setPoint},
			}
			return nil
		})
	} else if t.config.SetPointLow != nil {
		requestNames = append(requestNames, "setPointLow")
		readValues = append(readValues, *t.config.SetPointLow)
		resProcessors = append(resProcessors, func(response any) error {
			setPointLow, err := comm.Float64Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "setPointLow", Cause: err}
			}
			setPoint := setPointLow + float64(*t.config.SetPointDeadBand)
			data.TemperatureGoal = &traits.AirTemperature_TemperatureSetPoint{
				TemperatureSetPoint: &types.Temperature{ValueCelsius: setPoint},
			}
			return nil
		})
	}

	if t.config.AmbientTemperature != nil {
		requestNames = append(requestNames, "ambientTemperature")
		readValues = append(readValues, *t.config.AmbientTemperature)
		resProcessors = append(resProcessors, func(response any) error {
			ambientTemperature, err := comm.Float64Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "ambientTemperature", Cause: err}
			}
			data.AmbientTemperature = &types.Temperature{ValueCelsius: ambientTemperature}
			return nil
		})
	}
	responses := comm.ReadProperties(ctx, t.client, t.known, readValues...)
	var errs []error
	for i, response := range responses {
		err := resProcessors[i](response)
		if err != nil {
			errs = append(errs, err)
		}
	}
	status.UpdatePollErrorStatus(t.statuses, t.config.Name, "AirTemperature", requestNames, errs)
	if len(errs) > 0 {
		return nil, multierr.Combine(errs...)
	}

	return t.model.UpdateAirTemperature(data)
}
