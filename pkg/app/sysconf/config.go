package sysconf

import (
	"os"

	"github.com/vanti-dev/sc-bos/pkg/auth/policy"
	"github.com/vanti-dev/sc-bos/pkg/auto"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/system"
	"go.uber.org/zap"
)

// Load loads into dst any user supplied config from json files and CLI arguments.
func Load(dst *Config) error {
	if err := LoadAllFromJSON(dst); err != nil {
		return err
	}
	if _, err := LoadFromArgs(dst, os.Args[1:]...); err != nil {
		return err
	}
	if err := LoadFromDataDirJSON(dst); err != nil {
		return err
	}

	// do any post processing
	dst.CertConfig = dst.CertConfig.FillDefaults()

	return nil
}

// Config configures how the controller should run.
type Config struct {
	ConfigDirs  []string `json:"-"` // Dirs we look in for system config files. Config in DataDir is always loaded and will have higher priority.
	ConfigFiles []string `json:"-"` // Filenames we load in ConfigDirs for system config

	Logger      zap.Config `json:"logger,omitempty"`
	ListenGRPC  string     `json:"listenGrpc,omitempty"`
	ListenHTTPS string     `json:"listenHttps,omitempty"`

	DataDir       string `json:"dataDir,omitempty"`       // defaults to .data/controller
	StaticDir     string `json:"staticDir,omitempty"`     // hosts static files from this directory over HTTP if StaticDir is non-empty
	AppConfigFile string `json:"appConfigFile,omitempty"` // defaults to app.conf.json
	CertConfig    Certs  `json:"certs,omitempty"`

	Systems map[string]system.RawConfig `json:"systems,omitempty"`

	Policy        policy.Policy `json:"-"` // Override the policy used for RPC calls. Defaults to policy.Default
	DisablePolicy bool          `json:"-"` // Unsafe, disables any policy checking for the server. Can't be set by json config.

	DriverFactories map[string]driver.Factory `json:"-"` // keyed by driver name
	AutoFactories   map[string]auto.Factory   `json:"-"` // keyed by automation type
	SystemFactories map[string]system.Factory `json:"-"` // keyed by system type
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

func Default() Config {
	config := Config{
		ConfigDirs:  []string{},
		ConfigFiles: []string{"system.conf.json", "system.json"},

		Logger:      zap.NewDevelopmentConfig(),
		ListenGRPC:  ":23557",
		ListenHTTPS: ":443",

		DataDir:       ".data/controller",
		StaticDir:     "",
		AppConfigFile: "app.conf.json",

		CertConfig: Certs{
			KeyFile:      "grpc.key.pem",
			CertFile:     "grpc.cert.pem",
			RootsFile:    "grpc.roots.pem",
			HTTPCert:     false,
			HTTPKeyFile:  "", // while these have defaults, we can't specify them and still have the "turn on if specified" feature
			HTTPCertFile: "",
		},
		Policy:        nil,
		DisablePolicy: false,
	}
	config.Logger.DisableStacktrace = true // because it's annoying
	return config
}

func (c Certs) FillDefaults() Certs {
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
