package occupancyemail

import (
	"time"

	"github.com/vanti-dev/sc-bos/pkg/auto/occupancyemail/config"
)

type Attrs struct {
	Now   time.Time
	Stats []Stats
}

type Stats struct {
	Source    config.Source
	Last7Days OccupancyStats
}

type OccupancyStats struct {
	MaxPeopleCount int32
}
