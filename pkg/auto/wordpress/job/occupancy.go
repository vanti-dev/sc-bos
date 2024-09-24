package job

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
)

// OccupancyJob gets occupancy at current point in time
type OccupancyJob struct {
	BaseJob
	client  traits.OccupancySensorApiClient
	Sensors []string
}

func (o *OccupancyJob) GetName() string {
	return "occupancy"
}

func (o *OccupancyJob) GetClients() []any {
	return []any{&o.client}
}

func (o *OccupancyJob) Do(ctx context.Context, sendFn sender) error {
	sum := int32(0)

	for _, sensor := range o.Sensors {
		resp, err := o.client.GetOccupancy(ctx, &traits.GetOccupancyRequest{Name: sensor})

		if err != nil {
			o.Logger.Error("getting occupancy", zap.String("sensor", sensor), zap.Error(err))
			continue
		}

		sum += resp.PeopleCount
	}

	body := &TotalOccupancy{
		Meta: Meta{
			Site:      o.GetSite(),
			Timestamp: time.Now(),
		},
		TotalOccupancy: IntMeasure{
			Value: sum,
			Units: "People",
		},
	}

	bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return sendFn(ctx, o.GetUrl(), bytes)
}
