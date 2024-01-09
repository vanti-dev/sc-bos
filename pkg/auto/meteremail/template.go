package meteremail

import (
	"github.com/vanti-dev/sc-bos/pkg/auto/meteremail/config"

	"time"
)

type Attrs struct {
	Now   time.Time
	Stats []Stats
}

type Stats struct {
	Source        config.Source
	MeterReadings MeterReadings
}

type MeterReadings struct {
	Date time.Time
	Reading
}

type Reading struct {
	Reading float32
}
