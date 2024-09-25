package job

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/vanti-dev/sc-bos/pkg/auto/wordpress/config"
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
	GetClients() []any

	Do(ctx context.Context, sendFn sender) error
}

type sender func(ctx context.Context, url string, body []byte) error

// Mulpx for multiple Job easy chan fan-in
type Mulpx struct {
	C    chan Job
	Done chan struct{}
}

// Multiplex receives multiple jobs and fans all the jobs into a single chan
func Multiplex(jobs ...Job) *Mulpx {
	out := &Mulpx{
		C:    make(chan Job),
		Done: make(chan struct{}),
	}

	for _, job := range jobs {
		j := job

		go func() {
			for range j.GetTicker().C {
				out.C <- j
			}
		}()
	}

	// clean up
	go func() {
		defer close(out.C)
		<-out.Done
	}()

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

func FromConfig(cfg config.Root, logger *zap.Logger) []Job {
	var jobs []Job

	if cfg.Sources.Occupancy != nil && len(cfg.Sources.Occupancy.Sensors) > 0 {
		occ := &OccupancyJob{
			BaseJob: BaseJob{
				Site:   cfg.Site,
				Url:    fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Occupancy.Path),
				Ticker: time.NewTicker(cfg.Sources.Occupancy.Interval.Duration),
				Logger: logger,
			},
			Sensors: cfg.Sources.Occupancy.Sensors,
		}

		jobs = append(jobs, occ)
	}
	if cfg.Sources.Temperature != nil && len(cfg.Sources.Temperature.Sensors) > 0 {
		temperature := &TemperatureJob{
			BaseJob: BaseJob{
				Site:   cfg.Site,
				Url:    fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Temperature.Path),
				Ticker: time.NewTicker(cfg.Sources.Temperature.Interval.Duration),
				Logger: logger,
			},
			Sensors: cfg.Sources.Temperature.Sensors,
		}

		jobs = append(jobs, temperature)
	}
	if cfg.Sources.Energy != nil && len(cfg.Sources.Energy.Meters) > 0 {
		energy := &EnergyJob{
			BaseJob: BaseJob{
				Site:   cfg.Site,
				Url:    fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Energy.Path),
				Ticker: time.NewTicker(cfg.Sources.Energy.Interval.Duration),
				Logger: logger,
			},
			Meters:   cfg.Sources.Energy.Meters,
			Interval: cfg.Sources.Energy.Interval.Duration,
		}

		jobs = append(jobs, energy)
	}
	if cfg.Sources.AirQuality != nil && len(cfg.Sources.AirQuality.Sensors) > 0 {
		air := &AirQualityJob{
			BaseJob: BaseJob{
				Site:   cfg.Site,
				Url:    fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.AirQuality.Path),
				Ticker: time.NewTicker(cfg.Sources.AirQuality.Interval.Duration),
				Logger: logger,
			},
			Sensors: cfg.Sources.AirQuality.Sensors,
		}

		jobs = append(jobs, air)
	}
	if cfg.Sources.Water != nil && len(cfg.Sources.Water.Meters) > 0 {
		water := &WaterJob{
			BaseJob: BaseJob{
				Site:   cfg.Site,
				Url:    fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Water.Path),
				Ticker: time.NewTicker(cfg.Sources.Water.Interval.Duration),
				Logger: logger,
			},
			Meters:   cfg.Sources.Energy.Meters,
			Interval: cfg.Sources.Energy.Interval.Duration,
		}

		jobs = append(jobs, water)
	}

	return jobs
}
