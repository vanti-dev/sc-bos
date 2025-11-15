package config

import (
	"encoding/json"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

// Default values for config fields
const (
	DefaultConnectTimeout    = 1 * time.Second
	DefaultRefreshOccupancy  = 10 * time.Second
	DefaultRefreshStatus     = 1 * time.Minute
	DefaultRxBufferSize      = 1024
	DefaultRxBufferSizeMin   = 1024
	DefaultRxBufferSizeMax   = 65536
	DefaultRxTimeout         = 1 * time.Second
	DefaultTxTimeout         = 5 * time.Millisecond
	DefaultSendPacketTimeout = 5 * time.Millisecond
	DefaultPort              = ":50000"
)

// Root represents the root configuration for the HelvarNet driver.
type Root struct {
	driver.BaseConfig

	ConnectTimeout *jsontypes.Duration `json:"connectTimeout,omitempty"`

	EmergencyLights []*Device `json:"emergencyLights,omitempty"`
	Lights          []*Device `json:"lights,omitempty"`
	LightingGroups  []*Device `json:"lightingGroups,omitempty"`
	Pirs            []*Device `json:"pirs,omitempty"`
	// RefreshOccupancy is the duration at which the pir sensors refresh their occupancy status
	// Defaults to every 10 seconds
	RefreshOccupancy *jsontypes.Duration `json:"refreshOccupancy,omitempty,omitzero"`
	// RefreshStatus is the duration between each command to query the device state
	// Defaults to every 1 min
	RefreshStatus *jsontypes.Duration `json:"refreshStatus,omitempty,omitzero"`
	// RxBufferSize is the size of the receive buffer for the TCP connection, defaults to 1024 bytes
	RxBufferSize *int `json:"rxBufferSize,omitempty,omitzero"`
	// RxTimeout is the duration to wait for a response from querying the device, defaults to 1 second
	RxTimeout *jsontypes.Duration `json:"rxTimeout,omitempty,omitzero"`
	// TxTimeout is the send timeout, defaults to 5 milliseconds
	TxTimeout *jsontypes.Duration `json:"txTimeout,omitempty,omitzero"`
	// SendPacketTimeout is the timeout used for sending packets, defaults to 5 milliseconds
	SendPacketTimeout *jsontypes.Duration `json:"sendPacketTimeout,omitempty,omitzero"`
	// Port is the TCP port used to connect to Helvarnet devices, defaults to ":50000"
	Port *string `json:"port,omitempty"`
	// RetrySleepDuration is the duration to wait before retrying a failed operation, defaults 500 microseconds
	RetrySleepDuration *jsontypes.Duration `json:"retrySleepDuration,omitempty,omitzero"`
}

// Device represents a HelvarNet device, which can be a light, lighting group, or PIR sensor.
//
// Name is the Smart Core device name.
// Address is the Helvarnet device address, in the format <cluster>.<router>.<subnet>.<device>.(<subdevice> - optional).
// GroupNumber is the Helvarnet group number.
// IpAddress is the device's IP address.
// Meta contains additional metadata for the device.
// DurationTestLength is the length of the duration test for emergency lights, if known.
// TopicPrefix is the topic prefix to use for the UDMI automation, without the trailing '/'. If empty, the device name will be used.
type Device struct {
	Name               string              `json:"name,omitempty"`
	Address            string              `json:"address,omitempty"`
	GroupNumber        *int                `json:"groupNumber,omitempty"`
	IpAddress          string              `json:"ipAddress,omitempty"`
	Meta               *traits.Metadata    `json:"meta,omitempty"`
	DurationTestLength *jsontypes.Duration `json:"durationTestLength,omitempty,omitzero"`
	TopicPrefix        string              `json:"topicPrefix,omitempty"`
}

// Scene represents a HelvarNet lighting scene, which is a combination of a block (address), scene, and title.
//
// Block is the address block for the scene.
// Scene is the scene identifier.
// Title is a human-readable title for the scene.
type Scene struct {
	Block string `json:"block,omitempty"`
	Scene string `json:"scene,omitempty"`
	Title string `json:"title,omitempty"`
}

// ParseConfig parses the JSON configuration data into a Root struct and sets default values for optional fields.
//
// It returns the parsed Root configuration and any error encountered during parsing.
func ParseConfig(data []byte) (Root, error) {
	root := Root{}

	err := json.Unmarshal(data, &root)

	if err != nil {
		return Root{}, err
	}

	if root.ConnectTimeout == nil {
		root.ConnectTimeout = &jsontypes.Duration{Duration: DefaultConnectTimeout}
	}

	if root.RefreshOccupancy == nil {
		root.RefreshOccupancy = &jsontypes.Duration{Duration: DefaultRefreshOccupancy}
	}

	if root.RefreshStatus == nil {
		root.RefreshStatus = &jsontypes.Duration{Duration: DefaultRefreshStatus}
	}

	if root.RxTimeout == nil {
		root.RxTimeout = &jsontypes.Duration{Duration: DefaultRxTimeout}
	}

	if root.RxBufferSize == nil || *root.RxBufferSize < DefaultRxBufferSizeMin || *root.RxBufferSize > DefaultRxBufferSizeMax {
		root.RxBufferSize = new(int)
		*root.RxBufferSize = DefaultRxBufferSize
	}

	if root.TxTimeout == nil {
		root.TxTimeout = &jsontypes.Duration{Duration: DefaultTxTimeout}
	}

	if root.SendPacketTimeout == nil {
		root.SendPacketTimeout = &jsontypes.Duration{Duration: DefaultSendPacketTimeout}
	}

	if root.Port == nil {
		root.Port = new(string)
		*root.Port = DefaultPort
	}

	if root.RetrySleepDuration == nil {
		root.RetrySleepDuration = &jsontypes.Duration{Duration: 500 * time.Microsecond}
	}

	for _, device := range root.EmergencyLights {
		if device.TopicPrefix == "" {
			device.TopicPrefix = device.Name
		}
	}
	for _, device := range root.Lights {
		if device.TopicPrefix == "" {
			device.TopicPrefix = device.Name
		}
	}
	for _, device := range root.LightingGroups {
		if device.TopicPrefix == "" {
			device.TopicPrefix = device.Name
		}
	}
	for _, device := range root.Pirs {
		if device.TopicPrefix == "" {
			device.TopicPrefix = device.Name
		}
	}

	return root, nil
}

type State struct {
	State       string              `json:"state,omitempty"`
	Description string              `json:"description,omitempty"`
	FlagValue   uint32              `json:"flagValue,omitempty"`
	Level       gen.StatusLog_Level `json:"smartcoreStatusLevel,omitempty"`
}

// DeviceStatuses lists all possible device states and their associated metadata.
var DeviceStatuses = []State{
	{"Disabled", "Device or subdevice has been disabled, usually an IR subdevice or a DMX channel", 0x00000001, gen.StatusLog_NON_FUNCTIONAL},
	{"LampFailure", "Unspecified lamp problem", 0x00000002, gen.StatusLog_NON_FUNCTIONAL},
	{"Missing", "The device previously existed but is not currently present", 0x00000004, gen.StatusLog_NON_FUNCTIONAL},
	{"Faulty", "Ran out of addresses (DALI subnet) / unknown Digidim control ndevice / DALI load that keeps responding with multi-replies", 0x00000008, gen.StatusLog_NON_FUNCTIONAL},
	{"Refreshing", "DALI subnet, DALI load or Digidim control device is being discovered", 0x00000010, gen.StatusLog_NOTICE},
	{"Resting", "The load is intentionally off whilst the control gear is being powered by the emergency supply", 0x00000100, gen.StatusLog_REDUCED_FUNCTION},
	{"InEmergency", "No mains power is being supplied", 0x00000400, gen.StatusLog_REDUCED_FUNCTION},
	{"InProlong", "Mains has been restored but device is still using the emergency supply", 0x00000800, gen.StatusLog_NOTICE},
	{"FTInProgress", "The Functional Test is in progress (brief test where the control gear is being powered by the emergency supply)", 0x00001000, gen.StatusLog_NOTICE},
	{"DTInProgress", "The Duration Test is in progress. This test involves operating the control gear using the battery until the battery is completely discharged. The duration that the control gear was operational for is recorded, and then the battery recharges itself from the mains supply", 0x00002000, gen.StatusLog_NOTICE},
	{"DTPending", "The Duration Test has been requested but has not yet commenced. The test can be delayed if the battery is not fully charged", 0x00010000, gen.StatusLog_NOTICE},
	{"FTPending", "The Functional Test has been requested but has not yet commenced. The test can be delayed if there is not enough charge in the battery", 0x00020000, gen.StatusLog_NOTICE},
	{"BatteryFail", "Battery has failed", 0x00040000, gen.StatusLog_REDUCED_FUNCTION},
	{"Inhibit", "Prevents an emergency fitting from going into emergency mode", 0x00200000, gen.StatusLog_NOTICE},
	{"FTRequested", "Emergency Function Test has been requested", 0x00400000, gen.StatusLog_NOTICE},
	{"DTRequested", "Emergency Duration Test has been requested", 0x00800000, gen.StatusLog_NOTICE},
	{"Unknown", "Initial state of an emergency fitting", 0x01000000, gen.StatusLog_NOTICE},
	{"OverTemperature", "Load is over temperature/heating", 0x02000000, gen.StatusLog_REDUCED_FUNCTION},
	{"OverCurrent", "Too much current is being drawn by the load", 0x04000000, gen.StatusLog_REDUCED_FUNCTION},
	{"CommsError", "Communications error", 0x08000000, gen.StatusLog_REDUCED_FUNCTION},
	{"SevereError", "Indicates that a load is either over temperature or drawing too much current, or both", 0x10000000, gen.StatusLog_REDUCED_FUNCTION},
	{"BadReply", "Indicates that a reply to a query was malformed", 0x20000000, gen.StatusLog_NOTICE},
	{"DeviceMismatch", "The actual load type does not match the expected type", 0x80000000, gen.StatusLog_NOTICE},
}

func GetStatusListFromFlag(flag uint32) []string {

	if flag == 0 {
		return []string{"OK"}
	}

	var statusList []string
	for _, ds := range DeviceStatuses {
		if flag&ds.FlagValue != 0 {
			statusList = append(statusList, ds.State)
		}
	}

	if len(statusList) == 0 {
		// There are some flags which are NSReserved / Internal use only, so if none of the known flags are set, just return OK
		statusList = append(statusList, "OK")
	}

	return statusList
}
