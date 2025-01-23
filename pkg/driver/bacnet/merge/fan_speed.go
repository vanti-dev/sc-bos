package merge

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/fanspeed"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/comm"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/status"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

type fanSpeedConfig struct {
	config.Trait
	Speed   *config.ValueSource `json:"speed,omitempty"`
	Presets map[string]float32  `json:"presets,omitempty"`
}

func readFanSpeedConfig(raw []byte) (cfg fanSpeedConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	return
}

type fanSpeed struct {
	client   *gobacnet.Client
	known    known.Context
	statuses *statuspb.Map
	logger   *zap.Logger

	model *fanspeed.Model
	*fanspeed.ModelServer
	config   fanSpeedConfig
	pollTask *task.Intermittent
}

func newFanSpeed(client *gobacnet.Client, devices known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*fanSpeed, error) {
	cfg, err := readFanSpeedConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	model := fanspeed.NewModel()
	t := &fanSpeed{
		client:      client,
		known:       devices,
		statuses:    statuses,
		logger:      logger,
		model:       model,
		ModelServer: fanspeed.NewModelServer(model),
		config:      cfg,
	}
	t.pollTask = task.NewIntermittent(t.startPoll)
	initTraitStatus(statuses, cfg.Name, "FanSpeed")
	return t, nil
}

func (t *fanSpeed) startPoll(init context.Context) (stop task.StopFn, err error) {
	return startPoll(init, "fan speed", t.config.PollPeriodDuration(), t.config.PollTimeoutDuration(), t.logger, func(ctx context.Context) error {
		_, err := t.pollPeer(ctx)
		return err
	})
}

func (t *fanSpeed) AnnounceSelf(a node.Announcer) node.Undo {
	return a.Announce(t.config.Name, node.HasTrait(trait.FanSpeed, node.WithClients(fanspeed.WrapApi(t))))
}

func (t *fanSpeed) GetFanSpeed(ctx context.Context, request *traits.GetFanSpeedRequest) (*traits.FanSpeed, error) {
	_, err := t.pollPeer(ctx)
	if err != nil {
		return nil, err
	}
	return t.ModelServer.GetFanSpeed(ctx, request)
}

func (t *fanSpeed) UpdateFanSpeed(ctx context.Context, request *traits.UpdateFanSpeedRequest) (*traits.FanSpeed, error) {
	newPreset := request.GetFanSpeed().GetPreset()
	newFanSpeed := request.GetFanSpeed().GetPercentage()
	if newPreset != "" {
		presetSpeed, ok := t.config.Presets[newPreset]
		if !ok {
			return nil, fmt.Errorf("unknown preset %s", newPreset)
		}
		newFanSpeed = presetSpeed
	}
	err := comm.WriteProperty(ctx, t.client, t.known, *t.config.Speed, newFanSpeed, 0)
	if err != nil {
		return nil, err
	}

	// todo: not strictly correct as we're not paying attention to the require customisation properties that ModelServer would give us
	return pollUntil(ctx, 5, t.pollPeer, func(data *traits.FanSpeed) bool {
		return data.GetPercentage() == newFanSpeed
	})
}

func (t *fanSpeed) PullFanSpeed(request *traits.PullFanSpeedRequest, server traits.FanSpeedApi_PullFanSpeedServer) error {
	_ = t.pollTask.Attach(server.Context())
	return t.ModelServer.PullFanSpeed(request, server)
}

func (t *fanSpeed) speedToPreset(speed float32) string {
	for preset, candidate := range t.config.Presets {
		if candidate == speed {
			return preset
		}
	}
	return ""
}

// pollPeer fetches data from the peer device and saves the data locally.
func (t *fanSpeed) pollPeer(ctx context.Context) (*traits.FanSpeed, error) {
	speed, err := readPropertyFloat32(ctx, t.client, t.known, *t.config.Speed)
	status.UpdatePollErrorStatus(t.statuses, t.config.Name, "FanSpeed", []string{"speed"}, []error{err})
	if err != nil {
		return nil, comm.ErrReadProperty{Prop: "speed", Cause: err}
	}
	data := &traits.FanSpeed{
		Percentage: speed,
		Direction:  traits.FanSpeed_DIRECTION_UNSPECIFIED,
		Preset:     t.speedToPreset(speed),
	}
	return t.model.UpdateFanSpeed(data)
}
