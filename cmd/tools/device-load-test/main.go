// A tool to load test the server by simulating the request pattern of the Ops UI loading a layered graphic with
// many devices.
//
// It reads UI config files to determine which devices and traits to load, then connects to the server via gRPC
// and pulls data for all those device traits concurrently, reporting the time taken to receive the first value
// for each trait.
// This provides a lower bound on the time the UI would take to load and populate such a graphic with data for all
// its devices.
//
// Unlike the Ops UI, this tool uses native gRPC rather than gRPC-web, and does not have to render anything,
// so we expect it to be faster than the UI.
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

var (
	flagServer string
	flagConfig string
	flagPage   string
)

func init() {
	flag.StringVar(&flagServer, "server", "127.0.0.1:23557", "server gRPC address (host:port)")
	flag.StringVar(&flagConfig, "config", "", "path to UI config file")
	flag.StringVar(&flagPage, "page", "", "page path to load, under /ops/overview (e.g. building/Ground Floor)")
}

func run(ctx context.Context) error {
	flag.Parse()
	if flagConfig == "" {
		return fmt.Errorf("missing flag -config")
	}
	if flagPage == "" {
		return fmt.Errorf("missing flag -page")
	}

	allLayers, err := loadUIConfigLayers(flagConfig)
	if err != nil {
		return fmt.Errorf("failed to load UI config: %w", err)
	}

	var pageLayers []layerConfig
	for p, l := range allLayers {
		if p == flagPage {
			pageLayers = l
		}
	}
	if len(pageLayers) == 0 {
		return fmt.Errorf("no layers found for page %q", flagPage)
	}

	var dts []deviceTrait
	for _, layer := range pageLayers {
		layerTraits, err := layerDeviceTraits(layer)
		if err != nil {
			return fmt.Errorf("failed to get device traits for layer: %w", err)
		}
		dts = append(dts, layerTraits...)
	}
	slices.SortFunc(dts, compareDeviceTrait)

	conn, err := grpc.NewClient(flagServer, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("failed to close gRPC connection: %v", err)
		}
	}()

	// to track which device traits are still remaining to be pulled
	remaining := map[deviceTrait]struct{}{}
	for _, dt := range dts {
		remaining[dt] = struct{}{}
	}

	var grp sync.WaitGroup
	defer grp.Wait()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	completed := make(chan pullResult)

	start := time.Now()
	for _, dt := range dts {
		grp.Go(func() {
			pull(ctx, completed, conn, dt)
		})
	}

	for len(remaining) > 0 {
		select {
		case res := <-completed:
			delete(remaining, res.deviceTrait)
			fmt.Printf("first value: %s\t\t\t%s\t\tat %v\n", res.Name, res.Trait.String(), res.Time)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	cancel()
	elapsed := time.Since(start)
	log.Printf("completed pulling first value of %d device traits in %s", len(dts), elapsed)

	return nil
}
