package config

import (
	"encoding/json"
	"fmt"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

type Root struct {
	driver.BaseConfig
	// smart core name prefix to bind the supported traits to
	ScNamePrefix string `json:"scNamePrefix"`
	// devices to be controlled by the modbus driver
	Devices []Device `json:"devices"`
}

type TCPHandler struct {
	// IP or hostname address of the modbus over IP device
	Address string `json:"address"`
	// port number of the modbus over IP device
	Port int `json:"port"`
	// timeout for the client to wait for
	Timeout *jsontypes.Duration `json:"timeout"`
	// slave id of the modbus over IP device
	SlaveId byte `json:"slaveId"`
}

type RTUHandler struct {
	// bus address of the modbus device
	Address string `json:"address"`
	// timeout for the client to wait for
	Timeout *jsontypes.Duration `json:"timeout"`
	// Baud rate (default 19200)
	BaudRate int `json:"baudRate"`
	// Data bits: 5, 6, 7 or 8 (default 8)
	DataBits int `json:"dataBits"`
	// Stop bits: 1 or 2 (default 1)
	StopBits int `json:"stopBits"`
	// Parity: N - None, E - Even, O - Odd (default E)
	// (The use of no parity requires 2 stop bits.)
	Parity string `json:"parity"`
	// slave id of the modbus device
	SlaveId byte `json:"slaveId"`
}

type Device struct {
	// unique name of the device (smart core name suffix)
	Name string `json:"name"`
	// tcp handler for the modbus over IP device
	TcpHandle *TCPHandler `json:"tcpHandle,omitempty"`
	// rtu handler for the modbus over IP device
	RTUHandle *RTUHandler `json:"rtuHandle,omitempty"`
	// traits this device supports
	Traits []DeviceTrait `json:"traits"`
	// meta data
	Metadata *traits.Metadata `json:"metadata,omitempty"`
}

type DeviceTrait struct {
	// smart core name of trait
	Name trait.Name `json:"name"`
	// where the device trait sits in side the modbus memory heap
	PDU *PDUAddress `json:"pdu"`
	// the address of the device trait
	Address uint16 `json:"address"`
	// quantity of the device trait
	Quantity uint16 `json:"quantity"`
	// interval to poll the device trait
	PollInterval *jsontypes.Duration `json:"pollInterval"`
	// scale factor to apply to the number held at address
	ScaleFactor float32 `json:"scaleFactor"`
}

type PDUAddress int

const (
	DiscreteInput PDUAddress = iota
	Coil
	InputRegister
	HoldingRegister
)

var (
	pduAddress = map[PDUAddress]string{
		DiscreteInput:   "DiscreteInput",
		Coil:            "Coil",
		InputRegister:   "InputRegister",
		HoldingRegister: "HoldingRegister",
	}

	addressPdu = map[string]PDUAddress{
		"DiscreteInput":   DiscreteInput,
		"Coil":            Coil,
		"InputRegister":   InputRegister,
		"HoldingRegister": HoldingRegister,
	}
)

func (p *PDUAddress) MarshalJSON() ([]byte, error) {
	if str, ok := pduAddress[*p]; ok {
		return json.Marshal(str)
	}

	return nil, fmt.Errorf("%d pdu address is not valid", *p)
}

func (p *PDUAddress) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if pdu, ok := addressPdu[str]; ok {
		*p = pdu
		return nil
	}

	return fmt.Errorf("%s pdu address is not valid", str)
}
