package helvarnet

import (
	"context"
	"net"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/lightpb"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensorpb"
	"github.com/vanti-dev/sc-bos/pkg/driver"
	"github.com/vanti-dev/sc-bos/pkg/driver/helvarnet/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task/service"
)

const (
	DriverName = "helvarnet"
)

var Factory driver.Factory = factory{}

type factory struct{}

type Driver struct {
	*service.Service[config.Root]
	announcer node.Announcer
	logger    *zap.Logger
	clients   map[string]*tcpClient
}

func (f factory) New(services driver.Services) service.Lifecycle {
	logger := services.Logger.Named(DriverName)

	d := &Driver{
		logger:    logger,
		announcer: services.Node,
	}

	d.Service = service.New(
		service.MonoApply(d.applyConfig),
		service.WithParser[config.Root](config.ParseConfig),
		service.WithRetry[config.Root](service.RetryWithLogger(func(logCtx service.RetryContext) {
			logCtx.LogTo("applyConfig", logger)
		}), service.RetryWithMinDelay(10*time.Second)),
	)

	return d
}

func (d *Driver) applyConfig(ctx context.Context, cfg config.Root) error {

	announcer, undo := node.AnnounceScope(d.announcer)
	grp, ctx := errgroup.WithContext(ctx)
	d.clients = make(map[string]*tcpClient)

	for _, l := range cfg.LightingGroups {
		if _, ok := d.clients[l.IpAddress]; !ok {
			tcpAddr, err := net.ResolveTCPAddr("tcp", l.IpAddress+*cfg.Port)
			if err != nil {
				return err
			}
			d.clients[l.IpAddress] = newTcpClient(tcpAddr, d.logger, &cfg)
		}

		if l.GroupNumber == nil {
			d.logger.Warn("lighting group without group number", zap.String("name", l.Name))
			continue
		}
		lightingGroup := newLightingGroup(d.clients[l.IpAddress], d.logger, l, *l.GroupNumber)
		// try to get the last scene on a restart of the area controller
		lightingGroup.getLastScene()
		err := lightingGroup.getSceneNames()
		if err != nil {
			d.logger.Error("getSceneNames error", zap.String("error", err.Error()))
		}
		announcer.Announce(l.Name,
			node.HasTrait(trait.Light,
				node.WithClients(lightpb.WrapApi(lightingGroup)),
				node.WithClients(lightpb.WrapInfo(lightingGroup))),
			node.HasMetadata(l.Meta))
	}

	for _, l := range cfg.Lights {
		if _, ok := d.clients[l.IpAddress]; !ok {
			tcpAddr, err := net.ResolveTCPAddr("tcp", l.IpAddress+*cfg.Port)
			if err != nil {
				return err
			}
			d.clients[l.IpAddress] = newTcpClient(tcpAddr, d.logger, &cfg)
		}
		lum := newLight(d.clients[l.IpAddress], d.logger, l)

		announcer.Announce(l.Name,
			node.HasTrait(trait.Light,
				node.WithClients(lightpb.WrapApi(lum))),
			node.HasTrait(statuspb.TraitName,
				node.WithClients(gen.WrapStatusApi(lum))),
			node.HasMetadata(l.Meta))
		grp.Go(func() error {
			return lum.runHealthCheck(ctx, cfg.RefreshStatus.Duration)
		})
	}

	for _, pir := range cfg.Pirs {
		if _, ok := d.clients[pir.IpAddress]; !ok {
			tcpAddr, err := net.ResolveTCPAddr("tcp", pir.IpAddress+*cfg.Port)
			if err != nil {
				return err
			}
			d.clients[pir.IpAddress] = newTcpClient(tcpAddr, d.logger, &cfg)
		}
		p := newPir(d.clients[pir.IpAddress], d.logger, pir)
		announcer.Announce(pir.Name,
			node.HasTrait(trait.OccupancySensor,
				node.WithClients(occupancysensorpb.WrapApi(p))),
			node.HasMetadata(pir.Meta))
		grp.Go(func() error {
			return p.runUpdateState(ctx, cfg.RefreshStatus.Duration)
		})
	}

	go func() {
		err := grp.Wait()
		for _, client := range d.clients {
			if client.conn != nil {
				client.close()
			}
		}
		undo()
		if err != nil {
			d.logger.Error("run error", zap.String("error", err.Error()))
		}
	}()
	return nil
}
