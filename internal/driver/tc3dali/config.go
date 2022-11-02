package tc3dali

import (
	"encoding/json"
	"fmt"

	"github.com/vanti-dev/bsp-ew/internal/driver"
	"github.com/vanti-dev/twincat3-ads-go/pkg/ads"
)

const (
	BridgeSuffix              = "_bridge"
	ResponseMailboxSuffix     = "_response"
	NotificationMailboxSuffix = "_notification"
)

type Config struct {
	driver.BaseConfig
	ADS   ADSConfig   `json:"ads"`
	Buses []BusConfig `json:"buses"`
}

type ADSConfig struct {
	NetID NetID  `json:"netID"`
	Port  uint16 `json:"port"`
}

type NetID ads.NetId

//goland:noinspection GoMixedReceiverTypes
func (n *NetID) UnmarshalJSON(buf []byte) error {
	var str string
	err := json.Unmarshal(buf, &str)
	if err != nil {
		return err
	}

	parsed, err := ParseNetID(str)
	if err != nil {
		return err
	}
	*n = parsed
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (n NetID) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("%d.%d.%d.%d.%d.%d", n[0], n[1], n[2], n[3], n[4], n[5])
	return json.Marshal(str)
}

func ParseNetID(raw string) (n NetID, err error) {
	_, err = fmt.Sscanf(raw, "%d.%d.%d.%d.%d.%d", &n[0], &n[1], &n[2], &n[3], &n[4], &n[5])
	return
}

type BusConfig struct {
	Name           string                `json:"name"`
	ControlGear    []ControlGearConfig   `json:"controlGear"`
	ControlDevices []ControlDeviceConfig `json:"controlDevices"`
	BridgePrefix   string                `json:"bridgePrefix"`
}

type ControlGearConfig struct {
	Name         string  `json:"name"`
	Emergency    bool    `json:"emergency"`
	ShortAddress uint8   `json:"shortAddress"`
	Groups       []uint8 `json:"groups"`
}

type ControlDeviceConfig struct {
	Name          string         `json:"name"`
	ShortAddress  uint8          `json:"shortAddress"`
	InstanceTypes []InstanceType `json:"instanceTypes"`
}

type InstanceType string

const (
	InstanceTypeOccupancySensor InstanceType = "occupancySensor"
)

func (c *ControlDeviceConfig) hasInstance(want InstanceType) bool {
	for _, have := range c.InstanceTypes {
		if have == want {
			return true
		}
	}
	return false
}
