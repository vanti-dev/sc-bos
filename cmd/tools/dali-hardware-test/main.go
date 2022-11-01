//go:build !notc3dali

package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/bridge"
	"github.com/vanti-dev/twincat3-ads-go/pkg/ads"
	"github.com/vanti-dev/twincat3-ads-go/pkg/adsdll"
	"github.com/vanti-dev/twincat3-ads-go/pkg/device"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	ADS       tc3dali.ADSConfig `json:"ads"`
	Prefixes  []string          `json:"prefixes"`
	LowLevel  uint8             `json:"lowLevel"`
	HighLevel uint8             `json:"highLevel"`
	Delay     string            `json:"delay"`
}

var (
	flagConfig string
)

func init() {
	flag.StringVar(&flagConfig, "config", "config.json", "path to config file")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := run(ctx)
	if errors.Is(err, context.Canceled) {
		_, _ = fmt.Fprintln(os.Stderr, "Exit due to Interrupt")
	} else if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "FATAL ERROR: %s\n", err.Error())
	}
}

func run(ctx context.Context) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	delay, err := time.ParseDuration(config.Delay)
	if err != nil {
		return fmt.Errorf("invalid delay duration value %q: %w", config.Delay, err)
	}

	// connect to PLC device
	port, err := adsdll.Connect()
	if err != nil {
		return fmt.Errorf("connect to ADS router: %w", err)
	}
	defer func() {
		closeErr := port.Close()
		if closeErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ERROR: failed to close ADS port: %s\n", err.Error())
		}
	}()
	dev, err := device.Open(port, ads.Addr{
		NetId: ads.NetId(config.ADS.NetID),
		Port:  config.ADS.Port,
	})
	if err != nil {
		return fmt.Errorf("connect to ADS %v:%d: %w", config.ADS.NetID, config.ADS.Port, err)
	}

	// initialise bridge connections
	bridges := make([]bridge.Dali, 0, len(config.Prefixes))
	for _, prefix := range config.Prefixes {
		bridgeConfig := &bridge.Config{
			Device:                  dev,
			Logger:                  zap.NewNop(),
			BridgeFBName:            prefix + "_bridge",
			ResponseMailboxName:     prefix + "_response",
			NotificationMailboxName: prefix + "_notification",
		}
		b, err := bridgeConfig.Connect()
		if err != nil {
			return fmt.Errorf("init bridge %q: %w", prefix, err)
		}
		bridges = append(bridges, b)
	}

	// prepare by setting all buses to low level
	err = sendBroadcastLevel(ctx, config.LowLevel, bridges...)
	if err != nil {
		return err
	}

	// when we stop, attempt to turn all the lights off
	defer func() {
		// need a new context because the main one may have been cancelled by this point
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := sendBroadcastLevel(ctx, 0, bridges...)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ERROR: failed to turn off some lights: %s\n", err.Error())
		}
	}()

	for ctx.Err() == nil {
		for i, b := range bridges {
			fmt.Printf("highlighting bus %d - %q\n", i, config.Prefixes[i])

			err = sendBroadcastLevel(ctx, config.HighLevel, b)
			if err != nil {
				return err
			}

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}

			err = sendBroadcastLevel(ctx, config.LowLevel, b)
			if err != nil {
				return err
			}
		}

		// pause between iterations, so that if there is only one bus it still spends some time off
		fmt.Println("no highlight")
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return ctx.Err()
}

func loadConfig() (config Config, err error) {
	raw, err := os.ReadFile(flagConfig)
	if err != nil {
		return config, fmt.Errorf("read config from %q: %w", flagConfig, err)
	}

	// set defaults
	config.Delay = "2s"
	config.LowLevel = 50
	config.HighLevel = 200

	err = json.Unmarshal(raw, &config)
	if err != nil {
		return config, fmt.Errorf("decode config: %w", err)
	}
	return
}

func sendBroadcastLevel(ctx context.Context, level uint8, bridges ...bridge.Dali) error {
	var group errgroup.Group
	for _, b := range bridges {
		b := b
		group.Go(func() error {
			_, err := b.ExecuteCommand(ctx, bridge.Request{
				Command:     bridge.DirectArcPowerControl,
				AddressType: bridge.Broadcast,
				Data:        level,
			})
			return err
		})
	}

	return group.Wait()
}
