package auto

import (
	"context"
	"time"

	"golang.org/x/exp/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

func Electric(model *electric.Model) service.Lifecycle {
	s := service.New(service.MonoApply(func(ctx context.Context, _ string) error {
		go func() {
			timer := time.NewTimer((30 * time.Second) + time.Duration(rand.Float32())*time.Minute)
			for {
				state := &traits.ElectricDemand{
					Current:     float32Between(0, 40),
					Voltage:     ptr(float32Between(238, 243)),
					PowerFactor: ptr(float32Between(0.7, 1.3)),
				}
				state.ApparentPower = ptr(state.Current * *state.Voltage)
				state.RealPower = ptr(*state.ApparentPower * *state.PowerFactor)
				state.ReactivePower = ptr(*state.ApparentPower * (1 - *state.PowerFactor))
				_, _ = model.UpdateDemand(state)

				select {
				case <-ctx.Done():
					return
				case <-timer.C:
					timer = time.NewTimer((30 * time.Second) + time.Duration(rand.Float32())*time.Minute)
				}
			}
		}()
		return nil
	}))
	_, _ = s.Configure([]byte(`""`)) // ensure when start is called it actually starts
	return s
}
