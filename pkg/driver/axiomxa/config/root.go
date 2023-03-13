package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

const DefaultPort = 60001

type Root struct {
	driver.BaseConfig
	HTTP         *HTTP        `json:"http,omitempty"`
	MessagePorts MessagePorts `json:"messagePorts,omitempty"`
	Database     *Database    `json:"database,omitempty"`

	Devices []Device `json:"devices,omitempty"`

	QR *QR `json:"qr,omitempty"`
}

func ReadBytes(data []byte) (root Root, err error) {
	err = json.Unmarshal(data, &root)
	if err != nil {
		return
	}
	if root.MessagePorts.Bind == "" {
		root.MessagePorts.Bind = fmt.Sprintf(":%d", DefaultPort)
	}
	return
}

type Device struct {
	Name string // Smart Core name
	// NetworkDesc is the human specified name identifying an Axiom controller.
	NetworkDesc string
	// DeviceDesc is the human specified name identifying a card reader managed by NetworkDesc.
	DeviceDesc string
	// UDMITopicPrefix is used for telemetry and config when using UDMI.
	// Defaults to Name.
	UDMITopicPrefix string
}

type HTTP struct {
	BaseURL      string     `json:"baseUrl,omitempty"`
	Username     string     `json:"username,omitempty"`
	Password     string     `json:"password,omitempty"`
	PasswordFile string     `json:"passwordFile,omitempty"`
	TLS          *TLSClient `json:"tls,omitempty"`
}

func (h *HTTP) Credentials() (username, password string, err error) {
	username = h.Username
	password = h.Password
	if password != "" {
		return
	}
	if h.PasswordFile != "" {
		var passFileBody []byte
		passFileBody, err = os.ReadFile(h.PasswordFile)
		if err != nil {
			return
		}
		password = strings.TrimSpace(string(passFileBody))
		return
	}

	return // no password I guess
}

type TLSClient struct {
	Disabled           bool   `json:"disabled,omitempty"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify,omitempty"`
	Roots              string `json:"roots,omitempty"` // PEM encoded certs
	RootsFile          string `json:"rootsFile,omitempty"`
}

func (t *TLSClient) TLSConfig() (*tls.Config, error) {
	if t == nil {
		return &tls.Config{}, nil // default to default tls config
	}
	if t.Disabled {
		return nil, nil
	}
	dst := &tls.Config{}
	if t.InsecureSkipVerify {
		dst.InsecureSkipVerify = true
	}
	if t.Roots != "" {
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM([]byte(t.Roots)) {
			return nil, errors.New("tls.roots is present but contains no certificates, it should be a PEM encoded list of CERTIFICATE blocks")
		}
		dst.RootCAs = pool
	} else if t.RootsFile != "" {
		roots, err := os.ReadFile(t.RootsFile)
		if err != nil {
			return nil, err
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(roots) {
			return nil, errors.New("tls.rootsFile is present but contains no certificates, it should be a PEM encoded list of CERTIFICATE blocks")
		}
		dst.RootCAs = pool
	}

	return dst, nil
}

type Database struct {
	DSN          string `json:"dsn,omitempty"`
	PasswordFile string `json:"passwordFile,omitempty"`
}

type QR struct {
	AccessLevel uint               `json:"accessLevel,omitempty"`
	ExpireAfter jsontypes.Duration `json:"expireAfter,omitempty"`
}
