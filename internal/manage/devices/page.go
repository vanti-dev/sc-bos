package devices

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type PageToken struct {
	ParentPageToken string `json:"t,omitempty"`
	PageIndex       int    `json:"i,omitempty"`
}

func (pt PageToken) encode() (string, error) {
	bytes, err := json.Marshal(pt)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(bytes), nil
}

func (pt PageToken) String() string {
	return fmt.Sprintf("%s@%d", pt.ParentPageToken, pt.PageIndex)
}

func decodePageToken(str string) (PageToken, error) {
	bytes, err := base64.RawStdEncoding.DecodeString(str)
	if err != nil {
		return PageToken{}, err
	}
	var pt PageToken
	return pt, json.Unmarshal(bytes, &pt)
}
