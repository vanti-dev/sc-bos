package pgxalerts

import (
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-bos/pkg/gen"
)

func Test_fieldMaskIncludesPath(t *testing.T) {
	tests := []struct {
		name string
		m    string
		p    string
		want bool
	}{
		{"nil", "", "", true},
		{"nil prop", "", "prop", true},
		{"has", "prop", "prop", true},
		{"includes", "bar,prop,foo", "prop", true},
		{"not includes", "bar,prop,foo", "baz", false},
		{"parent", "parent.child", "parent", true},
		{"parent.child", "parent.child", "parent.child", true},
		{"match parent", "parent", "parent.child", true},
		{"invert parent", "parent.child", "child", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m *fieldmaskpb.FieldMask
			if tt.m != "" {
				m = &fieldmaskpb.FieldMask{Paths: strings.Split(tt.m, ",")}
			}
			if got := fieldMaskIncludesPath(m, tt.p); got != tt.want {
				t.Errorf("fieldMaskIncludesPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertChangeForQuery(t *testing.T) {
	tests := []struct {
		name   string
		q      *gen.Alert_Query
		change *gen.PullAlertsResponse_Change
		want   *gen.PullAlertsResponse_Change
	}{
		{
			"nil query",
			nil,
			&gen.PullAlertsResponse_Change{Type: types.ChangeType_ADD, NewValue: &gen.Alert{Id: "01", Description: "Add alert"}},
			&gen.PullAlertsResponse_Change{Type: types.ChangeType_ADD, NewValue: &gen.Alert{Id: "01", Description: "Add alert"}},
		},
		{
			"convert to add",
			&gen.Alert_Query{Floor: "1"},
			&gen.PullAlertsResponse_Change{Type: types.ChangeType_UPDATE,
				OldValue: &gen.Alert{Id: "01", Floor: "2"},
				NewValue: &gen.Alert{Id: "01", Floor: "1"},
			},
			&gen.PullAlertsResponse_Change{Type: types.ChangeType_ADD, NewValue: &gen.Alert{Id: "01", Floor: "1"}},
		},
		{
			"convert to remove",
			&gen.Alert_Query{Floor: "1"},
			&gen.PullAlertsResponse_Change{Type: types.ChangeType_UPDATE,
				OldValue: &gen.Alert{Id: "01", Floor: "1"},
				NewValue: &gen.Alert{Id: "01", Floor: "2"},
			},
			&gen.PullAlertsResponse_Change{Type: types.ChangeType_REMOVE, OldValue: &gen.Alert{Id: "01", Floor: "1"}},
		},
		{
			"update still applies",
			&gen.Alert_Query{Floor: "1"},
			&gen.PullAlertsResponse_Change{Type: types.ChangeType_UPDATE,
				OldValue: &gen.Alert{Id: "01", Floor: "1", Zone: "Z1"},
				NewValue: &gen.Alert{Id: "01", Floor: "1", Zone: "Z2"},
			},
			&gen.PullAlertsResponse_Change{Type: types.ChangeType_UPDATE,
				OldValue: &gen.Alert{Id: "01", Floor: "1", Zone: "Z1"},
				NewValue: &gen.Alert{Id: "01", Floor: "1", Zone: "Z2"},
			},
		},
		{
			"add doesn't match",
			&gen.Alert_Query{Floor: "1"},
			&gen.PullAlertsResponse_Change{Type: types.ChangeType_ADD,
				NewValue: &gen.Alert{Id: "01", Floor: "2"},
			},
			nil,
		},
		{
			"delete doesn't match",
			&gen.Alert_Query{Floor: "1"},
			&gen.PullAlertsResponse_Change{Type: types.ChangeType_REMOVE,
				OldValue: &gen.Alert{Id: "01", Floor: "2"},
			},
			nil,
		},
		{
			"update doesn't match",
			&gen.Alert_Query{Floor: "1"},
			&gen.PullAlertsResponse_Change{Type: types.ChangeType_UPDATE,
				OldValue: &gen.Alert{Id: "01", Floor: "2"},
				NewValue: &gen.Alert{Id: "01", Floor: "3"},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertChangeForQuery(tt.q, tt.change)
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("convertChangeForQuery() (-want,+got)\n%s", diff)
			}
		})
	}
}

func Test_alertMatchesQuery(t *testing.T) {
	ack := true
	noAck := false
	tests := []struct {
		name string
		q    *gen.Alert_Query
		a    *gen.Alert
		want bool
	}{
		{"nil query", nil, &gen.Alert{Id: "any"}, true},
		{"empty query", &gen.Alert_Query{}, &gen.Alert{Id: "any"}, true},
		{"floor yes", &gen.Alert_Query{Floor: "1"}, &gen.Alert{Floor: "1"}, true},
		{"floor no", &gen.Alert_Query{Floor: "1"}, &gen.Alert{Floor: "2"}, false},
		{"floor absent", &gen.Alert_Query{Floor: "1"}, &gen.Alert{}, false},
		{"zone yes", &gen.Alert_Query{Zone: "1"}, &gen.Alert{Zone: "1"}, true},
		{"zone no", &gen.Alert_Query{Zone: "1"}, &gen.Alert{Zone: "2"}, false},
		{"zone absent", &gen.Alert_Query{Zone: "1"}, &gen.Alert{}, false},
		{"source yes", &gen.Alert_Query{Source: "1"}, &gen.Alert{Source: "1"}, true},
		{"source no", &gen.Alert_Query{Source: "1"}, &gen.Alert{Source: "2"}, false},
		{"source absent", &gen.Alert_Query{Source: "1"}, &gen.Alert{}, false},
		{"acknowledged yes", &gen.Alert_Query{Acknowledged: &ack}, &gen.Alert{Acknowledgement: &gen.Alert_Acknowledgement{}}, true},
		{"acknowledged no", &gen.Alert_Query{Acknowledged: &ack}, &gen.Alert{}, false},
		{"not acknowledged yes", &gen.Alert_Query{Acknowledged: &noAck}, &gen.Alert{}, true},
		{"not acknowledged no", &gen.Alert_Query{Acknowledged: &noAck}, &gen.Alert{Acknowledgement: &gen.Alert_Acknowledgement{}}, false},
		{"create_time before", &gen.Alert_Query{CreatedNotBefore: timestamppb.New(time.Unix(100, 1))}, &gen.Alert{CreateTime: timestamppb.New(time.Unix(100, 0))}, false},
		{"create_time start", &gen.Alert_Query{CreatedNotBefore: timestamppb.New(time.Unix(100, 0))}, &gen.Alert{CreateTime: timestamppb.New(time.Unix(100, 0))}, true},
		{"create_time after", &gen.Alert_Query{CreatedNotBefore: timestamppb.New(time.Unix(100, 0))}, &gen.Alert{CreateTime: timestamppb.New(time.Unix(100, 1))}, true},
		{"create_time early", &gen.Alert_Query{CreatedNotAfter: timestamppb.New(time.Unix(100, 1))}, &gen.Alert{CreateTime: timestamppb.New(time.Unix(100, 0))}, true},
		{"create_time end", &gen.Alert_Query{CreatedNotAfter: timestamppb.New(time.Unix(100, 0))}, &gen.Alert{CreateTime: timestamppb.New(time.Unix(100, 0))}, true},
		{"create_time late", &gen.Alert_Query{CreatedNotAfter: timestamppb.New(time.Unix(100, 0))}, &gen.Alert{CreateTime: timestamppb.New(time.Unix(100, 1))}, false},
		{"create_time within", &gen.Alert_Query{
			CreatedNotBefore: timestamppb.New(time.Unix(100, 0)),
			CreatedNotAfter:  timestamppb.New(time.Unix(200, 0)),
		}, &gen.Alert{CreateTime: timestamppb.New(time.Unix(150, 0))}, true},
		{"severity low", &gen.Alert_Query{SeverityNotBelow: 2, SeverityNotAbove: 5}, &gen.Alert{Severity: gen.Alert_Severity(1)}, false},
		{"severity bottom", &gen.Alert_Query{SeverityNotBelow: 2, SeverityNotAbove: 5}, &gen.Alert{Severity: gen.Alert_Severity(2)}, true},
		{"severity within", &gen.Alert_Query{SeverityNotBelow: 2, SeverityNotAbove: 5}, &gen.Alert{Severity: gen.Alert_Severity(4)}, true},
		{"severity top", &gen.Alert_Query{SeverityNotBelow: 2, SeverityNotAbove: 5}, &gen.Alert{Severity: gen.Alert_Severity(5)}, true},
		{"severity high", &gen.Alert_Query{SeverityNotBelow: 2, SeverityNotAbove: 5}, &gen.Alert{Severity: gen.Alert_Severity(6)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := alertMatchesQuery(tt.q, tt.a); got != tt.want {
				t.Errorf("alertMatchesQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
