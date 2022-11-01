//go:build !notc3dali

package tc3dali

import (
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/bridge"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/dali"
	"github.com/vanti-dev/twincat3-ads-go/pkg/ads"
	"github.com/vanti-dev/twincat3-ads-go/pkg/adsdll"
	"github.com/vanti-dev/twincat3-ads-go/pkg/device"
	"go.uber.org/zap"
)

func newBusBuilder(adsConfig ADSConfig) (busBuilder, error) {
	port, err := adsdll.Connect()
	if err != nil {
		return nil, err
	}

	return &bridgeBusBuilder{
		port: port,
		addr: ads.Addr{
			NetId: ads.NetId(adsConfig.NetID),
			Port:  adsConfig.Port,
		},
	}, nil
}

type bridgeBusBuilder struct {
	port ads.Port
	addr ads.Addr
}

func (bb *bridgeBusBuilder) buildBus(config BusConfig, logger *zap.Logger) (dali.Dali, error) {
	dev, err := device.Open(bb.port, bb.addr)
	if err != nil {
		logger.Error("failed to connect to ADS PLC Device", zap.Error(err),
			zap.Uint8s("netID", bb.addr.NetId[:]), zap.Uint16("port", bb.addr.Port))
		return nil, err
	}

	bridgeConfig := &bridge.Config{
		Device:                  dev,
		Logger:                  logger,
		BridgeFBName:            config.BridgePrefix + BridgeSuffix,
		ResponseMailboxName:     config.BridgePrefix + ResponseMailboxSuffix,
		NotificationMailboxName: config.BridgePrefix + NotificationMailboxSuffix,
	}
	bus, err := bridgeConfig.Connect()
	if err != nil {
		logger.Error("DALI bus bridge initialisation failure", zap.Error(err),
			zap.String("busName", config.Name),
			zap.String("prefix", config.BridgePrefix))
		return nil, err
	}
	return bus, nil
}
