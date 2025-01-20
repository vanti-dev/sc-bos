package merge

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/emergency"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/gobacnet/enum/lifesafetystate"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/comm"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/status"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

// AlarmConfig allows configuring a specific bacnet point to raise an Emergency if either:
//   - value read from that point is anything other than OkValue
//   - value read from that point is less than OkAbove
type AlarmConfig struct {
	config.ValueSource
	// Only one of the below OK values can be configured.
	OkValue     *int64 `json:"okValue,omitempty"` // what we expect to read from the point when it is ok, any other value is an emergency. Not supported for float values.
	OkAbove     *int64 `json:"okAbove,omitempty"` // the point is OK if it is equal to or greater than this value, if it isn't we have an emergency
	AlarmReason string `json:"alarmReason"`       // the reason of the alarm
}

type emergencyConfig struct {
	config.Trait
	Level       *config.ValueSource `json:"level,omitempty"`
	AlarmConfig *AlarmConfig        `json:"alarmConfig,omitempty"`
}

func readEmergencyConfig(raw []byte) (cfg emergencyConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		return
	}

	if cfg.AlarmConfig != nil {
		if cfg.AlarmConfig.OkValue != nil && cfg.AlarmConfig.OkAbove != nil {
			return cfg, multierr.Combine(err, errors.New("cannot set both okValue and okAbove"))
		}
	}
	return
}

type emergencyImpl struct {
	client   *gobacnet.Client
	known    known.Context
	statuses *statuspb.Map
	logger   *zap.Logger

	model *emergency.MemoryDevice
	traits.EmergencyApiServer
	config   emergencyConfig
	pollTask *task.Intermittent
}

func newEmergency(client *gobacnet.Client, devices known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*emergencyImpl, error) {
	cfg, err := readEmergencyConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	model := emergency.NewMemoryDevice()
	t := &emergencyImpl{
		client:             client,
		known:              devices,
		statuses:           statuses,
		logger:             logger,
		model:              model,
		EmergencyApiServer: model,
		config:             cfg,
	}
	t.pollTask = task.NewIntermittent(t.startPoll)
	initTraitStatus(statuses, cfg.Name, "Emergency")
	return t, nil
}

func (t *emergencyImpl) startPoll(init context.Context) (stop task.StopFn, err error) {
	return startPoll(init, "emergency", t.config.PollPeriodDuration(), t.config.PollTimeoutDuration(), t.logger, func(ctx context.Context) error {
		_, err := t.pollPeer(ctx)
		return err
	})
}

func (t *emergencyImpl) AnnounceSelf(a node.Announcer) node.Undo {
	return a.Announce(t.config.Name, node.HasTrait(trait.Emergency, node.WithClients(emergency.WrapApi(t))))
}

func (t *emergencyImpl) GetEmergency(ctx context.Context, request *traits.GetEmergencyRequest) (*traits.Emergency, error) {
	_, err := t.pollPeer(ctx)
	if err != nil {
		return nil, err
	}
	return t.EmergencyApiServer.GetEmergency(ctx, request)
}

func (t *emergencyImpl) UpdateEmergency(ctx context.Context, request *traits.UpdateEmergencyRequest) (*traits.Emergency, error) {
	return traits.UnimplementedEmergencyApiServer{}.UpdateEmergency(ctx, request)
}

func (t *emergencyImpl) PullEmergency(request *traits.PullEmergencyRequest, server traits.EmergencyApi_PullEmergencyServer) error {
	_ = t.pollTask.Attach(server.Context())
	return t.EmergencyApiServer.PullEmergency(request, server)
}

func (t *emergencyImpl) checkIntValueForEmergency(response any) (*traits.Emergency, error) {
	data := &traits.Emergency{}

	value, err := comm.IntValue(response)
	if err != nil {
		return nil, comm.ErrReadProperty{Prop: "alarmConfig", Cause: err}
	}

	if t.config.AlarmConfig.OkValue != nil {
		if value != *t.config.AlarmConfig.OkValue {
			data.Reason = t.config.AlarmConfig.AlarmReason
			data.Level = traits.Emergency_EMERGENCY
		} else {
			data.Level = traits.Emergency_OK
		}
	}

	if t.config.AlarmConfig.OkAbove != nil {
		if value < *t.config.AlarmConfig.OkAbove {
			data.Reason = t.config.AlarmConfig.AlarmReason
			data.Level = traits.Emergency_EMERGENCY
		} else {
			data.Level = traits.Emergency_OK
		}
	}

	return data, nil
}

func (t *emergencyImpl) checkFloatValueForEmergency(response any) (*traits.Emergency, error) {
	data := &traits.Emergency{}

	value, err := comm.Float64Value(response)
	if err != nil {
		return nil, comm.ErrReadProperty{Prop: "alarmConfig", Cause: err}
	}

	if t.config.AlarmConfig.OkAbove != nil {
		if value < float64(*t.config.AlarmConfig.OkAbove) {
			data.Reason = t.config.AlarmConfig.AlarmReason
			data.Level = traits.Emergency_EMERGENCY
		} else {
			data.Level = traits.Emergency_OK
		}
		return data, nil
	}
	return data, nil
}

// pollPeer fetches data from the peer device and saves the data locally.
func (t *emergencyImpl) pollPeer(ctx context.Context) (*traits.Emergency, error) {
	data := &traits.Emergency{}
	var resProcessors []func(response any) error
	var readValues []config.ValueSource
	var requestNames []string

	if t.config.Level != nil {
		requestNames = append(requestNames, "level")
		readValues = append(readValues, *t.config.Level)
		resProcessors = append(resProcessors, func(response any) error {
			enum, err := comm.EnumValue(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "level", Cause: err}
			}
			level := lifesafetystate.LifeSafetyState(enum)
			switch level {
			case lifesafetystate.Quiet:
				data.Level = traits.Emergency_OK
			case lifesafetystate.PreAlarm, lifesafetystate.FaultPreAlarm:
				data.Level = traits.Emergency_WARNING
			default:
				data.Level = traits.Emergency_EMERGENCY
			}

			switch level {
			case lifesafetystate.TestActive, lifesafetystate.TestAlarm, lifesafetystate.TestFault, lifesafetystate.TestFaultAlarm, lifesafetystate.TestSupervisory,
				lifesafetystate.TestOEOAlarm, lifesafetystate.TestOEOEvacuate, lifesafetystate.TestOROPhase1Recall, lifesafetystate.TestOEOUnaffected, lifesafetystate.TestOEOUnavailable:
				data.Drill = true
			default:
				data.Drill = false
			}

			return nil
		})
	}

	if t.config.AlarmConfig != nil {
		requestNames = append(requestNames, "alarmConfig")
		readValues = append(readValues, t.config.AlarmConfig.ValueSource)
		resProcessors = append(resProcessors, func(response any) error {
			switch response.(type) {
			case int32, int16, int8, uint32, uint16, uint8:
				if e, err := t.checkIntValueForEmergency(response); err == nil {
					data = e
				}
			case float32, float64:
				if e, err := t.checkFloatValueForEmergency(response); err == nil {
					data = e
				}
			default:
				t.logger.Warn("alarmConfig unknown value type", zap.String("type", reflect.TypeOf(response).String()))
			}
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
	status.UpdatePollErrorStatus(t.statuses, t.config.Name, "Emergency", requestNames, errs)
	if len(errs) > 0 {
		return nil, multierr.Combine(errs...)
	}

	return t.model.UpdateEmergency(ctx, &traits.UpdateEmergencyRequest{Emergency: data})
}
