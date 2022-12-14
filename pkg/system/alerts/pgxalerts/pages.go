package pgxalerts

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

const (
	DefaultPageSize = 50
	MaxPageSize     = 1000
)

type PageToken struct {
	LastCreateTime time.Time `json:"ct"`
	LastID         string    `json:"id"`
}

func (p PageToken) Encode() (text string, err error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(data), nil
}

func DecodePageToken(s string) (PageToken, error) {
	var pt PageToken
	data, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil {
		return pt, err
	}
	err = json.Unmarshal(data, &pt)
	return pt, err
}

func normalizePageSize(size int32) int {
	var pageSize int32 = DefaultPageSize
	if size > 0 {
		pageSize = size
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return int(pageSize)
}
