package fetch

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func JSON(ctx context.Context, url string, into any, options ...Option) error {
	o := resolveOpts(options...)

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	client := o.httpClient
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return HTTPError{response.StatusCode, response.Status}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, into)
}
