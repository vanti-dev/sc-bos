package services

import (
	"encoding/base64"
	"encoding/json"
)

type PageToken struct {
	NextId string `json:"n,omitempty"`
}

func (p PageToken) Encode() (string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(data), nil
}

func DecodePageToken(data string) (PageToken, error) {
	decoded, err := base64.RawStdEncoding.DecodeString(data)
	if err != nil {
		return PageToken{}, err
	}
	var p PageToken
	err = json.Unmarshal(decoded, &p)
	if err != nil {
		return PageToken{}, err
	}
	return p, nil
}
