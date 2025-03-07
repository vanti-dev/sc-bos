package hpd

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type sensor interface {
	GetUpdate(response *SensorResponse) error
	GetName() string
}

type Poller struct {
	client       *Client
	pollInterval time.Duration
	logger       *zap.Logger
	sensors      []sensor
}

func NewPoller(client *Client, pollInterval time.Duration, logger *zap.Logger, sensors ...sensor) *Poller {
	if pollInterval <= 0 {
		pollInterval = time.Second * 60
	}

	return &Poller{
		client:       client,
		pollInterval: pollInterval,
		logger:       logger,
		sensors:      sensors,
	}
}

func (p *Poller) startPoll(ctx context.Context) {
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()
	p.process()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.process()
		}
	}
}

func (p *Poller) process() {
	response := SensorResponse{}
	if err := doGetRequest(p.client, &response, "sensor"); err != nil {
		p.logger.Error("failed to GET sensor", zap.Error(err))

		return
	}
	for _, sensor := range p.sensors {
		if err := sensor.GetUpdate(&response); err != nil {
			p.logger.Error("sensor failed refreshing data", zap.String("sensor", sensor.GetName()), zap.Error(err))
		}
	}
}
