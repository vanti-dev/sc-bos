package history

import (
	"encoding/base64"
	"errors"

	"google.golang.org/protobuf/proto"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/history"
)

func normPageSize(pageSize int32) int32 {
	if pageSize < 0 {
		return 50
	}
	if pageSize > 1000 {
		return 1000
	}
	return pageSize
}

var errPageTokenEmpty = errors.New("page token empty")

func unmarshalPageToken(pageToken string) (history.Record, error) {
	if pageToken == "" {
		return history.Record{}, errPageTokenEmpty
	}
	bs, err := base64.RawStdEncoding.DecodeString(pageToken)
	if err != nil {
		return history.Record{}, err
	}
	var pbRecord *gen.HistoryRecord
	err = proto.Unmarshal(bs, pbRecord)
	if err != nil {
		return history.Record{}, err
	}
	_, hRecord := protoRecordToStoreRecord(pbRecord)
	return hRecord, nil
}

func marshalPageToken(record history.Record) (string, error) {
	record.Payload = nil
	pbRecord := storeRecordToProtoRecord("", record)
	bs, err := proto.Marshal(pbRecord)
	return base64.RawStdEncoding.EncodeToString(bs), err
}
