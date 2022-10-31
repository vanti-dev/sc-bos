package merge

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/fanspeed"
	"github.com/vanti-dev/bsp-ew/internal/driver/bacnet/config"
	"github.com/vanti-dev/bsp-ew/internal/driver/bacnet/known"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/gobacnet"
	"math"
	"time"
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
	client *gobacnet.Client
	known  known.Context

	model *fanspeed.Model
	*fanspeed.ModelServer
	config fanSpeedConfig
}

func newFanSpeed(client *gobacnet.Client, ctx known.Context, config config.RawTrait) (*fanSpeed, error) {
	cfg, err := readFanSpeedConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	model := fanspeed.NewModel()
	return &fanSpeed{
		client:      client,
		known:       ctx,
		model:       model,
		ModelServer: fanspeed.NewModelServer(model),
		config:      cfg,
	}, nil
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
	err := writeProperty(ctx, t.client, t.known, *t.config.Speed, newFanSpeed, 0)
	if err != nil {
		return nil, err
	}

	// todo: not strictly correct as we're not paying attention to the require customisation properties that ModelServer would give us
	return t.pollUntil(ctx, 5, func(data *traits.FanSpeed) bool {
		return data.GetPercentage() == newFanSpeed
	})
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
	if err != nil {
		return nil, err
	}
	data := &traits.FanSpeed{
		Percentage: speed,
		Direction:  traits.FanSpeed_DIRECTION_UNSPECIFIED,
		Preset:     t.speedToPreset(speed),
	}
	return t.model.UpdateFanSpeed(data)
}

// pollUntil calls pollPeer until test returns true.
// Returns early with error if
//
//  1. ctx is done
//  2. the number of polls is tries
//  3. pollPeer returns an error
//
// An backoff delay will be added between each call to pollPeer
func (t *fanSpeed) pollUntil(ctx context.Context, tries int, test func(data *traits.FanSpeed) bool) (*traits.FanSpeed, error) {
	if tries == 0 {
		tries = math.MaxInt
	}

	var delay time.Duration
	delayMulti := 1.2
	var attempt int
	for {
		attempt++ // start with attempt 1 (not 0)

		res, err := t.pollPeer(ctx)
		if err != nil {
			return nil, err
		}

		if test(res) {
			return res, nil
		}

		if delay == 0 {
			delay = 10 * time.Millisecond
		} else {
			delay = time.Duration(float64(delay) * delayMulti)
		}

		if attempt >= tries {
			break
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}
	return nil, fmt.Errorf("ran out of tries: %d", tries)
}
