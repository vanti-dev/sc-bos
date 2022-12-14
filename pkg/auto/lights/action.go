package lights

import (
	"context"
	"fmt"

	"github.com/smart-core-os/sc-api/go/traits"

	"github.com/vanti-dev/sc-bos/pkg/node"
)

// actions defines the only side effects the automation can have.
// This is intended to allow easier testing of the business logic, a bit like a DAO would for database access.
type actions interface {
	// UpdateBrightness sends a LightApiClient.UpdateBrightness request and stores successful result in state.
	UpdateBrightness(ctx context.Context, req *traits.UpdateBrightnessRequest, state *WriteState) error
}

// newActions creates an actions backed by node.Clienter clients.
func newActions(clients node.Clienter) (actions, error) {
	res := &clientActions{}
	if err := clients.Client(&res.lightClient); err != nil {
		return nil, fmt.Errorf("%w traits.LightApiClient", err)
	}
	return res, nil
}

type clientActions struct {
	lightClient traits.LightApiClient
}

func (a *clientActions) UpdateBrightness(
	ctx context.Context, req *traits.UpdateBrightnessRequest, state *WriteState,
) error {
	got, err := a.lightClient.UpdateBrightness(ctx, req)
	if err != nil {
		return err
	}
	state.Brightness[req.Name] = got
	return nil
}

// updateBrightnessLevelIfNeeded sets all the names devices brightness levels to level and stores successful responses in state.
// This does not send requests if state already has a named brightness level equal to level.
func updateBrightnessLevelIfNeeded(
	ctx context.Context, state *WriteState, actions actions, level float32, names ...string,
) error {
	for _, name := range names {
		if val, ok := state.Brightness[name]; ok {
			// don't do requests that won't change the write state
			if val.LevelPercent == level {
				continue
			}
		}
		err := actions.UpdateBrightness(ctx, &traits.UpdateBrightnessRequest{
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
