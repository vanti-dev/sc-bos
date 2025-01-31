package history

import (
	"encoding/base64"
	"reflect"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/history"
)

func Test_unmarshalPageToken(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		wantRecord history.Record
		wantTotal  int
		wantErr    bool
	}{
		{"empty", "", history.Record{}, 0, true},
		{"invalid", "!@`", history.Record{}, 0, true},
		{"bad proto", "abc", history.Record{}, 0, true},
		{"token.id", makeToken(&gen.HistoryRecord{Id: "foo"}, 10), history.Record{ID: "foo"}, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRecord, gotTotal, err := unmarshalPageToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalPageToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRecord, tt.wantRecord) {
				t.Errorf("unmarshalPageToken() gotRecord = %v, wantRecord %v", gotRecord, tt.wantRecord)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("unmarshalPageToken() gotTotal = %v, wantTotal %v", gotTotal, tt.wantTotal)
			}
		})
	}
}

func makeToken(r *gen.HistoryRecord, total int) string {
	bs, err := proto.Marshal(&PageToken{Record: r, TotalSize: int32(total)})
	if err != nil {
		panic(err)
	}
	return base64.RawStdEncoding.EncodeToString(bs)
}
