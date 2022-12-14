package merge

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperature"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"math"
	"time"
)

type airTemperatureConfig struct {
	config.Trait
	SetPoint           *config.ValueSource `json:"setPoint,omitempty"`
	AmbientTemperature *config.ValueSource `json:"ambientTemperature,omitempty"`
	AmbientHumidity    *config.ValueSource `json:"ambientHumidity,omitempty"`
}

func readAirTemperatureConfig(raw []byte) (cfg airTemperatureConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	return
}

type airTemperature struct {
	client *gobacnet.Client
	known  known.Context

	model *airtemperature.Model
	*airtemperature.ModelServer
	config airTemperatureConfig
}

func newAirTemperature(client *gobacnet.Client, ctx known.Context, config config.RawTrait) (*airTemperature, error) {
	cfg, err := readAirTemperatureConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	model := airtemperature.NewModel(&traits.AirTemperature{})
	return &airTemperature{
		client:      client,
		known:       ctx,
		model:       model,
		ModelServer: airtemperature.NewModelServer(model),
		config:      cfg,
	}, nil
}

func (t *airTemperature) AnnounceSelf(a node.Announcer) node.Undo {
	return a.Announce(t.config.Name, node.HasTrait(trait.AirTemperature, node.WithClients(airtemperature.WrapApi(t))))
}

func (t *airTemperature) GetAirTemperature(ctx context.Context, request *traits.GetAirTemperatureRequest) (*traits.AirTemperature, error) {
	_, err := t.pollPeer(ctx)
	if err != nil {
		return nil, err
	}
	return t.ModelServer.GetAirTemperature(ctx, request)
}

func (t *airTemperature) UpdateAirTemperature(ctx context.Context, request *traits.UpdateAirTemperatureRequest) (*traits.AirTemperature, error) {
	newSetPoint := float32(request.GetState().GetTemperatureSetPoint().GetValueCelsius())
	err := writeProperty(ctx, t.client, t.known, *t.config.SetPoint, newSetPoint, 0)
	if err != nil {
		return nil, err
	}

	// todo: not strictly correct as we're not paying attention to the require customisation properties that ModelServer would give us
	return t.pollUntil(ctx, 5, func(temperature *traits.AirTemperature) bool {
		return temperature.GetTemperatureSetPoint().ValueCelsius == float64(newSetPoint)
	})
}

// pollPeer fetches data from the peer device and saves the data locally.
func (t *airTemperature) pollPeer(ctx context.Context) (*traits.AirTemperature, error) {
	setPoint, err := readPropertyFloat64(ctx, t.client, t.known, *t.config.SetPoint)
	if err != nil {
		return nil, err
	}
	ambientTemperature, err := readPropertyFloat64(ctx, t.client, t.known, *t.config.AmbientTemperature)
	if err != nil {
		return nil, err
	}
	data := &traits.AirTemperature{
		AmbientTemperature: &types.Temperature{ValueCelsius: ambientTemperature},
		TemperatureGoal: &traits.AirTemperature_TemperatureSetPoint{
			TemperatureSetPoint: &types.Temperature{ValueCelsius: setPoint},
		},
	}
	return t.model.UpdateAirTemperature(data)
}

// pollUntil calls pollPeer until test returns true.
// Returns early with error if
//
//  1. ctx is done
//  2. the number of polls is tries
//  3. pollPeer returns an error
//
// An backoff delay will be added between each call to pollPeer
func (t *airTemperature) pollUntil(ctx context.Context, tries int, test func(temperature *traits.AirTemperature) bool) (*traits.AirTemperature, error) {
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
