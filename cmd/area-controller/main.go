package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/vanti-dev/bsp-ew/internal/app"
	"github.com/vanti-dev/bsp-ew/internal/testapi"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"google.golang.org/grpc"
)

var (
	flagListenGRPC  string
	flagListenHTTPS string
	flagDataDir     string
	flagStaticDir   string
)

func init() {
	flag.StringVar(&flagListenGRPC, "listen-grpc", ":23557", "address (host:port) to host a Smart Core gRPC server on")
	flag.StringVar(&flagListenHTTPS, "listen-https", ":443", "address (host:port) to host a HTTPS server on")
	flag.StringVar(&flagDataDir, "data-dir", ".data/area-controller-01", "path to local data storage directory")
	flag.StringVar(&flagStaticDir, "static-dir", "ui/dist", "path for HTTP static resources")
}

func main() {
	flag.Parse()
	c := &app.Controller{
		DataDir:     flagDataDir,
		ListenGRPC:  flagListenGRPC,
		ListenHTTPS: flagListenHTTPS,
		Routes: map[string]http.Handler{
			"/": http.FileServer(http.Dir(flagStaticDir)),
		},
		RegisterServices: func(server *grpc.Server) {
			gen.RegisterTestApiServer(server, testapi.NewAPI())
		},
	}

	os.Exit(app.RunUntilInterrupt(c.Run))
}
