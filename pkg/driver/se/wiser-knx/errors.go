package wiser_knx

import (
	"encoding/json"
	"fmt"
	"io"
)

type APIError struct {
	Code   int      `json:"code"`
	Info   string   `json:"info"`
	Detail []string `json:"detail,omitempty"`
}

func readError(src io.Reader) error {
	var res APIError
	rawJSON, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(rawJSON, &res)
	if err != nil {
		return err
	}
	return res
}

func (err APIError) Error() string {
	return fmt.Sprintf("%d - %s", err.Code, err.Info)
}
