package meter

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/vanti-dev/sc-bos/pkg/gen"
)

func Test_mergeMeterReading(t *testing.T) {
	err := errors.New("expected error")
	reading := func(val float32) *gen.MeterReading {
		return &gen.MeterReading{
			Usage: val,
		}
	}

	tests := []struct {
		in      []value
		want    *gen.MeterReading
		wantErr bool
	}{
		// simple cases
		{nil, nil, true},
		{[]value{{}}, nil, true},
		{[]value{{err: err}}, nil, true},
		{[]value{{err: err, val: reading(10)}}, nil, true},
		// all present
		{[]value{{val: reading(10)}}, reading(10), false},
		{[]value{{val: reading(10)}, {val: reading(20)}}, reading(30), false},
		// some missing
		{[]value{{val: reading(10)}, {}}, nil, true},
		{[]value{{}, {val: reading(10)}}, nil, true},
		// some errors
		{[]value{{err: err}, {}}, nil, true},
		{[]value{{}, {err: err}}, nil, true},
		// mixed missing and error
		{[]value{{}, {err: err}, {val: reading(10)}}, nil, true},
	}
	for _, tt := range tests {
		name := ""
		if len(tt.in) == 0 {
			name = "empty"
		} else {
			var names []string
			for _, v := range tt.in {
				switch {
				case v.err != nil && v.val != nil:
					names = append(names, "err+val")
				case v.err != nil:
					names = append(names, "err")
				case v.val == nil:
					names = append(names, "nil")
				default:
					names = append(names, fmt.Sprintf("%v", v.val.Usage))
				}
			}
			name = strings.Join(names, ",")
		}
		t.Run(name, func(t *testing.T) {
			got, err := mergeMeterReading(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("mergeMeterReading() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("mergeMeterReading() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
