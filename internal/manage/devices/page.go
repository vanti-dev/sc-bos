package devices

import (
	"encoding/base64"
	"encoding/json"
)

type PageToken struct {
	LastName string `json:"n"`
}

func (pt PageToken) encode() (string, error) {
	bytes, err := json.Marshal(pt)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(bytes), nil
}

func decodePageToken(str string) (PageToken, error) {
	bytes, err := base64.RawStdEncoding.DecodeString(str)
	if err != nil {
		return PageToken{}, err
	}
	var pt PageToken
	return pt, json.Unmarshal(bytes, &pt)
}
