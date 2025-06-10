package opcua

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"go.uber.org/zap"
)

type Client struct {
	client       *opcua.Client
	logger       *zap.Logger
	interval     time.Duration
	clientHandle uint32
}

func NewClient(client *opcua.Client, logger *zap.Logger, interval time.Duration, handle uint32) *Client {
	return &Client{
		client:       client,
		clientHandle: handle,
		interval:     interval,
		logger:       logger,
	}
}

func (c *Client) Subscribe(ctx context.Context, nodeId *ua.NodeID) (<-chan *opcua.PublishNotificationData, error) {
	notifyCh := make(chan *opcua.PublishNotificationData)
	sub, err := c.client.Subscribe(ctx, &opcua.SubscriptionParameters{
		Interval: c.interval,
	}, notifyCh)
	if err != nil {
		return nil, err
	}
	valueReq := opcua.NewMonitoredItemCreateRequestWithDefaults(nodeId, ua.AttributeIDValue, c.clientHandle)
	res, err := sub.Monitor(ctx, ua.TimestampsToReturnNeither, valueReq)
	if err != nil {
		return nil, err
	}
	if len(res.Results) > 1 || len(res.Results) == 0 {
		c.logger.Warn("expected one result", zap.Int("count", len(res.Results)), zap.Any("results", res.Results))
		return nil, fmt.Errorf("expected one result, got %d", len(res.Results))
	}
	if !errors.Is(res.Results[0].StatusCode, ua.StatusOK) {
		return nil, fmt.Errorf("error monitoring node: %s", res.Results[0].StatusCode.Error())
	}
	return notifyCh, nil
}
