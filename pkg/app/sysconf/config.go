// Package sysconf provides system level configuration.
package sysconf

import (
	"os"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/app/http"
	"github.com/vanti-dev/sc-bos/pkg/auth/policy"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"github.com/vanti-dev/sc-bos/pkg/util/netutil"
	"github.com/vanti-dev/sc-bos/pkg/zone"
)

// Load loads into dst any user supplied config from json files and CLI arguments.
func Load(dst *Config) error {
	if err := LoadAllFromJSON(dst); err != nil {
		return err
	}
	if _, err := LoadFromArgs(dst, os.Args[1:]...); err != nil {
		return err
	}
	if err := LoadFromConfigDirJSON(dst); err != nil {
		return err
	}

	// do any post processing
	dst.CertConfig = dst.CertConfig.FillDefaults()

	return nil
}

// Config configures how the controller should run.
type Config struct {
	ConfigDirs  []string `json:"-"` // Dirs we look in for system config files. Config in ConfigDir is always loaded and will have higher priority.
	ConfigFiles []string `json:"-"` // Filenames we load in ConfigDirs for system config

	// The smart core name of the controller.
	// Can be overridden by app config.
	Name string `json:"name,omitempty"`

	Logger      *zap.Config `json:"logger,omitempty"`
	ListenGRPC  string      `json:"listenGrpc,omitempty"`
	ListenHTTPS string      `json:"listenHttps,omitempty"`
	// FooAddr are preferred IP/host others use to connect to us.
	// Defaults to netutil.PublicAddress
	GRPCAddr string `json:"grpcAddr,omitempty"`
	HTTPAddr string `json:"httpAddr,omitempty"`

	ConfigDir     string                     `json:"configDir,omitempty"` // defaults to .data/controller
	DataDir       string                     `json:"dataDir,omitempty"`   // defaults to .data/controller
	StaticHosting []http.StaticHostingConfig `json:"staticHosting"`
	AppConfigFile string                     `json:"appConfigFile,omitempty"` // defaults to app.conf.json
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
		ConfigDirs:  []string{},
		ConfigFiles: []string{"system.conf.json", "system.json"},

		Logger:      &logConf,
		ListenGRPC:  ":23557",
		ListenHTTPS: ":443",

		Cors: http.CorsConfig{
			DebugMode: false,
			// todo: this should really default to the default host
			CorsOrigins: []string{"*"},
		},
		StaticHosting: []http.StaticHostingConfig{},
		ConfigDir:     ".conf/controller",
		DataDir:       ".data/controller",
		AppConfigFile: "app.conf.json",

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

	if localIP, err := netutil.OutboundAddr(); err == nil {
		config.GRPCAddr = localIP.String()
		config.HTTPAddr = localIP.String()
	}

	return config
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
