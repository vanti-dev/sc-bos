package http

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

type Client struct {
	client  *retryablehttp.Client
	logger  *zap.Logger
	headers http.Header
}

type Option func(cli *Client)

func New(opts ...Option) *Client {
	httpClient := retryablehttp.NewClient()
	// some retry config
	httpClient.RetryMax = 3
	httpClient.RetryWaitMax = time.Second * 10

	cli := &Client{
		client: httpClient,
		headers: http.Header{
			"Content-Type": []string{"application/json"},
			"User-Agent":   []string{"sc-bos"},
		},
	}

	for _, opt := range opts {
		opt(cli)
	}

	return cli
}

func WithAuthorizationBearer(token string) Option {
	return func(cli *Client) {
		cli.headers.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

func WithAuthorizationBasic(base64Token string) Option {
	return func(cli *Client) {
		cli.headers.Add("Authorization", fmt.Sprintf("Basic %s", base64Token))
	}
}

func WithLogger(noop bool, loggers ...*zap.Logger) Option {
	return func(cli *Client) {
		if noop || len(loggers) < 1 {
			cli.client.Logger = nil
			return
		}
		cli.logger = loggers[0]
		cli.client.Logger = &logWrapper{Logger: loggers[0]}
	}
}

func (c *Client) Post(ctx context.Context, url string, body []byte) error {
	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	req.Header = c.headers.Clone()

	var authHeaders []string
	for _, header := range c.headers.Values("Authorization") {
		authHeaders = append(authHeaders, header)
	}

	req.Header.Set("Authorization", strings.Join(authHeaders, " , "))

	resp, err := c.client.Do(req)

	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Error("closing response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code [%d] not 200 http.OK", resp.StatusCode)
	}

	return nil
}
