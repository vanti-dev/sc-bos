package config

import (
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

const PointsEventTopicSuffix = "/event/pointset"

type ScDevice struct {
	Meta   *traits.Metadata `json:"meta,omitempty"`
	ScName string           `json:"scName,omitempty"`
}

type Root struct {
	driver.BaseConfig
	HTTP           *HTTP  `json:"http,omitempty"`
	ScNamePrefix   string `json:"scNamePrefix,omitempty"`
	CaPath         string `json:"caPath,omitempty"`
	ClientCertPath string `json:"clientCertPath,omitempty"`
	ClientKeyPath  string `json:"clientKeyPath,omitempty"`
	// poll the cardholders api for updates on this schedule, defaults to once per minute
	RefreshCardholders *jsontypes.Schedule `json:"refreshCardholders,omitempty"`
	// poll the alerts API for updates on this schedule, defaults to once per minute
	RefreshAlarms *jsontypes.Schedule `json:"refreshAlerts,omitempty"`
	// poll the doors on this schedule, defaults to once per day
	RefreshDoors *jsontypes.Schedule `json:"refreshDoors,omitempty"`
	// the interval between udmi exports, defaults to every 5 seconds
	// the duration after which a security event will expire and be removed from any List/Pull, defaults to 24 hours
	SecurityEventExpiry jsontypes.Duration `json:"securityEventExpiry,omitempty"`
	UdmiExportInterval  jsontypes.Duration `json:"udmiExportInterval,omitempty"`
	TopicPrefix         string             `json:"topicPrefix,omitempty"`

	RefreshOccupancyInterval *jsontypes.Duration `json:"refreshOccupancyInterval,omitempty"`

	OccupancyCountEnabled bool `json:"occupancyCountEnabled,omitempty"`
}

type HTTP struct {
	BaseURL    string `json:"baseUrl,omitempty"`
	ApiKeyFile string `json:"apiKeyFile,omitempty"`
}

func (cfg *Root) ApplyDefaults() {
	if cfg.RefreshCardholders == nil {
		cfg.RefreshCardholders = jsontypes.MustParseSchedule("* * * * *")
	}

	if cfg.UdmiExportInterval.Duration == 0 {
		cfg.UdmiExportInterval.Duration = 5 * time.Second
	}

	if cfg.RefreshDoors == nil {
		cfg.RefreshDoors = jsontypes.MustParseSchedule("0 0 * * *")
	}

	if cfg.RefreshAlarms == nil {
		cfg.RefreshAlarms = jsontypes.MustParseSchedule("* * * * *")
	}

	if cfg.SecurityEventExpiry.Duration == 0 {
		cfg.SecurityEventExpiry.Duration = 24 * time.Hour
	}
}
