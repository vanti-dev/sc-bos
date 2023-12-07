package dps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/multierr"
)

func httpGetJSON(ctx context.Context, url, auth string, resBody any) (http.Header, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("construct request: %w", err)
	}
	req.Header.Set("authorization", auth)
	resBodyBytes, header, err := httpDo(req)
	if err != nil {
		return header, err
	}
	return header, json.Unmarshal(resBodyBytes, resBody)
}

func httpPutJSON(ctx context.Context, url, auth string, reqBody, resBody any) (http.Header, error) {
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("construct request: %w", err)
	}
	req.Header.Set("authorization", auth)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("content-encoding", "utf-8")
	resBodyBytes, header, err := httpDo(req)
	if err != nil {
		return header, err
	}

	err = json.Unmarshal(resBodyBytes, resBody)
	return header, err
}

func httpDo(req *http.Request) (resBody []byte, resHeader http.Header, err error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	// limit to 10 MiB because this function is only intended to fetch small JSON payloads.
	// anything larger should use some custom code to stream the data properly
	resBodyBytes, err := io.ReadAll(io.LimitReader(res.Body, 10*1024*1024))
	err = multierr.Append(err, res.Body.Close())
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		err = multierr.Append(err, fmt.Errorf("HTTP response: %d %s", res.StatusCode, res.Status))
	}
	return resBodyBytes, res.Header, err
}
