package auto

import (
	"context"
	"time"

	"golang.org/x/exp/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensorpb"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

func OccupancySensorAuto(model *occupancysensorpb.Model) *service.Service[string] {
	slc := service.New(service.MonoApply(func(ctx context.Context, _ string) error {
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			for {
				state := traits.Occupancy_State(rand.Intn(3) + 1)
				occupancy := &traits.Occupancy{State: state}
				if state == traits.Occupancy_OCCUPIED {
					occupancy.PeopleCount = int32(rand.Intn(10) + 1)
				}
				_, _ = model.SetOccupancy(occupancy, resource.WithUpdatePaths("state", "people_count"))
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
				}
			}
		}()
		return nil
	}), service.WithParser(func(data []byte) (string, error) {
		return string(data), nil
	}))
	_, _ = slc.Configure([]byte{}) // call configure to ensure we load when start is called.
	return slc
}
