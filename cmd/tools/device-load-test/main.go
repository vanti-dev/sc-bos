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
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/benchmark/latency"
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
	flagServer  string
	flagConfig  string
	flagPage    string
	flagLatency time.Duration
	flagVerbose bool
)

func init() {
	flag.StringVar(&flagServer, "server", "127.0.0.1:23557", "server gRPC address (host:port)")
	flag.StringVar(&flagConfig, "config", "", "path to UI config file")
	flag.StringVar(&flagPage, "page", "", "page path to load, under /ops/overview (e.g. building/Ground Floor)")
	flag.DurationVar(&flagLatency, "latency", 0, "artificial latency to add on stream open")
	flag.BoolVar(&flagVerbose, "verbose", false, "enable verbose logging")
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

	counter := &channelCounter{}
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})),
		grpc.WithUnaryInterceptor(counter.UnaryInterceptor()),
		grpc.WithStreamInterceptor(counter.StreamInterceptor()),
	}
	if flagLatency > 0 {
		log.Printf("adding artificial latency of %s", flagLatency)
		network := latency.Network{
			Latency: flagLatency,
		}
		baseDialer := &net.Dialer{}
		latencyDialer := network.ContextDialer(baseDialer.DialContext)
		dialOptions = append(dialOptions,
			grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
				return latencyDialer(ctx, "tcp", addr)
			}),
		)
	}
	conn, err := grpc.NewClient(flagServer, dialOptions...)
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

	grp.Go(func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				printCounts(counter)
			}
		}
	})
	start := time.Now()
	for _, dt := range dts {
		grp.Go(func() {
			pull(ctx, completed, conn, dt)
		})
	}
	fmt.Printf("%d device traits subscribed to.\n", len(dts))

	for len(remaining) > 0 {
		select {
		case res := <-completed:
			if res.Err != nil {
				fmt.Printf("failed to pull %s (%s): %v\n", res.Name, res.Trait, res.Err)
			} else if _, ok := remaining[res.deviceTrait]; ok {
				if flagVerbose {
					fmt.Printf("first value: %s\t\t\t%s\t\tat %v\n", res.Name, res.Trait.String(), res.Time)
				}
			} else {
				if flagVerbose {
					fmt.Printf("extra value: %s\t\t\t%s\t\tat %v\n", res.Name, res.Trait.String(), res.Time)
				}
			}
			delete(remaining, res.deviceTrait)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	cancel()
	elapsed := time.Since(start)
	log.Printf("completed pulling first value of %d device traits in %s", len(dts), elapsed)

	maxCounts := counter.MaxCounts()
	log.Printf("max open channels %d, unary: %d, stream: %d", maxCounts.Channel, maxCounts.Unary, maxCounts.Stream)

	return nil
}

func printCounts(counts *channelCounter) {
	current := counts.CurrentCounts()
	log.Printf("current open channels %d, unary: %d, stream: %d", current.Channel, current.Unary, current.Stream)
}
