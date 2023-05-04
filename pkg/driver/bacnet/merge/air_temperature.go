package merge

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperature"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
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
	client   *gobacnet.Client
	known    known.Context
	statuses *statuspb.Map
	logger   *zap.Logger

	model *airtemperature.Model
	*airtemperature.ModelServer
	config   airTemperatureConfig
	pollTask *task.Intermittent
}

func newAirTemperature(client *gobacnet.Client, devices known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*airTemperature, error) {
	cfg, err := readAirTemperatureConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	model := airtemperature.NewModel(&traits.AirTemperature{})
	t := &airTemperature{
		client:      client,
		known:       devices,
		statuses:    statuses,
		logger:      logger,
		model:       model,
		ModelServer: airtemperature.NewModelServer(model),
		config:      cfg,
	}
	t.pollTask = task.NewIntermittent(t.startPoll)
	return t, nil
}

func (t *airTemperature) startPoll(init context.Context) (stop task.StopFn, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(t.config.PollPeriodDuration())
	go func() {
		cleanup := func() {}
		defer func() { cleanup() }()
		for {
			cleanup()
			ctx, stop := context.WithTimeout(ctx, t.config.PollTimeoutDuration())
			cleanup = stop
			_, err := t.pollPeer(ctx)
			LogPollError(t.logger, "air temperature poll error", err)
			select {
			case <-ticker.C:
			case <-ctx.Done():
				return
			}
		}
	}()
	return cancel, nil
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
	if request.GetState().GetTemperatureSetPoint() == nil {
		return t.GetAirTemperature(ctx, &traits.GetAirTemperatureRequest{Name: request.Name})
	}
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

func (t *airTemperature) PullAirTemperature(request *traits.PullAirTemperatureRequest, server traits.AirTemperatureApi_PullAirTemperatureServer) error {
	_ = t.pollTask.Attach(server.Context())
	return t.ModelServer.PullAirTemperature(request, server)
}

// pollPeer fetches data from the peer device and saves the data locally.
func (t *airTemperature) pollPeer(ctx context.Context) (*traits.AirTemperature, error) {
	responses := readProperties(ctx, t.client, t.known, *t.config.SetPoint, *t.config.AmbientTemperature)
	var errs []error
	setPoint, err := float64Value(responses[0])
	if err != nil {
		errs = append(errs, ErrReadProperty{Prop: "setPoint", Cause: err})
	}
	ambientTemperature, err := float64Value(responses[1])
	if err != nil {
		errs = append(errs, ErrReadProperty{Prop: "ambientTemperature", Cause: err})
	}
	updatePollErrorStatus(t.statuses, t.config.Name, 2, errs...)
	if len(errs) > 0 {
		return nil, multierr.Combine(errs...)
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
