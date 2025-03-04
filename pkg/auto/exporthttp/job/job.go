package job

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/auto/exporthttp/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/node"
)

var (
	// strict:compile
	_ Job = (*OccupancyJob)(nil)
	_ Job = (*TemperatureJob)(nil)
	_ Job = (*EnergyJob)(nil)
	_ Job = (*AirQualityJob)(nil)
	_ Job = (*WaterJob)(nil)

	errNoSensorsRetrieved = errors.New("no sensors retrieved")
)

// Job represents a WordPress automation task that executes Do to send a POST request
type Job interface {
	GetName() string
	GetUrl() string
	GetSite() string
	GetTicker() *time.Ticker

	Do(ctx context.Context, sendFn sender) error
	Stop()
}

type sender func(ctx context.Context, url string, body []byte) error

// Mulpx for multiple Job easy chan fan-in
type Mulpx struct {
	C     chan Job
	group errgroup.Group
}

func (m *Mulpx) WaitForDone() {
	// no errors are returned from any of the group.Go() calls
	_ = m.group.Wait()
	close(m.C)
}

// Multiplex receives multiple jobs and fans all the jobs into a single chan
func Multiplex(jobs ...Job) *Mulpx {
	out := &Mulpx{
		C: make(chan Job),
		// no need to listen to parent thread's context
		// as there should already be a listener for context.Done in the parent
		group: errgroup.Group{},
	}

	for _, job := range jobs {
		j := job

		out.group.Go(func() error {
			for range j.GetTicker().C {
				out.C <- j
			}

			return nil
		})
	}

	return out
}

// BaseJob shared fields
type BaseJob struct {
	Url    string
	Ticker *time.Ticker
	Site   string
	Logger *zap.Logger
}

func (b *BaseJob) GetUrl() string {
	return b.Url
}

func (b *BaseJob) GetTicker() *time.Ticker {
	return b.Ticker
}

func (b *BaseJob) GetSite() string {
	return b.Site
}

func (b *BaseJob) Stop() {
	b.Ticker.Stop()
}

func FromConfig(cfg config.Root, logger *zap.Logger, node *node.Node) []Job {
	var jobs []Job

	if cfg.Sources.Occupancy != nil && len(cfg.Sources.Occupancy.Sensors) > 0 {
		occ := &OccupancyJob{
			BaseJob: BaseJob{
				Site:   cfg.Site,
				Url:    fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Occupancy.Path),
				Ticker: time.NewTicker(cfg.Sources.Occupancy.Interval.Or(time.Minute)),
				Logger: logger,
			},
			Sensors: cfg.Sources.Occupancy.Sensors,
			client:  traits.NewOccupancySensorApiClient(node.ClientConn()),
		}

		jobs = append(jobs, occ)
	}
	if cfg.Sources.Temperature != nil && len(cfg.Sources.Temperature.Sensors) > 0 {
		temperature := &TemperatureJob{
			BaseJob: BaseJob{
				Site:   cfg.Site,
				Url:    fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Temperature.Path),
				Ticker: time.NewTicker(cfg.Sources.Temperature.Interval.Or(time.Minute)),
				Logger: logger,
			},
			Sensors: cfg.Sources.Temperature.Sensors,
			client:  traits.NewAirTemperatureApiClient(node.ClientConn()),
		}

		jobs = append(jobs, temperature)
	}
	if cfg.Sources.Energy != nil && len(cfg.Sources.Energy.Meters) > 0 {
		interval := cfg.Sources.Energy.Interval.Or(24 * time.Hour)
		energy := &EnergyJob{
			BaseJob: BaseJob{
				Site:   cfg.Site,
				Url:    fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Energy.Path),
				Ticker: time.NewTicker(interval),
				Logger: logger,
			},
			Meters:     cfg.Sources.Energy.Meters,
			Interval:   interval,
			client:     gen.NewMeterHistoryClient(node.ClientConn()),
			infoClient: gen.NewMeterInfoClient(node.ClientConn()),
		}

		jobs = append(jobs, energy)
	}
	if cfg.Sources.AirQuality != nil && len(cfg.Sources.AirQuality.Sensors) > 0 {
		air := &AirQualityJob{
			BaseJob: BaseJob{
				Site:   cfg.Site,
				Url:    fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.AirQuality.Path),
				Ticker: time.NewTicker(cfg.Sources.AirQuality.Interval.Or(time.Minute)),
				Logger: logger,
			},
			Sensors: cfg.Sources.AirQuality.Sensors,
			client:  traits.NewAirQualitySensorApiClient(node.ClientConn()),
		}

		jobs = append(jobs, air)
	}
	if cfg.Sources.Water != nil && len(cfg.Sources.Water.Meters) > 0 {
		interval := cfg.Sources.Water.Interval.Or(24 * time.Hour)
		water := &WaterJob{
			BaseJob: BaseJob{
				Site:   cfg.Site,
				Url:    fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Water.Path),
				Ticker: time.NewTicker(interval),
				Logger: logger,
			},
			Meters:     cfg.Sources.Energy.Meters,
			Interval:   interval,
			client:     gen.NewMeterHistoryClient(node.ClientConn()),
			infoClient: gen.NewMeterInfoClient(node.ClientConn()),
		}

		jobs = append(jobs, water)
	}

	return jobs
}
