package meteremail

import (
	"github.com/vanti-dev/sc-bos/pkg/auto/meteremail/config"

	"time"
)

const MeterTypeWater = 1
const MeterTypeElectric = 2

type MeterType uint8

type Attrs struct {
	Now   time.Time
	Stats []Stats
}

type Stats struct {
	Source       config.Source
	MeterReading MeterReading
}

type MeterReading struct {
	MeterType MeterType
	Date      time.Time
	Reading   float32
}
