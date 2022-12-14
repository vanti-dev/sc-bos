package jsonapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/multierr"
)

func (c *Client) url(suffix string, args ...any) string {
	return c.baseURL + fmt.Sprintf(suffix, args...)
}

func (c *Client) postNoAuth(ctx context.Context, url string, reqJSON any, resJSON any) error {
	var req *http.Request
	if reqJSON == nil {
		var err error
		req, err = http.NewRequestWithContext(ctx, "POST", c.url(url), nil)
		if err != nil {
			return err
		}
	} else {
		bodyBytes, err := json.Marshal(reqJSON)
		if err != nil {
			return err
		}
		req, err = http.NewRequestWithContext(ctx, "POST", c.url(url), bytes.NewReader(bodyBytes))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		_ = consumeAndClose(res.Body)
		return fmt.Errorf("http %s %d %s", req.URL, res.StatusCode, res.Status)
	}

	if resJSON == nil {
		return consumeAndClose(res.Body)
	}

	return json.NewDecoder(res.Body).Decode(resJSON)
}

func (c *Client) post(ctx context.Context, url string, reqJSON any, resJSON any) error {
	return c.withToken(ctx, func(token string) error {
		var req *http.Request
		if reqJSON == nil {
			var err error
			req, err = http.NewRequestWithContext(ctx, "POST", c.url(url), nil)
			if err != nil {
				return err
			}
		} else {
			bodyBytes, err := json.Marshal(reqJSON)
			if err != nil {
				return err
			}
			req, err = http.NewRequestWithContext(ctx, "POST", c.url(url), bytes.NewReader(bodyBytes))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
		}
		req.Header.Set("Authorization", "Bearer "+token)

		res, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode == 401 {
			return ErrUnauthorizedResponse
		}
		if res.StatusCode != 200 {
			_ = consumeAndClose(res.Body)
			return fmt.Errorf("http %s %d %s", req.URL, res.StatusCode, res.Status)
		}

		if resJSON == nil {
			return consumeAndClose(res.Body)
		}

		return json.NewDecoder(res.Body).Decode(resJSON)
	})
}

func (c *Client) get(ctx context.Context, url string, resJSON any) error {
	return c.withToken(ctx, func(token string) error {
		req, err := http.NewRequestWithContext(ctx, "GET", c.url(url), nil)
		if err != nil {
			return err
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode == 401 {
			return ErrUnauthorizedResponse
		}
		if res.StatusCode != 200 {
			_ = consumeAndClose(res.Body)
			return fmt.Errorf("http %s %d %s", req.URL, res.StatusCode, res.Status)
		}
		if resJSON == nil {
			return consumeAndClose(res.Body)
		}

		return json.NewDecoder(res.Body).Decode(resJSON)
	})
}

func consumeAndClose(r io.ReadCloser) error {
	_, rErr := io.ReadAll(r)
	return multierr.Combine(rErr, r.Close())
}
