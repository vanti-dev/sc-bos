package sysconf

import (
	"flag"
	"fmt"
	"path"
	"strings"
)

// LoadFromArgs populates fields in dst from the given command line arguments.
// Args should not include the program name, just like flag.FlagSet.Parse.
// Returns an error if the flags failed to parse or unknown flags exist.
func LoadFromArgs(dst *Config, args ...string) ([]string, error) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Var(&sysConfArg{dst}, "sysconf", "path to system config file")
	fs.Var(appConfArg{dst}, "appconf", "path to application config file")
	fs.StringVar(&dst.DataDir, "data-dir", dst.DataDir, "path to local data storage directory")
	fs.StringVar(&dst.ListenGRPC, "listen-grpc", dst.ListenGRPC, "address (host:port) to host a Smart Core gRPC server on")
	fs.StringVar(&dst.ListenHTTPS, "listen-https", dst.ListenHTTPS, "address (host:port) to host a HTTPS server on")
	fs.Var(disablePolicy{dst}, "insecure-disable-policy", "Deprecated. Equivalent to --policy-mode=off")
	fs.Var(&dst.PolicyMode, "policy-mode", `Configure how requests are compared against the authorization policy.
One of:
 - on (default): Only permit requests that pass the policy
 - off: Insecure! Permit all requests regardless of any policy or authentication
 - check: Insecure! If requests have a token or client certificate, check them 
   against the policy, otherwise permit them`)

	// todo: add support for staticHosting config via flag

	err := fs.Parse(args)
	return fs.Args(), err
}

type sysConfArg struct {
	dst *Config
}

func (a sysConfArg) String() string {
	return strings.Join(a.dst.ConfigFiles, ",")
}

func (a sysConfArg) Set(s string) error {
	str := strings.Split(s, ",")
	a.dst.ConfigDirs = []string{path.Base("")} // make config relative to current location
	for _, f := range str {
		a.dst.ConfigFiles = append(a.dst.ConfigFiles, f)
	}
	return nil
}

type appConfArg struct {
	dst *Config
}

func (a appConfArg) String() string {
	return strings.Join(a.dst.AppConfig, ",")
}

func (a appConfArg) Set(s string) error {
	str := strings.Split(s, ",")
	for _, f := range str {
		a.dst.AppConfig = append(a.dst.AppConfig, path.Join(path.Base(""), f))
	}
	return nil
}

// disablePolicy maps from the legacy DisablePolicy flag to the new PolicyMode flag.
// todo: remove this once all clients have migrated to the new PolicyMode flag.
type disablePolicy struct {
	dst *Config
}

func (d disablePolicy) String() string {
	if d.dst != nil && d.dst.PolicyMode == PolicyOff {
		return "true"
	}
	return "false"
}

func (d disablePolicy) Set(s string) error {
	if s == "true" {
		d.dst.PolicyMode = PolicyOff
	}
	return nil
}

func (d disablePolicy) IsBoolFlag() bool {
	return true
}

func (pm *PolicyMode) String() string {
	if pm == nil {
		return string(PolicyOn)
	}
	return string(*pm)
}

func (pm *PolicyMode) Set(s string) error {
	switch PolicyMode(s) {
	case PolicyOn:
		*pm = PolicyOn
	case PolicyOff:
		*pm = PolicyOff
	case PolicyCheck:
		*pm = PolicyCheck
	default:
		return fmt.Errorf("supported [on,off,check]")
	}
	return nil
}
