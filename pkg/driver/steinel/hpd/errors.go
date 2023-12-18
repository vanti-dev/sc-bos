package hpd

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

type Code int

const (
	CodeInvalidSettings          Code = 10010
	CodeInvalidContentType       Code = 11000
	CodeEmptyRequest             Code = 11001
	CodeContentLengthRequired    Code = 11002
	CodeInvalidJSON              Code = 11100
	CodeInvalidTimeFormat        Code = 11200
	CodeTimeInPast               Code = 11201
	CodeUpdateDowngrade          Code = 21000
	CodeUpdateUploading          Code = 21001
	CodeUpdateInstalling         Code = 21002
	CodeUpdateInvalid            Code = 21003
	CodeBackupIPSettingsChanged  Code = 21100
	CodeBackupNewerVersion       Code = 21101
	CodeDataPushIDModeError      Code = 31000
	CodeDataPushAgentTriggerBusy Code = 31001
)
