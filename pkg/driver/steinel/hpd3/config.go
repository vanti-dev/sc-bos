package hpd3

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/vanti-dev/sc-bos/pkg/driver"
)

type DriverConfig struct {
	driver.BaseConfig
	Devices []DeviceConfig `json:"devices"`
}

func ParseDriverConfig(buf []byte) (conf DriverConfig, err error) {
	err = json.Unmarshal(buf, &conf)
	return
}

type DeviceConfig struct {
	Name     string   `json:"name"`
	Host     string   `json:"host"`
	Password Password `json:"password"`
}

// Password represents a secret that can be specified directly inline or loaded from a file
type Password struct {
	Value string `json:"-"`
	File  string `json:"file,omitempty"`
}

func (p *Password) UnmarshalJSON(buf []byte) error {
	var value string
	err := json.Unmarshal(buf, &value)
	if err == nil {
		p.Value = value
		p.File = ""
	}
	if typeErr := (&json.UnmarshalTypeError{}); !errors.As(err, &typeErr) {
		return err
	}

	type wrapped Password
	return json.Unmarshal(buf, (*wrapped)(p))
}

func (p *Password) Read() (string, error) {
	if p.Value != "" {
		return p.Value, nil
	}
	if p.File != "" {
		buf, err := os.ReadFile(p.File)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(buf)), nil
	}
	return "", errors.New("no password specified")
}
