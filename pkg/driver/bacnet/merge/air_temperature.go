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

type modeDataPoints struct {
	FanOnValue     *float32
	HeatingOnValue *float32
	CoolingOnValue *float32
}

type airTempModeConfig struct {
	FanOn *config.ValueSource `json:"fanOn,omitempty"`
	// When FanOn reads equal to or above this value, the fan is considered on. Defaults to 1.
	FanOnThreshold *float32            `json:"fanOnThreshold,omitempty"`
	HeatingOn      *config.ValueSource `json:"heatingOn,omitempty"`
	// When HeatingOn reads equal to or above this value, the heating is considered on. Defaults to 1.
	HeatingOnThreshold *float32            `json:"heatingOnThreshold,omitempty"`
	CoolingOn          *config.ValueSource `json:"coolingOn,omitempty"`
	// When CoolingOn reads equal to or above this value, the cooling is considered on. Defaults to 1.
	CoolingOnThreshold *float32 `json:"coolingOnThreshold,omitempty"`
}

type airTemperatureConfig struct {
	config.Trait
	SetPoint           *config.ValueSource `json:"setPoint,omitempty"`
	AmbientTemperature *config.ValueSource `json:"ambientTemperature,omitempty"`
	AmbientHumidity    *config.ValueSource `json:"ambientHumidity,omitempty"`
	SetPointLow        *config.ValueSource `json:"setPointLow,omitempty"`
	SetPointHigh       *config.ValueSource `json:"setPointHigh,omitempty"`
	// SetPointDeadBand should be defined when SetPointLow & SetPointHigh are, defaults to 1
	SetPointDeadBand *float32           `json:"deadBand,omitempty,omitzero"`
	ModeConfig       *airTempModeConfig `json:"modeConfig,omitempty"`
}

func readAirTemperatureConfig(raw []byte) (cfg airTemperatureConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	if err == nil {
		if cfg.SetPointDeadBand == nil || *cfg.SetPointDeadBand == 0 {
			cfg.SetPointDeadBand = new(float32)
			*cfg.SetPointDeadBand = 1
		}
		if cfg.ModeConfig != nil {
			if cfg.ModeConfig.FanOnThreshold == nil {
				cfg.ModeConfig.FanOnThreshold = new(float32)
				*cfg.ModeConfig.FanOnThreshold = 1
			}
			if cfg.ModeConfig.HeatingOnThreshold == nil {
				cfg.ModeConfig.HeatingOnThreshold = new(float32)
				*cfg.ModeConfig.HeatingOnThreshold = 1
			}
			if cfg.ModeConfig.CoolingOnThreshold == nil {
				cfg.ModeConfig.CoolingOnThreshold = new(float32)
				*cfg.ModeConfig.CoolingOnThreshold = 1
			}
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
	modeData := &modeDataPoints{}

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

	if t.config.ModeConfig != nil {
		proc, vals, names := t.getModePoints(modeData)
		resProcessors = append(resProcessors, proc...)
		readValues = append(readValues, vals...)
		requestNames = append(requestNames, names...)
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
	if t.config.ModeConfig != nil {
		updateMode(t.config.ModeConfig, modeData, data)
	}
	return t.model.UpdateAirTemperature(data)
}

// getModePoints fetches all the data points needed to calculate the current AirTemperature.Mode.
func (t *airTemperature) getModePoints(data *modeDataPoints) ([]func(response any) error, []config.ValueSource, []string) {
	var processors []func(response any) error
	var values []config.ValueSource
	var names []string
	if t.config.ModeConfig.FanOn != nil {
		names = append(names, "fanOn")
		values = append(values, *t.config.ModeConfig.FanOn)
		processors = append(processors, func(response any) error {
			fanOn, err := comm.Float64Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "fanOn", Cause: err}
			}
			data.FanOnValue = new(float32)
			*data.FanOnValue = float32(fanOn)
			return nil
		})
	}

	if t.config.ModeConfig.HeatingOn != nil {
		names = append(names, "heatingOn")
		values = append(values, *t.config.ModeConfig.HeatingOn)
		processors = append(processors, func(response any) error {
			heatingOn, err := comm.Float64Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "heatingOn", Cause: err}
			}
			data.HeatingOnValue = new(float32)
			*data.HeatingOnValue = float32(heatingOn)
			return nil
		})
	}

	if t.config.ModeConfig.CoolingOn != nil {
		names = append(names, "coolingOn")
		values = append(values, *t.config.ModeConfig.CoolingOn)
		processors = append(processors, func(response any) error {
			coolingOn, err := comm.Float64Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "coolingOn", Cause: err}
			}
			data.CoolingOnValue = new(float32)
			*data.CoolingOnValue = float32(coolingOn)
			return nil
		})
	}
	return processors, values, names
}

// updateMode updates the AirTemperature.Mode based on the current values of the mode data points.
func updateMode(modeCfg *airTempModeConfig, data *modeDataPoints, airTemp *traits.AirTemperature) {
	m := traits.AirTemperature_MODE_UNSPECIFIED

	allBelow := true
	if data.CoolingOnValue != nil && *data.CoolingOnValue >= *modeCfg.CoolingOnThreshold {
		allBelow = false
		if data.FanOnValue == nil {
			m = traits.AirTemperature_COOL
		} else if *data.FanOnValue >= *modeCfg.FanOnThreshold {
			m = traits.AirTemperature_COOL
		} else {
			// we have a fan but it is off, therefore the unit is off
			m = traits.AirTemperature_OFF
		}
	}
	if data.HeatingOnValue != nil && *data.HeatingOnValue >= *modeCfg.HeatingOnThreshold {
		allBelow = false
		if data.FanOnValue == nil {
			if m == traits.AirTemperature_MODE_UNSPECIFIED {
				m = traits.AirTemperature_HEAT
			} else {
				m = traits.AirTemperature_HEAT_COOL
			}
		} else if *data.FanOnValue >= *modeCfg.FanOnThreshold {
			if m == traits.AirTemperature_MODE_UNSPECIFIED {
				m = traits.AirTemperature_HEAT
			} else {
				m = traits.AirTemperature_HEAT_COOL
			}
		} else {
			m = traits.AirTemperature_OFF
		}
	}

	if data.FanOnValue != nil && *data.FanOnValue >= *modeCfg.FanOnThreshold {
		allBelow = false
		if m == traits.AirTemperature_MODE_UNSPECIFIED {
			m = traits.AirTemperature_FAN_ONLY
		}
	}

	if allBelow &&
		(data.FanOnValue != nil || data.HeatingOnValue != nil || data.CoolingOnValue != nil) {
		m = traits.AirTemperature_OFF
	}
	airTemp.Mode = m
}
