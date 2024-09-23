package job

import (
	"context"
	"time"
)

var (
	// strict:compile
	_ Job = (*OccupancyJob)(nil)
	_ Job = (*TemperatureJob)(nil)
	_ Job = (*EnergyJob)(nil)
	_ Job = (*AirQualityJob)(nil)
	_ Job = (*WaterJob)(nil)
)

// Job represents a wordpress automation task that executes Do to send a POST request
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
