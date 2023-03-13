package sysconf

import (
	"flag"

	"github.com/vanti-dev/sc-bos/pkg/app/http"
)

// LoadFromArgs populates fields in dst from the given command line arguments.
// Args should not include the program name, just like flag.FlagSet.Parse.
// Returns an error if the flags failed to parse or unknown flags exist.
func LoadFromArgs(dst *Config, args ...string) ([]string, error) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&dst.ListenGRPC, "listen-grpc", dst.ListenGRPC, "address (host:port) to host a Smart Core gRPC server on")
	fs.StringVar(&dst.ListenHTTPS, "listen-https", dst.ListenHTTPS, "address (host:port) to host a HTTPS server on")
	fs.StringVar(&dst.DataDir, "data-dir", dst.DataDir, "path to local data storage directory")
	var static string
	fs.StringVar(&static, "static-dir", "", "(optional) path to directory to host static files over HTTP")
	fs.BoolVar(&dst.DisablePolicy, "insecure-disable-policy", dst.DisablePolicy, "Insecure! Disable checking requests against the security policy. This option opens up the server to any request.")

	err := fs.Parse(args)
	dst.StaticHosting = []http.StaticHostingConfig{{
		FilePath: static,
		Path:     "/",
	}}
	return fs.Args(), err
}
