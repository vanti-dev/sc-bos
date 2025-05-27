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
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
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

const defaultTimeout = time.Second * 5

// Job represents an exporthttp automation task that executes Do to send a POST request
type Job interface {
	GetName() string
	GetUrl() string
	GetSite() string
	GetExecutionAfter(t time.Time) <-chan time.Time
	SetPreviousExecution(t time.Time)

	Do(ctx context.Context, sendFn sender) error
}

type sender func(ctx context.Context, url string, body []byte) error

// Mulpx for multiple Job easy chan fan-in
type Mulpx struct {
	C     chan Job
	group *errgroup.Group
}

func (m *Mulpx) WaitForDone() {
	// no errors are returned from any of the group.Go() calls
	_ = m.group.Wait()
	close(m.C)
}

// Multiplex receives multiple jobs and fans all the jobs into a single chan
func Multiplex(ctx context.Context, jobs ...Job) *Mulpx {
	group, ctx := errgroup.WithContext(ctx)

	out := &Mulpx{
		C:     make(chan Job, len(jobs)),
		group: group,
	}

	for _, job := range jobs {
		j := job

		out.group.Go(func() error {
			out.C <- j
			current := time.Now().UTC()
			for {
				select {
				case currentExecTime := <-j.GetExecutionAfter(current):
					out.C <- j
					current = currentExecTime
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})
	}

	return out
}

// BaseJob shared fields
type BaseJob struct {
	Url               string
	Schedule          *jsontypes.Schedule
	Timeout           *jsontypes.Duration
	PreviousExecution time.Time
	Site              string
	Logger            *zap.Logger
}

func (b *BaseJob) GetUrl() string {
	return b.Url
}

func (b *BaseJob) GetExecutionAfter(t time.Time) <-chan time.Time {
	if t.IsZero() {
		t = time.Now().UTC()
	}
	return time.After(time.Until(b.Schedule.Next(t)))
}

func (b *BaseJob) SetPreviousExecution(t time.Time) {
	b.PreviousExecution = t
}

func (b *BaseJob) GetSite() string {
	return b.Site
}

func FromConfig(cfg config.Root, logger *zap.Logger, node *node.Node) []Job {
	var jobs []Job

	now := time.Now().UTC()

	if cfg.Sources.Occupancy != nil && len(cfg.Sources.Occupancy.Sensors) > 0 {
		occ := &OccupancyJob{
			BaseJob: BaseJob{
				Site:              cfg.Site,
				Url:               fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Occupancy.Path),
				Schedule:          cfg.Sources.Occupancy.Schedule,
				Logger:            logger,
				PreviousExecution: now,
				Timeout:           cfg.Sources.Occupancy.Timeout,
			},
			Sensors: cfg.Sources.Occupancy.Sensors,
			client:  traits.NewOccupancySensorApiClient(node.ClientConn()),
		}

		jobs = append(jobs, occ)
	}
	if cfg.Sources.Temperature != nil && len(cfg.Sources.Temperature.Sensors) > 0 {
		temperature := &TemperatureJob{
			BaseJob: BaseJob{
				Site:              cfg.Site,
				Url:               fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Temperature.Path),
				Schedule:          cfg.Sources.Temperature.Schedule,
				PreviousExecution: now,
				Logger:            logger,
				Timeout:           cfg.Sources.Temperature.Timeout,
			},
			Sensors: cfg.Sources.Temperature.Sensors,
			client:  traits.NewAirTemperatureApiClient(node.ClientConn()),
		}

		jobs = append(jobs, temperature)
	}
	if cfg.Sources.Energy != nil && len(cfg.Sources.Energy.Meters) > 0 {
		energy := &EnergyJob{
			BaseJob: BaseJob{
				Site:              cfg.Site,
				Url:               fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Energy.Path),
				Schedule:          cfg.Sources.Energy.Schedule,
				PreviousExecution: now,
				Logger:            logger,
				Timeout:           cfg.Sources.Energy.Timeout,
			},
			Meters:     cfg.Sources.Energy.Meters,
			client:     gen.NewMeterHistoryClient(node.ClientConn()),
			infoClient: gen.NewMeterInfoClient(node.ClientConn()),
		}

		jobs = append(jobs, energy)
	}
	if cfg.Sources.AirQuality != nil && len(cfg.Sources.AirQuality.Sensors) > 0 {
		air := &AirQualityJob{
			BaseJob: BaseJob{
				Site:              cfg.Site,
				Url:               fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.AirQuality.Path),
				Schedule:          cfg.Sources.AirQuality.Schedule,
				PreviousExecution: now,
				Logger:            logger,
				Timeout:           cfg.Sources.AirQuality.Timeout,
			},
			Sensors: cfg.Sources.AirQuality.Sensors,
			client:  traits.NewAirQualitySensorApiClient(node.ClientConn()),
		}

		jobs = append(jobs, air)
	}
	if cfg.Sources.Water != nil && len(cfg.Sources.Water.Meters) > 0 {
		water := &WaterJob{
			BaseJob: BaseJob{
				Site:              cfg.Site,
				Url:               fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Water.Path),
				Schedule:          cfg.Sources.Water.Schedule,
				PreviousExecution: now,
				Logger:            logger,
				Timeout:           cfg.Sources.Water.Timeout,
			},
			Meters:     cfg.Sources.Water.Meters,
			client:     gen.NewMeterHistoryClient(node.ClientConn()),
			infoClient: gen.NewMeterInfoClient(node.ClientConn()),
		}

		jobs = append(jobs, water)
	}

	return jobs
}
