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
	// FetchSensorData mutates out to insert the data from a single data point.
	// If the data cannot be fetched, then out is left unmodified and an error is returned.
	FetchSensorData(ctx context.Context, point string, out *PointData) error
}

type HTTPClient struct {
	Client   *http.Client
	Host     string
	Password string
}

func (c *HTTPClient) FetchSensorData(ctx context.Context, point string, out *PointData) error {
	sensorPath := path.Join("rest", "sensor", url.PathEscape(point))
	err := c.getJSON(ctx, sensorPath, out)
	return err
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

func FetchPoints(ctx context.Context, client Client, points ...string) (PointData, error) {
	var errs error
	var out PointData
	for _, point := range points {
		err := client.FetchSensorData(ctx, point, &out)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
	}
	return out, errs
}
