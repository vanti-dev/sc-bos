package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

const (
	AuthenticationBearer = "Bearer"
	AuthenticationBasic  = "Basic"
)

type Root struct {
	auto.Config

	BaseUrl string            `json:"baseUrl"`
	Site    string            `json:"site"`
	Auths   []*Authentication `json:"authentication"`
	Logs    bool              `json:"logs"`
	Sources struct {
		Occupancy   *Occupancy   `json:"occupancy,omitempty"`
		Temperature *Temperature `json:"temperature,omitempty"`
		Energy      *Energy      `json:"energy,omitempty"`
		AirQuality  *AirQuality  `json:"airQuality,omitempty"`
		Water       *Water       `json:"water,omitempty"`
	} `json:"sources"`
}

type Authentication struct {
	Type       string `json:"type"`
	SecretFile string `json:"secretFile"`
	Token      string `json:"-"`
}

type Source struct {
	Path     string             `json:"path"`
	Interval jsontypes.Duration `json:"interval"`
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

	for _, auth := range cfg.Auths {
		err := resolveAuth(auth)

		if err != nil {
			return cfg, err
		}
	}

	return
}

func resolveAuth(auth *Authentication) error {
	switch auth.Type {
	case AuthenticationBearer:
		tok, err := os.ReadFile(auth.SecretFile)
		if err != nil {
			return err
		}
		auth.Token = strings.TrimSpace(string(tok))
	case AuthenticationBasic:
		usernamePassword, err := os.ReadFile(auth.SecretFile)
		if err != nil {
			return err
		}
		auth.Token = base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(string(usernamePassword))))
	default:
		return fmt.Errorf("authentication type %s not yet supported", auth.Type)
	}

	return nil
}
