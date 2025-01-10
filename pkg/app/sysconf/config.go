// Package sysconf provides system level configuration.
package sysconf

import (
	"os"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/app/http"
	"github.com/vanti-dev/sc-bos/pkg/auth/policy"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/block"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/util/netutil"
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

// Load loads into dst any user supplied config from json files and CLI arguments. CLI arguments take precedence.
func Load(dst *Config) error {
	// We call LoadFromArgs twice because args can be used to specify config file paths,
	// but also args should override config values specified in json files.

	if _, err := LoadFromArgs(dst, os.Args[1:]...); err != nil {
		return err
	}
	if err := LoadAllFromJSON(dst); err != nil {
		return err
	}
	if err := LoadFromConfigFilesJSON(dst); err != nil {
		return err
	}
	if _, err := LoadFromArgs(dst, os.Args[1:]...); err != nil {
		return err
	}

	// do any post processing
	dst.Normalize()

	return nil
}

// Config configures how the controller should run.
type Config struct {
	ConfigDirs  []string `json:"-"` // Dirs to look in for system config files. Defaults to [".conf"]
	ConfigFiles []string `json:"-"` // Filenames we load in ConfigDirs for system config. Defaults to ["system.conf.json", "system.json"]

	// The smart core name of the controller.
	// Can be overridden by app config.
	Name string `json:"name,omitempty"`

	Logger      *zap.Config `json:"logger,omitempty"`
	ListenGRPC  string      `json:"listenGrpc,omitempty"`
	ListenHTTPS string      `json:"listenHttps,omitempty"`
	// Preferred IP/host others use to connect to us.
	// Typically used when the controller constructs and shares its own address with others,
	// for example during enrollment or when producing download links.
	GRPCAddr string `json:"grpcAddr,omitempty"` // Defaults to netutil.OutboundAddr
	HTTPAddr string `json:"httpAddr,omitempty"` // Defaults to GRPCAddr

	SANs []string `json:"sans,omitempty"` // Subject Alternative Names for the self-signed cert

	AppConfig []string `json:"appConfig,omitempty"` // defaults to [".conf/app.conf.json"]
	DataDir   string   `json:"dataDir,omitempty"`   // defaults to .data/

	StaticHosting []http.StaticHostingConfig `json:"staticHosting"`
	CertConfig    *Certs                     `json:"certs,omitempty"`
	Cors          http.CorsConfig            `json:"cors,omitempty"`

	Systems map[string]system.RawConfig `json:"systems,omitempty"`

	Policy     policy.Policy `json:"-"` // Override the policy used for RPC calls. Defaults to policy.Default
	PolicyMode PolicyMode    `json:"-"` // How to apply the policy. Unsafe and can disable security checks. Defaults to PolicyOn.

	DriverFactories map[string]driver.Factory `json:"-"` // keyed by driver name
	AutoFactories   map[string]auto.Factory   `json:"-"` // keyed by automation type
	SystemFactories map[string]system.Factory `json:"-"` // keyed by system type
	ZoneFactories   map[string]zone.Factory   `json:"-"` // keyed by zone type
}

// DriverConfigBlocks returns a map of driver type to a block list that describes the config for that driver.
// It scans DriverFactories for factories that implement the BlockSource interface.
// Drivers that do not implement BlockSource are not included in the output.
func (c *Config) DriverConfigBlocks() map[string][]block.Block {
	return extractConfigBlocks(c, c.DriverFactories)
}

// AutoConfigBlocks returns a map of automation type to a block list that describes the config for that automation.
// It scans AutoFactories for factories that implement the BlockSource interface.
// Automations that do not implement BlockSource are not included in the output.
func (c *Config) AutoConfigBlocks() map[string][]block.Block {
	return extractConfigBlocks(c, c.AutoFactories)
}

// ZoneConfigBlocks returns a map of zone type to a block list that describes the config for that zone.
// It scans ZoneFactories for factories that implement the BlockSource interface.
// Zone types that do not implement BlockSource are not included in the output.
func (c *Config) ZoneConfigBlocks() map[string][]block.Block {
	return extractConfigBlocks(c, c.ZoneFactories)
}

func extractConfigBlocks[Factory any](c *Config, factories map[string]Factory) map[string][]block.Block {
	blocks := make(map[string][]block.Block)
	for name, factory := range factories {
		switch source := any(factory).(type) {
		case BlockSource:
			blocks[name] = source.ConfigBlocks()
		case BlockSource2:
			blocks[name] = source.ConfigBlocks(c)
		}
	}
	return blocks
}

// Certs encapsulates different settings used for loading and present certificates to clients and servers.
type Certs struct {
	KeyFile   string `json:"keyFile,omitempty"`
	CertFile  string `json:"certFile,omitempty"`
	RootsFile string `json:"rootsFile,omitempty"`

	HTTPCert     bool   `json:"httpCert,omitempty"` // have the https stack (grpc-web and hosting) use different pki.Source from the grpc stack
	HTTPKeyFile  string `json:"httpKeyFile,omitempty"`
	HTTPCertFile string `json:"httpCertFile,omitempty"`
}

type PolicyMode string

const (
	PolicyOn    PolicyMode = "on"    // Always check requests against the policy.
	PolicyOff   PolicyMode = "off"   // Never check requests against the policy, allow all requests.
	PolicyCheck PolicyMode = "check" // Check requests against the policy if the request has a token or client cert.
)

func Default() Config {
	logConf := zap.NewDevelopmentConfig()
	config := Config{
		ConfigDirs:  []string{".conf"},
		ConfigFiles: []string{"system.conf.json", "system.json"},

		Logger:      &logConf,
		ListenGRPC:  ":23557",
		ListenHTTPS: ":443",

		AppConfig: []string{".conf/app.conf.json"},
		DataDir:   ".data",

		Cors: http.CorsConfig{
			DebugMode: false,
			// todo: this should really default to the default host
			CorsOrigins: []string{"*"},
		},
		StaticHosting: []http.StaticHostingConfig{},

		CertConfig: &Certs{
			KeyFile:      "grpc.key.pem",
			CertFile:     "grpc.cert.pem",
			RootsFile:    "grpc.roots.pem",
			HTTPCert:     false,
			HTTPKeyFile:  "", // while these have defaults, we can't specify them and still have the "turn on if specified" feature
			HTTPCertFile: "",
		},
		Policy:     nil,
		PolicyMode: PolicyOn,
	}
	config.Logger.DisableStacktrace = true // because it's annoying

	return config
}

// Normalize adjusts c to apply defaults that are based on the values of other fields.
// Normalize should be called explicitly if not using Load.
func (c *Config) Normalize() {
	if c.GRPCAddr == "" {
		if addr, err := netutil.OutboundAddr(); err == nil {
			c.GRPCAddr = addr.String()
		}
	}
	if c.HTTPAddr == "" {
		c.HTTPAddr = c.GRPCAddr
	}

	c.CertConfig = c.CertConfig.FillDefaults()
}

func (c *Certs) FillDefaults() *Certs {
	or := func(a *string, b string) {
		if *a == "" {
			*a = b
		}
	}

	// if the config specifies http key or cert file paths, assume they want to use it
	if c.HTTPKeyFile != "" || c.HTTPCertFile != "" {
		c.HTTPCert = true
	}
	or(&c.HTTPKeyFile, "https.key.pem")
	or(&c.HTTPCertFile, "https.cert.pem")

	return c
}

// BlockSource is an interface that can be implemented by a factory to provide a list of blocks that describe the config
// for that driver/automation/zone type.
// This is used to produce granular diffs for config changes.
type BlockSource interface {
	ConfigBlocks() []block.Block
}

// BlockSource2 is like BlockSource but allows the factory to receive the config as an argument.
type BlockSource2 interface {
	ConfigBlocks(cfg *Config) []block.Block
}
