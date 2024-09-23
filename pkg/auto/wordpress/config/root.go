package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/vanti-dev/sc-bos/pkg/auto"
)

type Root struct {
	auto.Config

	BaseUrl string         `json:"base_url"`
	Site    string         `json:"site"`
	Auth    Authentication `json:"authentication"`
	Logs    bool           `json:"logs"`
	Sources struct {
		Occupancy   *Occupancy   `json:"occupancy,omitempty"`
		Temperature *Temperature `json:"temperature,omitempty"`
		Energy      *Energy      `json:"energy,omitempty"`
		AirQuality  *AirQuality  `json:"air_quality,omitempty"`
		Water       *Water       `json:"water,omitempty"`
	} `json:"sources"`
}

type Authentication struct {
	Type       string `json:"type"`
	SecretPath string `json:"secretFile"`
	Token      string `json:"-"`
}

type Source struct {
	Path     string        `json:"path"`
	Interval string        `json:"interval"`
	Duration time.Duration `json:"-"` // actual parsed interval
}

type Occupancy struct {
	Source
	Sensors []string `json:"sensors"`
}

type Temperature struct {
	Source
	Sensors []string `json:"sensors"`
}

type Energy struct {
	Source
	Meters []string `json:"meters"`
}

type AirQuality struct {
	Source
	Sensors []string `json:"sensors"`
}

type Water struct {
	Source
	Meters []string `json:"meters"`
}

func ReadBytes(data []byte) (cfg Root, err error) {
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return
	}

	if cfg.BaseUrl == "" {
		err = fmt.Errorf("baseUrl not specified")
		return
	}

	switch cfg.Auth.Type {
	case "Authorization Bearer":
		var tok []byte
		tok, err = os.ReadFile(cfg.Auth.SecretPath)
		if err != nil {
			return
		}
		cfg.Auth.Token = string(tok)
	default:
		err = fmt.Errorf("authentication type %s not yet supported", cfg.Auth.Type)
		return
	}

	if cfg.Sources.Occupancy != nil {
		cfg.Sources.Occupancy.Duration, err = time.ParseDuration(cfg.Sources.Occupancy.Interval)
		if err != nil {
			return
		}
	}
	if cfg.Sources.Temperature != nil {
		cfg.Sources.Temperature.Duration, err = time.ParseDuration(cfg.Sources.Temperature.Interval)
		if err != nil {
			return
		}
	}
	if cfg.Sources.Energy != nil {
		cfg.Sources.Energy.Duration, err = time.ParseDuration(cfg.Sources.Energy.Interval)
		if err != nil {
			return
		}
	}
	if cfg.Sources.AirQuality != nil {
		cfg.Sources.AirQuality.Duration, err = time.ParseDuration(cfg.Sources.AirQuality.Interval)
		if err != nil {
			return
		}
	}

	if cfg.Sources.Water != nil {
		cfg.Sources.Water.Duration, err = time.ParseDuration(cfg.Sources.Water.Interval)
		if err != nil {
			return
		}
	}
	return
}
