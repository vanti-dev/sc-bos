package apistore

import (
	"context"
	"sync"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/history"
)

// New creates a new history.Store backed by the given client and name.
// All records created or read by this store will have the given source.
func New(client gen.HistoryAdminApiClient, name, source string) *Store {
	return &Store{
		slice: slice{
			client: client,
			name:   name,
			source: source,
		},
	}
}

// Store implements history.Store backed by a gen.HistoryAdminApiClient.
// As a quirk of how the api works, it's more efficient to call Read then Len, than Len then Read.
type Store struct {
	slice
}

var _ history.Store = (*Store)(nil)

func (s *Store) Append(ctx context.Context, payload []byte) (history.Record, error) {
	pbRecord, err := s.client.CreateHistoryRecord(ctx, &gen.CreateHistoryRecordRequest{
		Name: s.name,
		Record: &gen.HistoryRecord{
			Source:  s.source,
			Payload: payload,
		},
	})
	if err != nil {
		return history.Record{}, err
	}
	_, hRecord := protoRecordToStoreRecord(pbRecord)
	return hRecord, nil
}

type slice struct {
	client gen.HistoryAdminApiClient
	name   string

	source   string
	from, to history.Record

	totalSizeMu sync.Mutex
	totalSize   int
	totalSizeOk bool // true if we've successfully read totalSize from client
}

func (s *slice) Slice(from, to history.Record) history.Slice {
	return &slice{
		client: s.client,
		name:   s.name,
		source: s.source,
		from:   from,
		to:     to,
	}
}

func (s *slice) Read(ctx context.Context, into []history.Record) (int, error) {
	req := s.newListRequest(int32(len(into)))

	i := 0
	for {
		res, err := s.client.ListHistoryRecords(ctx, req)
		if err != nil {
			return 0, err
		}

		s.totalSizeMu.Lock()
		s.totalSize = int(res.TotalSize)
		s.totalSizeOk = true
		s.totalSizeMu.Unlock()

		for _, record := range res.Records {
			if i >= len(into) {
				// should only happen if the server doesn't respect our page size
				return i, nil
			}
			_, into[i] = protoRecordToStoreRecord(record)
			i++
		}

		req.PageToken = res.NextPageToken
		if req.PageToken == "" {
			return i, nil
		}
	}
}

func (s *slice) newListRequest(pageSize int32) *gen.ListHistoryRecordsRequest {
	req := &gen.ListHistoryRecordsRequest{
		Name:     s.name,
		PageSize: pageSize,
		Query: &gen.HistoryRecord_Query{
			Source: &gen.HistoryRecord_Query_SourceEqual{SourceEqual: s.source},
		},
	}
	if !s.from.IsZero() {
		req.Query.FromRecord = storeRecordToProtoRecord("", s.from)
	}
	if !s.to.IsZero() {
		req.Query.ToRecord = storeRecordToProtoRecord("", s.to)
	}
	return req
}

func (s *slice) Len(ctx context.Context) (int, error) {
	s.totalSizeMu.Lock()
	defer s.totalSizeMu.Unlock()
	if !s.totalSizeOk {
		req := s.newListRequest(1)
		res, err := s.client.ListHistoryRecords(ctx, req)
		if err != nil {
			return 0, err
		}
		s.totalSizeOk = true
		s.totalSize = int(res.TotalSize)
	}
	return s.totalSize, nil
}

func protoRecordToStoreRecord(r *gen.HistoryRecord) (string, history.Record) {
	return r.GetSource(), history.Record{
		ID:         r.GetId(),
		CreateTime: r.GetCreateTime().AsTime(),
		Payload:    r.GetPayload(),
	}
}

func storeRecordToProtoRecord(source string, r history.Record) *gen.HistoryRecord {
	return &gen.HistoryRecord{
		Id:         r.ID,
		Source:     source,
		CreateTime: timestamppb.New(r.CreateTime),
		Payload:    r.Payload,
	}
}
