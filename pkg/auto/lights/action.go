package lights

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"

	"github.com/vanti-dev/sc-bos/pkg/node"
)

const brightnessCacheValidity = 45 * time.Second

// actions defines the only side effects the automation can have.
// This is intended to allow easier testing of the business logic, a bit like a DAO would for database access.
type actions interface {
	// UpdateBrightness sends a LightApiClient.UpdateBrightness request and stores successful result in state.
	UpdateBrightness(ctx context.Context, now time.Time, req *traits.UpdateBrightnessRequest, state *WriteState) error
}

// newClientActions creates an actions backed by node.Clienter clients.
func newClientActions(clients node.Clienter) (actions, error) {
	res := &clientActions{}
	if err := clients.Client(&res.lightClient); err != nil {
		return nil, fmt.Errorf("%w traits.LightApiClient", err)
	}
	return res, nil
}

type clientActions struct {
	lightClient traits.LightApiClient
}

func (a *clientActions) UpdateBrightness(ctx context.Context, now time.Time, req *traits.UpdateBrightnessRequest, state *WriteState) error {
	got, err := a.lightClient.UpdateBrightness(ctx, req)
	state.Brightness[req.Name] = Value[*traits.Brightness]{
		V:   got,
		At:  now,
		Err: err,
	}
	return err
}

// updateBrightnessLevelIfNeeded sets all the names devices brightness levels to level and stores successful responses in state.
// This does not send requests if state already has a named brightness level equal to level.
func updateBrightnessLevelIfNeeded(ctx context.Context, now time.Time, state *WriteState, actions actions, level float32, logger *zap.Logger, names ...deviceName) error {
	for _, name := range names {
		if val, ok := state.Brightness[name]; ok && val.V != nil {
			expired := now.After(val.At.Add(brightnessCacheValidity))
			// don't do requests that won't change the write state unless the entry is expired
			if val.V.LevelPercent == level && !expired {
				continue
			}
		}

		logger.Debug("Setting brightness for light fitting", zap.String("fitting name", name), zap.Float32("level", level))
		err := actions.UpdateBrightness(ctx, now, &traits.UpdateBrightnessRequest{
			Name: name,
			Brightness: &traits.Brightness{
				LevelPercent: level,
			},
		}, state)
		if err != nil {
			return err
		}
	}
	return nil
}

func refreshBrightnessLevel(ctx context.Context, now time.Time, state *WriteState, actions actions, logger *zap.Logger, names ...deviceName) error {
	for _, name := range names {
		val, ok := state.Brightness[name]
		if !ok || val.V == nil {
			continue
		}
		expired := now.After(val.At.Add(brightnessCacheValidity))
		if !expired {
			// don't need to refresh if recently written
			continue
		}

		logger.Debug("refreshing brightness for light fitting",
			zap.String("fitting name", name),
			zap.Float32("level", val.V.LevelPercent),
		)
		err := actions.UpdateBrightness(ctx, now, &traits.UpdateBrightnessRequest{
			Name:       name,
			Brightness: val.V,
		}, state)
		if err != nil {
			return err
		}
	}
	return nil
}
