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
	client   *opcua.Client
	logger   *zap.Logger
	interval time.Duration
}

func NewClient(client *opcua.Client, logger *zap.Logger, interval time.Duration) *Client {
	return &Client{
		client:   client,
		interval: interval,
		logger:   logger,
	}
}

func (c *Client) Subscribe(ctx context.Context, node string) (<-chan *opcua.PublishNotificationData, error) {
	nodeId, err := ua.ParseNodeID(node)
	if err != nil {
		return nil, err
	}
	notifyCh := make(chan *opcua.PublishNotificationData)
	sub, err := c.client.Subscribe(ctx, &opcua.SubscriptionParameters{
		Interval: c.interval,
	}, notifyCh)
	if err != nil {
		return nil, err
	}
	valueReq := opcua.NewMonitoredItemCreateRequestWithDefaults(nodeId, ua.AttributeIDValue, 42)
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
