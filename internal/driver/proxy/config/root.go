package config

import (
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/bsp-ew/internal/driver"
)

// Root describes the configuration available to the proxy driver.
type Root struct {
	driver.BaseConfig

	// Nodes represent Smart Core nodes that this controller is proxying.
	Nodes []Node `json:"nodes"`
}

// Node is a networked Smart Core node, identified by its host.
type Node struct {
	Host string `json:"host,omitempty"` // for accepted values see grpc.Dial

	// TLS allows us to override the default enrollment managed TLS configuration.
	TLS TLS `json:"tls,omitempty"`

	// Name is the Smart Core name for the remote node.
	// Used when discovering traits.
	// If absent then the remote node should support empty/default names for Parent requests.
	Name string `json:"name,omitempty"`

	// Traits defines the exact named traits we proxy for this remote.
	// If absent or empty the remote will be inspected using the Parent trait including all found traits.
	//Traits []Trait `json:"traits,omitempty"` // todo: support manual traits
}

type TLS struct {
	// These override the default enrollment PKI and are useful for testing and if running without a management node.
	InsecureNoClientCert bool `json:"insecureNoClientCert,omitempty"` // don't present a client certificate when connecting to proxy servers
	InsecureSkipVerify   bool `json:"insecureSkipVerify,omitempty"`   // don't verify proxy server certificates
}

type Trait struct {
	Name  string     `json:"name,omitempty"`
	Trait trait.Name `json:"trait,omitempty"`
}
