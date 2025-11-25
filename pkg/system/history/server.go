package history

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/historypb"
	"github.com/smart-core-os/sc-bos/pkg/history"
)

type storeServer struct {
	gen.UnimplementedHistoryAdminApiServer
	store func(source string) history.Store
}

func (s *storeServer) CreateHistoryRecord(ctx context.Context, request *gen.CreateHistoryRecordRequest) (*gen.HistoryRecord, error) {
	if err := validateCreateRequest(request); err != nil {
		return nil, err
	}
	record := request.GetRecord()
	r, err := s.store(record.GetSource()).Append(ctx, record.Payload)
	if err != nil {
		return nil, err
	}
	return storeRecordToProtoRecord(record.Source, r), nil
}

func (s *storeServer) ListHistoryRecords(ctx context.Context, request *gen.ListHistoryRecordsRequest) (*gen.ListHistoryRecordsResponse, error) {
	source := request.GetQuery().GetSourceEqual()
	if source == "" {
		return nil, status.Error(codes.InvalidArgument, "source_equal must be set")
	}
	_, from := protoRecordToStoreRecord(request.GetQuery().GetFromRecord())
	_, to := protoRecordToStoreRecord(request.GetQuery().GetToRecord())

	store := s.store(source)
	pager := newPageReader(source)
	page, size, nextToken, err := pager.ListRecordsBetween(ctx, store, from, to, int(request.GetPageSize()), request.GetPageToken(), request.GetOrderBy())
	if err != nil {
		return nil, err
	}

	return &gen.ListHistoryRecordsResponse{
		Records:       page,
		NextPageToken: nextToken,
		TotalSize:     int32(size),
	}, nil
}

func newPageReader(source string) historypb.PageReader[*gen.HistoryRecord] {
	pr := historypb.NewPageReader(func(r history.Record) (*gen.HistoryRecord, error) {
		return storeRecordToProtoRecord(source, r), nil
	})
	// we use create_time in our API, so override the default record_time parsing
	pr.OrderByParser = func(s string) historypb.OrderBy {
		sn := strings.ToLower(s)
		sn = strings.Join(strings.Fields(sn), " ") // normalise whitespace
		switch sn {
		case "", "createtime", "createtime asc", "create_time", "create_time asc":
			return historypb.OrderByTimeAsc
		case "createtime desc", "create_time desc":
			return historypb.OrderByTimeDesc
		default:
			return historypb.OrderBy(s)
		}
	}
	return pr
}

func protoRecordToStoreRecord(r *gen.HistoryRecord) (string, history.Record) {
	hRecord := history.Record{
		ID:      r.GetId(),
		Payload: r.GetPayload(),
	}
	if r.GetCreateTime() != nil {
		hRecord.CreateTime = r.GetCreateTime().AsTime()
	}
	return r.GetSource(), hRecord
}

func storeRecordToProtoRecord(source string, r history.Record) *gen.HistoryRecord {
	pbRecord := &gen.HistoryRecord{
		Id:      r.ID,
		Source:  source,
		Payload: r.Payload,
	}
	if !r.CreateTime.IsZero() {
		pbRecord.CreateTime = timestamppb.New(r.CreateTime)
	}
	return pbRecord
}

func validateCreateRequest(request *gen.CreateHistoryRecordRequest) error {
	switch {
	case request.GetRecord().GetId() != "":
		return status.Error(codes.InvalidArgument, "id must not be set")
	case request.GetRecord().GetSource() == "":
		return status.Error(codes.InvalidArgument, "source must be set")
	case request.GetRecord().GetCreateTime() != nil:
		return status.Error(codes.InvalidArgument, "create_time must not be set")
	}
	return nil
}
