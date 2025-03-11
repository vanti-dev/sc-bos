package account

import (
	"encoding/base64"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func encodePageToken(token *PageToken) string {
	tokenBytes, err := proto.Marshal(token)
	if err != nil {
		// marshalling this type cannot error
		panic(fmt.Errorf("failed to marshal page token: %w", err))
	}
	return base64.URLEncoding.EncodeToString(tokenBytes)
}

func parsePageToken(token string, filter string) (*PageToken, error) {
	tokenBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	parsed := &PageToken{}
	err = proto.Unmarshal(tokenBytes, parsed)
	if err != nil {
		return nil, err
	}

	if parsed.Filter != filter {
		return nil, ErrInvalidPageToken
	}
	return parsed, nil
}

// like parsePageToken, but additionally validates that an integer last ID is present
func parseIntPageToken(token string, filter string) (*PageToken, error) {
	res, err := parsePageToken(token, filter)
	if err != nil {
		return nil, err
	}
	_, ok := res.GetLastId().(*PageToken_LastIdIntPk)
	if !ok {
		return nil, ErrInvalidPageToken
	}
	return res, nil
}

// like parsePageToken, but additionally validates that a string last ID is present
func parseStringPageToken(token string, filter string) (*PageToken, error) {
	res, err := parsePageToken(token, filter)
	if err != nil {
		return nil, err
	}
	_, ok := res.GetLastId().(*PageToken_LastNaturalId)
	if !ok {
		return nil, ErrInvalidPageToken
	}
	return res, nil
}
