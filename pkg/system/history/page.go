package history

import (
	"encoding/base64"
	"errors"

	"google.golang.org/protobuf/proto"

	"github.com/vanti-dev/sc-bos/pkg/history"
)

//go:generate protomod protoc -- -I . -I ../../../proto --go_out=paths=source_relative:. system_history_page.proto

func normPageSize(pageSize int32) int32 {
	if pageSize <= 0 {
		return 50
	}
	if pageSize > 1000 {
		return 1000
	}
	return pageSize
}

var errPageTokenEmpty = errors.New("page token empty")

func unmarshalPageToken(pageToken string) (history.Record, int, error) {
	fail := func(err error) (history.Record, int, error) {
		return history.Record{}, 0, err
	}
	if pageToken == "" {
		return fail(errPageTokenEmpty)
	}
	bs, err := base64.RawStdEncoding.DecodeString(pageToken)
	if err != nil {
		return fail(err)
	}
	pbRecord := &PageToken{}
	err = proto.Unmarshal(bs, pbRecord)
	if err != nil {
		return fail(err)
	}
	_, hRecord := protoRecordToStoreRecord(pbRecord.Record)
	return hRecord, int(pbRecord.TotalSize), nil
}

func marshalPageToken(record history.Record, totalSize int) (string, error) {
	record.Payload = nil
	pbRecord := storeRecordToProtoRecord("", record)
	bs, err := proto.Marshal(&PageToken{Record: pbRecord, TotalSize: int32(totalSize)})
	return base64.RawStdEncoding.EncodeToString(bs), err
}
