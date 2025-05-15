package config

import (
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

type Device struct {
	Name           string `json:"name,omitempty"` // the name of the device
	AbnormalHighId string
	AbnormalLowId  string
}

type AirTempConfig struct {
	OkRtLowerBound *float64 `json:"okLowerBound,omitempty,omitzero"` // if the return air temperature is equal to or greater than this value, it is ok.
	OkRtUpperBound *float64 `json:"okUpperBound,omitempty,omitzero"` // if the return air temperature is equal to or less than this value, it is ok.
	// if the return air temperature settles within the fcu set point tolerance within this time,
	// then it is considered ok. If not, a status is raised indicating the fcu might not be OK.
	// This duration should be in hours not minutes, as gets reset whenever the set point is adjusted
	OkSettlingTime *jsontypes.Duration `json:"okSettlingTime,omitempty,omitzero"`
	// The tolerance (+-) we have for monitoring whether the temperature settled at its target.
	// e.g. if the set point is 22°C and the tolerance is 2°C, then the temperature is considered settled
	// when it is between 20°C and 24°C.
	Tolerance *float64 `json:"deadband,omitempty,omitzero"`

	MonitorSchedule *jsontypes.Schedule `json:"monitorSchedule,omitempty"`
	Devices         []*Device           `json:"devices,omitempty"`
}

type Root struct {
	auto.Config

	AirTempConfig *AirTempConfig `json:"airTempConfig,omitempty"`

	// Now returns the current time. It's configurable for testing purposes, typically for testing the logic.
	Now func() time.Time `json:"-"`
}

func ReadBytes(data []byte) (cfg Root, err error) {
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return
	}

	if cfg.AirTempConfig.OkRtLowerBound == nil && cfg.AirTempConfig.OkRtUpperBound == nil {
		return cfg, fmt.Errorf("at least one or both of okLowerBound or okUpperBound must be set")
	}

	if cfg.AirTempConfig.OkSettlingTime == nil || cfg.AirTempConfig.OkSettlingTime.Duration == 0 {
		cfg.AirTempConfig.OkSettlingTime = &jsontypes.Duration{
			Duration: time.Hour * 4,
		}
	}

	if cfg.AirTempConfig.Tolerance == nil || *cfg.AirTempConfig.Tolerance == 0 {
		cfg.AirTempConfig.Tolerance = proto.Float64(2.0)
	}

	if cfg.AirTempConfig.MonitorSchedule == nil {
		cfg.AirTempConfig.MonitorSchedule = &jsontypes.Schedule{
			Schedule: jsontypes.MustParseSchedule("0 * * * *"),
		}
	}

	return cfg, nil
}
