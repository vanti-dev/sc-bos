package historypb

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"go.uber.org/multierr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	timepb "github.com/smart-core-os/sc-api/go/types/time"
	"github.com/vanti-dev/sc-bos/pkg/history"
)

//go:generate protomod protoc -- -I . -I ../../../proto --go_out=paths=source_relative:. historypb_page.proto

func newPageReader[R proto.Message](decodePayload func(r history.Record) (R, error)) pageReader[R] {
	return pageReader[R]{
		DefaultPageSize: 50,
		MaxPageSize:     1000,
		DecodePayload:   decodePayload,
	}
}

type pageReader[R proto.Message] struct {
	DefaultPageSize, MaxPageSize int

	DecodePayload func(r history.Record) (R, error)
}

func (pr pageReader[R]) listRecords(ctx context.Context, store history.Store, period *timepb.Period, pageSize int, pageToken, orderBy string) (page []R, totalSize int, nextPageToken string, err error) {
	if pageSize == 0 {
		pageSize = pr.DefaultPageSize
	}
	if pageSize > pr.MaxPageSize {
		pageSize = pr.MaxPageSize
	}

	tokenPb, err := unmarshalPageToken(pageToken)
	if err != nil {
		return nil, 0, "", status.Error(codes.InvalidArgument, "invalid page token")
	}

	from, to := periodToRecords(period)
	slice := store.Slice(from, to)
	totalSize = int(tokenPb.TotalSize)
	if totalSize == 0 {
		// avoid this potentially expensive step if possible
		totalSize, err = slice.Len(ctx)
		if err != nil {
			return nil, 0, "", err
		}
	}

	reader, pager, nexter, err := parseOrderBy(orderBy)
	if err != nil {
		return nil, 0, "", err
	}

	if tokenPb.RecordId != "" {
		slice = pager(slice, history.Record{ID: tokenPb.RecordId})
	}

	dst := make([]history.Record, pageSize+1) // +1 to know if there's a next page or not
	read, err := reader(slice, ctx, dst)
	if err != nil {
		return nil, 0, "", err
	}

	readRecords := read
	if readRecords > pageSize {
		readRecords = pageSize
	}
	page = make([]R, readRecords)

	var allErrs error
	for i := 0; i < readRecords; i++ {
		r, err := pr.DecodePayload(dst[i])
		if err != nil {
			allErrs = multierr.Append(allErrs, fmt.Errorf("%s %w", dst[i].ID, err))
			continue
		}
		page[i] = r
	}

	if read <= pageSize {
		return
	}

	nextPageToken, err = marshalPageToken(&PageToken{RecordId: nexter(dst).ID, TotalSize: int32(totalSize)})
	return
}

func periodToRecords(p *timepb.Period) (from, to history.Record) {
	if p == nil {
		return
	}
	if p.StartTime != nil {
		from.CreateTime = p.StartTime.AsTime()
	}
	if p.EndTime != nil {
		to.CreateTime = p.EndTime.AsTime()
	}
	return
}

func unmarshalPageToken(token string) (*PageToken, error) {
	if token == "" {
		return &PageToken{}, nil
	}
	data, err := base64.RawStdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	pb := &PageToken{}
	err = proto.Unmarshal(data, pb)
	return pb, err
}

func marshalPageToken(pb *PageToken) (string, error) {
	data, err := proto.Marshal(pb)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(data), nil
}

// These encapsulate the differences between ascending and descending order.
// sliceReadFunc calls Slice.Read or Slice.ReadDesc.
// slicePagerFunc applies page tokens to the slice, either slice[:i] or slice[i:].
// sliceNextFunc returns the record that will be used as the next page token, closely related to slicePagerFunc.
type (
	sliceReadFunc  = func(slice history.Slice, ctx context.Context, dst []history.Record) (int, error)
	slicePagerFunc = func(slice history.Slice, token history.Record) history.Slice
	sliceNextFunc  = func(all []history.Record) history.Record
)

func parseOrderBy(orderBy string) (sliceReadFunc, slicePagerFunc, sliceNextFunc, error) {
	// very simple parsing for now as we only support one field
	orderBy = strings.ToLower(orderBy)
	orderBy = strings.Join(strings.Fields(orderBy), " ") // normalize spaces
	switch orderBy {
	case "", "recordtime", "recordtime asc", "record_time", "record_time asc":
		return history.Slice.Read, ascPager, ascNext, nil
	case "recordtime desc", "record_time desc":
		return history.Slice.ReadDesc, descPager, descNext, nil
	default:
		return nil, nil, nil, status.Error(codes.InvalidArgument, "invalid order by")
	}
}

func ascPager(slice history.Slice, token history.Record) history.Slice {
	return slice.Slice(token, history.Record{})
}

func descPager(slice history.Slice, token history.Record) history.Slice {
	return slice.Slice(history.Record{}, token)
}

func ascNext(all []history.Record) history.Record {
	if len(all) < 1 {
		return history.Record{}
	}
	// The first record of the next page should be the last record read, we read +1 for this purpose
	return all[len(all)-1]
}

func descNext(all []history.Record) history.Record {
	if len(all) < 2 {
		return history.Record{}
	}
	// The first record of the next page should be the last record we return as it will be used as the endTime which
	// is exclusive. -2 because we've already read 1 past the end of the returned items to see if there will be more
	// pages.
	return all[len(all)-2]
}
