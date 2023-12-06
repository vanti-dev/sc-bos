package hpd3

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"

	"go.uber.org/multierr"
)

type Client interface {
	FetchSensorData(ctx context.Context, point string) (out any, err error)
}

type HTTPClient struct {
	Client   *http.Client
	Host     string
	Password string
}

func (c *HTTPClient) FetchSensorData(ctx context.Context, point string) (any, error) {
	sensorPath := path.Join("rest", "sensor", url.PathEscape(point))
	output := make(map[string]any)
	err := c.getJSON(ctx, sensorPath, output)
	if err != nil {
		return nil, err
	}
	value, ok := output[point]
	if !ok {
		return nil, errors.New("response did not contain requested point")
	}
	return value, nil
}

func (c *HTTPClient) getJSON(ctx context.Context, path string, out any) (err error) {
	u := (&url.URL{
		Scheme: "https",
		Host:   c.Host,
		Path:   path,
	}).String()
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", basicAuth(c.Password))

	client := c.Client
	if client == nil {
		client = http.DefaultClient
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		err = multierr.Append(err, res.Body.Close())
	}()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}
	err = json.Unmarshal(body, out)
	return
}

func basicAuth(password string) string {
	// the Steinel sensors don't use usernames
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(":"+password))
}

func FetchPoints(ctx context.Context, client Client, points ...string) (map[string]any, error) {
	var errs error
	values := make(map[string]any)
	for _, point := range points {
		value, err := client.FetchSensorData(ctx, point)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
		values[point] = value
	}
	return values, errs
}
