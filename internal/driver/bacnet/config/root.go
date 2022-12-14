package config

import (
	"bytes"
	"encoding/json"
	bactypes "github.com/vanti-dev/gobacnet/types"
	"github.com/vanti-dev/sc-bos/internal/driver"
	"io"
	"net/netip"
	"os"
)

// Root represents a full collection of related configuration properties
type Root struct {
	driver.BaseConfig

	LocalInterface string `json:"localInterface,omitempty"`
	LocalPort      uint16 `json:"localPort,omitempty"`

	Discovery                *Discovery `json:"discovery,omitempty"`
	ForceDiscovery           bool       `json:"forceDiscovery,omitempty"`
	IncludeDiscoveredDevices bool       `json:"includeDiscoveredDevices,omitempty"`

	DiscoverObjects bool `json:"discoverObjects,omitempty"`

	COV *COV `json:"cov,omitempty"`

	Devices []Device   `json:"devices,omitempty"`
	Traits  []RawTrait `json:"traits,omitempty"`
}

// ReadFile reads from the named file a config Root.
func ReadFile(name string) (root Root, err error) {
	bytes, err := os.ReadFile(name)
	if err != nil {
		return root, err
	}
	err = json.Unmarshal(bytes, &root)
	return root, err
}

// Read decodes r into a config Root.
func Read(r io.Reader) (root Root, err error) {
	err = json.NewDecoder(r).Decode(&root)
	return
}

// ReadBytes decodes bytes into a config Root.
func ReadBytes(data []byte) (root Root, err error) {
	return Read(bytes.NewReader(data))
}

type Discovery struct {
	Min        int      `json:"min,omitempty"`
	Max        int      `json:"max,omitempty"`
	Chunk      int      `json:"chunk,omitempty"`
	ChunkDelay Duration `json:"chunkDelay,omitempty"`
}

type Device struct {
	Name  string                  `json:"name,omitempty"`
	Title string                  `json:"title,omitempty"`
	Comm  *Comm                   `json:"comm,omitempty"`
	ID    bactypes.ObjectInstance `json:"id,omitempty"`

	COV *COV `json:"cov,omitempty"`

	DiscoverObjects *bool    `json:"discoverObjects,omitempty"`
	Objects         []Object `json:"objects,omitempty"`
}

type Comm struct {
	IP *netip.AddrPort `json:"ip,omitempty"`
}
