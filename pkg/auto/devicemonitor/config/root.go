package config

import (
	"encoding/json"
	"time"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

type Device struct {
	Name                 string `json:"name,omitempty"` // the name of the device
	AbnormalHighId       string
	AbnormalLowId        string
	SetPointNotReachedId string

	PreviousSetPoint    *float64
	SetPointChangedTime time.Time
}

type AirTempConfig struct {
	// if the return air temperature is equal to or greater than this value, it is ok.
	// This MUST be configured to run abnormalTemperatureCheck
	OkRtLowerBound *float64 `json:"okRtLowerBound,omitempty,omitzero"`
	// This MUST be configured to run abnormalTemperatureCheck
	OkRtUpperBound *float64 `json:"okRtUpperBound,omitempty,omitzero"` // if the return air temperature is equal to or less than this value, it is ok.

	// if the return air temperature settles within the fcu set point tolerance within this time,
	// then it is considered ok. If not, a status is raised indicating the fcu might not be OK.
	// This duration should be in hours not minutes, as gets reset whenever the set point is adjusted
	// This MUST be configured to run the set point reached check
	OkSettlingTime *jsontypes.Duration `json:"okSettlingTime,omitempty,omitzero"`
	// The tolerance (+-) we have for monitoring whether the temperature settled at its target.
	// e.g. if the set point is 22°C and the tolerance is 2°C, then the temperature is considered settled
	// when it is between 20°C and 24°C.
	// This MUST be configured to run the set point reached check
	Tolerance *float64 `json:"tolerance,omitempty,omitzero"`

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

	if cfg.AirTempConfig != nil && cfg.AirTempConfig.MonitorSchedule == nil {
		cfg.AirTempConfig.MonitorSchedule = &jsontypes.Schedule{
			Schedule: jsontypes.MustParseSchedule("0 * * * *"),
		}
	}

	return cfg, nil
}
