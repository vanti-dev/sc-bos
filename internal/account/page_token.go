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
