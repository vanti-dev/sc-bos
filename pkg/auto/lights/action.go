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

type actionCounts struct {
	TotalWrites int

	BrightnessUpdates map[float32][]deviceName // [brightness]names
	BrightnessWrites  []deviceName             // which devices have been written to
}

// Changes returns a summary of the changes made by the actions.
func (a actionCounts) Changes() []string {
	var res []string
	for brightness, names := range a.BrightnessUpdates {
		if l := len(names); l > 1 {
			res = append(res, fmt.Sprintf("%dx brightness=%.1f%%", l, brightness))
		} else {
			res = append(res, fmt.Sprintf("brightness=%.1f%%", brightness))
		}
	}
	return res
}

// newCountActions counts all invocations of a and stores them in actionCounts.
func newCountActions(a actions) (actions, *actionCounts) {
	impl := &countActions{actions: a, actionCounts: &actionCounts{}}
	return impl, impl.actionCounts
}

type countActions struct {
	actions
	*actionCounts
}

func (a *countActions) UpdateBrightness(ctx context.Context, now time.Time, req *traits.UpdateBrightnessRequest, state *WriteState) error {
	a.TotalWrites++
	a.BrightnessWrites = append(a.BrightnessWrites, req.Name)
	if a.BrightnessUpdates == nil {
		a.BrightnessUpdates = make(map[float32][]deviceName)
	}
	a.BrightnessUpdates[req.Brightness.LevelPercent] = append(a.BrightnessUpdates[req.Brightness.LevelPercent], req.Name)

	return a.actions.UpdateBrightness(ctx, now, req, state)
}

// nilActions is an actions that has no side effects.
// Actions are still recorded in the WriteState.
type nilActions struct{}

func (nilActions) UpdateBrightness(ctx context.Context, now time.Time, req *traits.UpdateBrightnessRequest, state *WriteState) error {
	state.Brightness[req.Name] = Value[*traits.Brightness]{
		V:  req.Brightness,
		At: now,
	}
	return nil
}

// newLogActions logs to logger each time an action is invoked.
func newLogActions(a actions, logger *zap.Logger) actions {
	return &logActions{actions: a, logger: logger}
}

type logActions struct {
	actions
	logger *zap.Logger
}

func (a *logActions) UpdateBrightness(ctx context.Context, now time.Time, req *traits.UpdateBrightnessRequest, state *WriteState) error {
	err := a.actions.UpdateBrightness(ctx, now, req, state)
	a.logger.Debug("actions.UpdateBrightness",
		zap.String("name", req.Name),
		zap.Float32("level", req.Brightness.LevelPercent),
		zap.Error(err),
	)
	return err
}

// updateBrightnessLevelIfNeeded sets all the names devices brightness levels to level and stores successful responses in state.
// This does not send requests if state already has a named brightness level equal to level.
func updateBrightnessLevelIfNeeded(ctx context.Context, now time.Time, state *WriteState, actions actions, level float32, names ...deviceName) error {
	for _, name := range names {
		if val, ok := state.Brightness[name]; ok && val.V != nil {
			expired := now.After(val.At.Add(brightnessCacheValidity))
			// don't do requests that won't change the write state unless the entry is expired
			if val.V.LevelPercent == level && !expired {
				continue
			}
		}

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

func refreshBrightnessLevel(ctx context.Context, now time.Time, state *WriteState, actions actions, names ...deviceName) error {
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
