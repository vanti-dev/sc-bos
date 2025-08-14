package parentdevices

import (
	"encoding/base64"
	"slices"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/types"
)

type listRequest interface {
	GetPageToken() string
	GetPageSize() int32
}

type namer interface {
	GetName() string
}

func listPage[T namer](req listRequest, list func() []T) (_ []T, totalSize int, nextPageToken string, _ error) {
	pageToken := &types.PageToken{}
	if err := decodePageToken(req.GetPageToken(), pageToken); err != nil {
		return nil, 0, "", err
	}

	lastName := pageToken.GetLastResourceName() // The name of the last item of the previous page
	pageSize := capPageSize(req.GetPageSize())

	all := list()
	slices.SortFunc(all, func(a, b T) int {
		return strings.Compare(a.GetName(), b.GetName())
	})

	nextIndex := 0
	if lastName != "" {
		var found bool
		nextIndex, found = slices.BinarySearchFunc(all, lastName, func(t T, s string) int {
			return strings.Compare(s, t.GetName())
		})
		if found {
			nextIndex++
		}
	}

	endOfPage := nextIndex + pageSize
	if endOfPage > len(all) {
		// no more pages
		endOfPage = len(all)
		pageToken = nil
	} else {
		pageToken.PageStart = &types.PageToken_LastResourceName{
			LastResourceName: all[endOfPage-1].GetName(),
		}
	}

	nextPageToken, err := encodePageToken(pageToken)
	if err != nil {
		return nil, 0, "", err
	}

	return all[nextIndex:endOfPage], len(all), nextPageToken, nil
}

const (
	defaultPageSize = 50
	maxPageSize     = 1000
)

func capPageSize(pageSize int32) int {
	if pageSize == 0 {
		return defaultPageSize
	}
	if pageSize > maxPageSize {
		return maxPageSize
	}
	return int(pageSize)
}

func decodePageToken(token string, pageToken *types.PageToken) error {
	if token != "" {
		tokenBytes, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "bad page token: %v", err)
		}
		if err := proto.Unmarshal(tokenBytes, pageToken); err != nil {
			return status.Errorf(codes.InvalidArgument, "bad page token: %v", err)
		}
	}
	return nil
}

func encodePageToken(pageToken *types.PageToken) (string, error) {
	if pageToken != nil {
		tokenBytes, err := proto.Marshal(pageToken)
		if err != nil {
			return "", status.Errorf(codes.Unknown, "unable to create page token: %v", err)
		}
		return base64.StdEncoding.EncodeToString(tokenBytes), nil
	}
	return "", nil
}
