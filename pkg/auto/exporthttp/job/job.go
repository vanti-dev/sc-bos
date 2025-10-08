package job

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/timshannon/bolthold"
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

// boltKeyTemplate is the template string used to generate the bolt key
// "${AutoName}_${config.Root.Name}_${Job.Name}"
// It is used to store the job's last previous execution timestamp in the database.
const boltKeyTemplate = "%s_%s_%s"

const defaultTimeout = time.Second * 5

// Job represents an exporthttp automation task that executes Do to send a POST request
type Job interface {
	GetName() string

	GetNextExecution() <-chan time.Time
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
			// Since cron schedules often occur at a minimum of every minute,
			// we add a throttle to avoid a hot loop if jobs take a while to execute.
			throttle := time.NewTicker(time.Minute)
			defer throttle.Stop()

			for {
				select {
				case <-j.GetNextExecution():
					out.C <- j
				case <-ctx.Done():
					return ctx.Err()
				}

				<-throttle.C // throttle to avoid a hot loop if jobs take a while to execute
			}
		})
	}

	return out
}

func shouldExecuteImmediately(schedule *jsontypes.Schedule, now, previous time.Time) bool {
	if schedule == nil {
		return false // no schedule means it should never execute
	}
	if previous.IsZero() {
		return true // no previous execution means it should execute initially
	}
	if now.Equal(previous) {
		return false
	}
	interval := schedule.Next(previous).Sub(previous)
	return now.Sub(previous) >= interval // now is at least one interval after the previous execution
}

// BaseJob shared fields
type BaseJob struct {
	Url               string
	Schedule          *jsontypes.Schedule
	Timeout           *jsontypes.Duration
	Db                *bolthold.Store
	AutoName          string // the name of the automation e.g. "exporthttp"
	ScName            string // the name of automation instance the job belongs to, e.g. "site-name"
	Name              string // the name of the job, e.g. "occupancy"
	PreviousExecution time.Time
	Site              string
	Logger            *zap.Logger
}

func (b *BaseJob) GetName() string {
	return b.Name
}

func (b *BaseJob) GetNextExecution() <-chan time.Time {
	t := time.Now().UTC()

	previous := b.getPreviousExecution()
	executeImmediately := shouldExecuteImmediately(b.Schedule, t, previous.UTC())

	b.Logger.Debug("previous execution time detected", zap.String("name", b.Name), zap.Time("previous", previous), zap.Time("current", t), zap.Bool("executeImmediately", executeImmediately))

	if executeImmediately {
		return time.After(0)
	}

	return time.After(time.Until(b.Schedule.Next(t)))
}

func (b *BaseJob) SetPreviousExecution(t time.Time) {
	b.PreviousExecution = t
	key := fmt.Sprintf(boltKeyTemplate, b.AutoName, b.ScName, b.Name)
	if err := b.Db.Upsert(key, &t); err != nil {
		b.Logger.Warn("failed to update execution time", zap.Error(err), zap.String("key", key))
	}
}

func FromConfig(cfg config.Root, db *bolthold.Store, autoName, scName string, logger *zap.Logger, node *node.Node) []Job {
	var jobs []Job

	if cfg.Sources.Occupancy != nil && len(cfg.Sources.Occupancy.Sensors) > 0 {
		occ := &OccupancyJob{
			BaseJob: BaseJob{
				Site:     cfg.Site,
				Url:      fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Occupancy.Path),
				Schedule: cfg.Sources.Occupancy.Schedule,
				Db:       db,
				AutoName: autoName,
				Name:     "occupancy",
				ScName:   scName,
				Logger:   logger,
				Timeout:  cfg.Sources.Occupancy.Timeout,
			},
			Sensors: cfg.Sources.Occupancy.Sensors,
			client:  traits.NewOccupancySensorApiClient(node.ClientConn()),
		}

		occ.PreviousExecution = occ.getPreviousExecution()
		jobs = append(jobs, occ)
	}
	if cfg.Sources.Temperature != nil && len(cfg.Sources.Temperature.Sensors) > 0 {
		temperature := &TemperatureJob{
			BaseJob: BaseJob{
				Site:     cfg.Site,
				Url:      fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Temperature.Path),
				Schedule: cfg.Sources.Temperature.Schedule,
				Db:       db,
				AutoName: autoName,
				ScName:   scName,
				Name:     "temperature",
				Logger:   logger,
				Timeout:  cfg.Sources.Temperature.Timeout,
			},
			Sensors: cfg.Sources.Temperature.Sensors,
			client:  traits.NewAirTemperatureApiClient(node.ClientConn()),
		}

		temperature.PreviousExecution = temperature.getPreviousExecution()
		jobs = append(jobs, temperature)
	}
	if cfg.Sources.Energy != nil && len(cfg.Sources.Energy.Meters) > 0 {
		energy := &EnergyJob{
			BaseJob: BaseJob{
				Site:     cfg.Site,
				Url:      fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Energy.Path),
				Schedule: cfg.Sources.Energy.Schedule,
				Db:       db,
				AutoName: autoName,
				ScName:   scName,
				Name:     "energy",
				Logger:   logger,
				Timeout:  cfg.Sources.Energy.Timeout,
			},
			Meters:     cfg.Sources.Energy.Meters,
			client:     gen.NewMeterHistoryClient(node.ClientConn()),
			infoClient: gen.NewMeterInfoClient(node.ClientConn()),
		}

		energy.PreviousExecution = energy.getPreviousExecution()
		jobs = append(jobs, energy)
	}
	if cfg.Sources.AirQuality != nil && len(cfg.Sources.AirQuality.Sensors) > 0 {
		air := &AirQualityJob{
			BaseJob: BaseJob{
				Site:     cfg.Site,
				Url:      fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.AirQuality.Path),
				Schedule: cfg.Sources.AirQuality.Schedule,
				Db:       db,
				AutoName: autoName,
				ScName:   scName,
				Name:     "air_quality",
				Logger:   logger,
				Timeout:  cfg.Sources.AirQuality.Timeout,
			},
			Sensors: cfg.Sources.AirQuality.Sensors,
			client:  traits.NewAirQualitySensorApiClient(node.ClientConn()),
		}

		air.PreviousExecution = air.getPreviousExecution()
		jobs = append(jobs, air)
	}
	if cfg.Sources.Water != nil && len(cfg.Sources.Water.Meters) > 0 {
		water := &WaterJob{
			BaseJob: BaseJob{
				Site:     cfg.Site,
				Url:      fmt.Sprintf("%s/%s", cfg.BaseUrl, cfg.Sources.Water.Path),
				Schedule: cfg.Sources.Water.Schedule,
				Db:       db,
				AutoName: autoName,
				ScName:   scName,
				Name:     "water",
				Logger:   logger,
				Timeout:  cfg.Sources.Water.Timeout,
			},
			Meters:     cfg.Sources.Water.Meters,
			client:     gen.NewMeterHistoryClient(node.ClientConn()),
			infoClient: gen.NewMeterInfoClient(node.ClientConn()),
		}

		water.PreviousExecution = water.getPreviousExecution()
		jobs = append(jobs, water)
	}

	return jobs
}

func (b *BaseJob) getPreviousExecution() time.Time {
	previous := time.Time{}
	key := fmt.Sprintf(boltKeyTemplate, b.AutoName, b.ScName, b.Name)
	if err := b.Db.Get(key, &previous); err != nil {
		b.Logger.Error("failed to get previous execution time", zap.String("name", b.Name), zap.Error(err), zap.String("key", key))
		// assume the job executed successfully one interval ago if we can't retrieve the previous execution time
		now := time.Now()
		next := b.Schedule.Next(now)
		interval := b.Schedule.Next(next).Sub(next)
		previous = now.Add(-interval)
	}

	return previous
}
