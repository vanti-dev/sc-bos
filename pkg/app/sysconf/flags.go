package sysconf

import (
	"flag"
	"strings"
)

// LoadFromArgs populates fields in dst from the given command line arguments.
// Args should not include the program name, just like flag.FlagSet.Parse.
// Returns an error if the flags failed to parse or unknown flags exist.
func LoadFromArgs(dst *Config, args ...string) ([]string, error) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&dst.ListenGRPC, "listen-grpc", dst.ListenGRPC, "address (host:port) to host a Smart Core gRPC server on")
	var https string
	fs.StringVar(&https, "listen-https", dst.HttpConfig.Host+":"+dst.HttpConfig.Port, "address (host:port) to host a HTTPS server on")
	fs.StringVar(&dst.DataDir, "data-dir", dst.DataDir, "path to local data storage directory")
	fs.BoolVar(&dst.DisablePolicy, "insecure-disable-policy", dst.DisablePolicy, "Insecure! Disable checking requests against the security policy. This option opens up the server to any request.")

	err := fs.Parse(args)

	h := strings.Split(https, ":")
	dst.HttpConfig.Host = h[0]
	dst.HttpConfig.Port = h[1]

	return fs.Args(), err
}
