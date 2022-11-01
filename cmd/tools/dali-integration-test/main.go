//go:build !notc3dali

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/bsp-ew/internal/driver"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/task"
	"github.com/vanti-dev/twincat3-ads-go/pkg/ads"
	"github.com/vanti-dev/twincat3-ads-go/pkg/adsdll"
	"github.com/vanti-dev/twincat3-ads-go/pkg/device"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var (
	flagAMSNetID  tc3dali.NetID
	flagADSPort   int
	flagBusPrefix string
)

func init() {
	flag.Func("ams-net-id", "TwinCAT 3 AMS NetID of the PLC instance to connect to", func(s string) error {
		parsed, err := tc3dali.ParseNetID(s)
		if err != nil {
			return err
		}
		flagAMSNetID = parsed
		return nil
	})
	flag.IntVar(&flagADSPort, "ads-port", 851, "TwinCAT 3 ADS Port of the PLC instance to connect to")
	flag.StringVar(&flagBusPrefix, "bus-prefix", "GVL_Bridges.bus_T1_1", "Bus variable prefix")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	err := run(ctx)
	cancel()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "FATAL ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	// TC3 setup
	port, err := adsdll.Connect()
	if err != nil {
		return fmt.Errorf("adsdll.Connect: %w", err)
	}

	logger.Info("connecting to TwinCAT", zap.String("netID", fmt.Sprint(flagAMSNetID)),
		zap.Uint16("port", uint16(flagADSPort)))
	dev, err := device.Open(port, ads.Addr{
		NetId: ads.NetId(flagAMSNetID),
		Port:  uint16(flagADSPort),
	})
	if err != nil {
		return fmt.Errorf("device.Open: %w", err)
	}

	// Driver setup
	services := driver.Services{
		Logger: logger.Named("tc3dali"),
		Node:   node.New("dali-integration-test"),
		Tasks:  &task.Group{},
	}
	busConfig := tc3dali.BusConfig{
		Name:         "dali-integration-test/bus/1",
		BridgePrefix: flagBusPrefix,
	}
	busTask := tc3dali.BusTask(busConfig, dev, services)

	// gRPC setup
	server := grpc.NewServer()
	defer server.Stop()
	services.Node.Register(server)
	listener := bufconn.Listen(1024 * 1024)
	go func() {
		if err := server.Serve(listener); err != nil {
			logger.Error("mock server stopped with error", zap.Error(err))
		}
		_ = listener.Close()
	}()
	conn, err := grpc.Dial("",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
	)
	if err != nil {
		return fmt.Errorf("grpc.Dial: %w", err)
	}

	// Run driver init routine
	err = task.Run(ctx, busTask)
	if err != nil {
		logger.Error("driver failed to initialise", zap.Error(err))
		return err
	}

	// use a gRPC client to do stuff
	lightClient := traits.NewLightApiClient(conn)
	for ctx.Err() == nil {
		_, err = lightClient.UpdateBrightness(ctx, &traits.UpdateBrightnessRequest{
			Name:       busConfig.Name,
			Brightness: &traits.Brightness{LevelPercent: 0},
		})

		time.Sleep(1)

		_, err = lightClient.UpdateBrightness(ctx, &traits.UpdateBrightnessRequest{
			Name:       busConfig.Name,
			Brightness: &traits.Brightness{LevelPercent: 100},
		})

		time.Sleep(1)
	}

	return nil
}
