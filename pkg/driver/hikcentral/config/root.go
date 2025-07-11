package config

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

type Root struct {
	driver.BaseConfig

	API      *API      `json:"api,omitempty"`
	Settings *Settings `json:"settings,omitempty"`

	// Metadata applied to all cameras
	Metadata *traits.Metadata `json:"metadata,omitempty"`
	Cameras  []*Camera        `json:"cameras,omitempty"`
}

type API struct {
	Address    string              `json:"address,omitempty"`
	AppKey     string              `json:"appKey,omitempty"`
	Secret     string              `json:"secret,omitempty"`
	SecretFile string              `json:"secretFile,omitempty"`
	Timeout    *jsontypes.Duration `json:"timeout,omitempty"`
}

type Settings struct {
	InfoPoll      *jsontypes.Duration `json:"infoPoll,omitempty"`
	OccupancyPoll *jsontypes.Duration `json:"occupancyPoll,omitempty"`
	EventsPoll    *jsontypes.Duration `json:"eventsPoll,omitempty"`
	StreamPoll    *jsontypes.Duration `json:"streamPoll,omitempty"`
}

type Camera struct {
	Name      string `json:"name,omitempty"`
	Topic     string `json:"topic,omitempty"`
	IndexCode string `json:"indexCode,omitempty"`
	// Metadata applied to this camera
	Metadata  *traits.Metadata `json:"metadata,omitempty"`
	IpAddress string           `json:"ipAddress,omitempty"`
}

func ReadBytes(raw []byte) (dst Root, err error) {
	dst = Root{}
	err = json.Unmarshal(raw, &dst)
	if err != nil {
		return dst, err
	}
	if dst.API.SecretFile != "" {
		dst.API.Secret, err = readSecret(dst.API.SecretFile)
		if err != nil {
			return dst, err
		}
	}
	if dst.API.Timeout == nil {
		dst.API.Timeout = &jsontypes.Duration{Duration: 5 * time.Second}
	}
	return
}

func readSecret(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	str := string(raw)
	str = strings.TrimSpace(str)
	return str, nil
}
