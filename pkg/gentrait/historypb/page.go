package historypb

import (
	"context"
	"fmt"

	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"

	timepb "github.com/smart-core-os/sc-api/go/types/time"
	"github.com/vanti-dev/sc-bos/pkg/history"
)

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

func (pr pageReader[R]) listRecords(ctx context.Context, store history.Store, period *timepb.Period, pageSize int, pageToken string) (page []R, totalSize int, nextPageToken string, err error) {
	if pageSize == 0 {
		pageSize = pr.DefaultPageSize
	}
	if pageSize > pr.MaxPageSize {
		pageSize = pr.MaxPageSize
	}

	from, to := periodToRecords(period)
	slice := store.Slice(from, to)
	totalSize, err = slice.Len(ctx)
	if err != nil {
		return nil, 0, "", err
	}

	if pageToken != "" {
		from = history.Record{ID: pageToken}
		slice = slice.Slice(from, to)
	}

	dst := make([]history.Record, pageSize+1) // +1 to know if there's a next page or not
	read, err := slice.Read(ctx, dst)
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

	nextPageToken = dst[len(dst)-1].ID
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
