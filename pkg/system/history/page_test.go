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
		name    string
		token   string
		want    history.Record
		wantErr bool
	}{
		{"empty", "", history.Record{}, true},
		{"invalid", "!@`", history.Record{}, true},
		{"bad proto", "abc", history.Record{}, true},
		{"token.id", makeToken(&gen.HistoryRecord{Id: "foo"}), history.Record{ID: "foo"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unmarshalPageToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalPageToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unmarshalPageToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func makeToken(r *gen.HistoryRecord) string {
	bs, err := proto.Marshal(r)
	if err != nil {
		panic(err)
	}
	return base64.RawStdEncoding.EncodeToString(bs)
}
