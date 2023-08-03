package history

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/history"
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
	slice := store.Slice(from, to)

	totalCount, _ := slice.Len(ctx) // ignore error, totalCount will be 0
	res := &gen.ListHistoryRecordsResponse{
		TotalSize: int32(totalCount),
	}

	pageSize := normPageSize(request.GetPageSize())
	pageToken, err := unmarshalPageToken(request.GetPageToken())
	switch {
	case errors.Is(err, errPageTokenEmpty):
	case err != nil:
		return nil, status.Errorf(codes.InvalidArgument, "page_token invalid %v", err)
	default:
		slice = slice.Slice(pageToken, to)
	}

	buf := make([]history.Record, pageSize+1) // +1 for nexPageToken calculation
	n, err := slice.Read(ctx, buf)
	if err != nil {
		return nil, err
	}
	if int32(n) > pageSize {
		// there's another page
		last := buf[n-1]
		buf = buf[:n-1]
		nextPageToken, err := marshalPageToken(last)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to marshal page token %v", err)
		}
		res.NextPageToken = nextPageToken
	} else {
		buf = buf[:n]
	}

	res.Records = make([]*gen.HistoryRecord, len(buf))
	for i, r := range buf {
		res.Records[i] = storeRecordToProtoRecord(source, r)
	}
	return res, nil
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
