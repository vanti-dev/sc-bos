package helvarnet

import (
	"context"
	"net"
	"time"

	"github.com/timshannon/bolthold"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/smart-core-os/sc-bos/pkg/driver"
	"github.com/smart-core-os/sc-bos/pkg/driver/helvarnet/config"
	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/emergencylightpb"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/statuspb"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/udmipb"
	"github.com/smart-core-os/sc-bos/pkg/node"
	"github.com/smart-core-os/sc-bos/pkg/task/service"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/lightpb"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensorpb"
)

const (
	DriverName = "helvarnet"
)

var Factory driver.Factory = factory{}

type factory struct{}

type Driver struct {
	*service.Service[config.Root]
	announcer *node.ReplaceAnnouncer
	logger    *zap.Logger
	clients   map[string]*tcpClient
	database  *bolthold.Store
}

func (f factory) New(services driver.Services) service.Lifecycle {
	logger := services.Logger.Named(DriverName)

	d := &Driver{
		logger:    logger,
		announcer: node.NewReplaceAnnouncer(services.Node),
		database:  services.Database,
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

	rootAnnouncer := d.announcer.Replace(ctx)
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
		rootAnnouncer.Announce(l.Name,
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
		lum := newLight(d.clients[l.IpAddress], d.logger, l, d.database, false)

		rootAnnouncer.Announce(l.Name,
			node.HasTrait(trait.Light,
				node.WithClients(lightpb.WrapApi(lum))),
			node.HasTrait(statuspb.TraitName,
				node.WithClients(gen.WrapStatusApi(lum))),
			node.HasTrait(udmipb.TraitName,
				node.WithClients(gen.WrapUdmiService(lum))),
			node.HasMetadata(l.Meta))
		grp.Go(func() error {
			return lum.queryDevice(ctx, cfg.RefreshStatus.Duration)
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
		rootAnnouncer.Announce(pir.Name,
			node.HasTrait(trait.OccupancySensor,
				node.WithClients(occupancysensorpb.WrapApi(p))),
			node.HasTrait(udmipb.TraitName,
				node.WithClients(gen.WrapUdmiService(p))),
			node.HasMetadata(pir.Meta))
		grp.Go(func() error {
			return p.runUpdateState(ctx, cfg.RefreshOccupancy.Duration)
		})
	}

	for _, em := range cfg.EmergencyLights {
		if _, ok := d.clients[em.IpAddress]; !ok {
			tcpAddr, err := net.ResolveTCPAddr("tcp", em.IpAddress+*cfg.Port)
			if err != nil {
				return err
			}
			d.clients[em.IpAddress] = newTcpClient(tcpAddr, d.logger, &cfg)
		}
		emergencyLight := newLight(d.clients[em.IpAddress], d.logger, em, d.database, true)
		err := emergencyLight.loadTestResults()
		if err != nil {
			d.logger.Error("loadTestResults error", zap.Error(err))
		}

		rootAnnouncer.Announce(em.Name,
			node.HasTrait(trait.Light,
				node.WithClients(lightpb.WrapApi(emergencyLight))),
			node.HasTrait(statuspb.TraitName,
				node.WithClients(gen.WrapStatusApi(emergencyLight))),
			node.HasTrait(emergencylightpb.TraitName,
				node.WithClients(gen.WrapEmergencyLightApi(emergencyLight))),
			node.HasTrait(udmipb.TraitName,
				node.WithClients(gen.WrapUdmiService(emergencyLight))),
			node.HasMetadata(em.Meta))
		grp.Go(func() error {
			return emergencyLight.queryDevice(ctx, cfg.RefreshStatus.Duration)
		})
	}

	go func() {
		err := grp.Wait()
		for _, client := range d.clients {
			if client.conn != nil {
				client.close()
			}
		}
		if err != nil {
			d.logger.Error("run error", zap.String("error", err.Error()))
		}
	}()
	return nil
}
