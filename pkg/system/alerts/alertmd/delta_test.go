package alertmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-golang/pkg/resource"
)

func Test_ApplyMdDelta(t *testing.T) {
	before := &gen.AlertMetadata{
		TotalCount: 100,
		AcknowledgedCounts: map[bool]uint32{
			true:  80,
			false: 20,
		},
		ResolvedCounts: map[bool]uint32{
			true:  20,
			false: 80,
		},
		NeedsAttentionCounts: map[string]uint32{
			"ack_resolved":    10,
			"ack_unresolved":  70,
			"nack_resolved":   5,
			"nack_unresolved": 15,
		},
		FloorCounts: map[string]uint32{
			"Floor1": 20,
			"Floor2": 30,
			"Floor3": 50,
			"Floor4": 0,
		},
		ZoneCounts: map[string]uint32{
			"Zone1": 20,
			"Zone2": 30,
			"Zone3": 50,
			"Zone4": 0,
		},
		SeverityCounts: map[int32]uint32{
			1: 20,
			2: 30,
			3: 50,
			4: 0,
		},
		SubsystemCounts: map[string]uint32{
			"Subsystem1": 20,
			"Subsystem2": 30,
		},
	}
	// base structs for an added and removed alert without acknowledgement
	added := patch(before, &gen.AlertMetadata{
		TotalCount:         101,
		AcknowledgedCounts: map[bool]uint32{false: 21},
		ResolvedCounts:     map[bool]uint32{false: 81},
		NeedsAttentionCounts: map[string]uint32{
			"nack_unresolved": 16,
		},
	})
	removed := patch(before, &gen.AlertMetadata{
		TotalCount:         99,
		AcknowledgedCounts: map[bool]uint32{false: 19},
		ResolvedCounts:     map[bool]uint32{false: 79},
		NeedsAttentionCounts: map[string]uint32{
			"nack_unresolved": 14,
		},
	})

	tests := []struct {
		name          string
		before, after *gen.AlertMetadata
		e             *gen.PullAlertsResponse_Change
		wantErr       bool
	}{
		{"no change", before, before, &gen.PullAlertsResponse_Change{}, false},
		{"no change (zero metadata)", &gen.AlertMetadata{}, &gen.AlertMetadata{}, &gen.PullAlertsResponse_Change{}, false},

		{"add empty", before, added, &gen.PullAlertsResponse_Change{NewValue: &gen.Alert{}}, false},
		{"remove empty", before, removed, &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{}}, false},
		{"add ack", before, patch(added, &gen.AlertMetadata{
			AcknowledgedCounts:   map[bool]uint32{true: 81, false: 20},
			NeedsAttentionCounts: map[string]uint32{"nack_unresolved": 15, "ack_unresolved": 71},
		}), &gen.PullAlertsResponse_Change{NewValue: &gen.Alert{Acknowledgement: &gen.Alert_Acknowledgement{AcknowledgeTime: timestamppb.Now()}}}, false},
		{"remove ack", before, patch(removed, &gen.AlertMetadata{
			AcknowledgedCounts:   map[bool]uint32{true: 79, false: 20},
			NeedsAttentionCounts: map[string]uint32{"nack_unresolved": 15, "ack_unresolved": 69},
		}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Acknowledgement: &gen.Alert_Acknowledgement{AcknowledgeTime: timestamppb.Now()}}}, false},
		{"add floor", before, patch(added, &gen.AlertMetadata{FloorCounts: map[string]uint32{"Floor1": 21}}), &gen.PullAlertsResponse_Change{NewValue: &gen.Alert{Floor: "Floor1"}}, false},
		{"remove floor", before, patch(removed, &gen.AlertMetadata{FloorCounts: map[string]uint32{"Floor1": 19}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Floor: "Floor1"}}, false},
		{"add zone", before, patch(added, &gen.AlertMetadata{ZoneCounts: map[string]uint32{"Zone1": 21}}), &gen.PullAlertsResponse_Change{NewValue: &gen.Alert{Zone: "Zone1"}}, false},
		{"remove zone", before, patch(removed, &gen.AlertMetadata{ZoneCounts: map[string]uint32{"Zone1": 19}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Zone: "Zone1"}}, false},
		{"add severity", before, patch(added, &gen.AlertMetadata{SeverityCounts: map[int32]uint32{1: 21}}), &gen.PullAlertsResponse_Change{NewValue: &gen.Alert{Severity: 1}}, false},
		{"remove severity", before, patch(removed, &gen.AlertMetadata{SeverityCounts: map[int32]uint32{1: 19}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Severity: 1}}, false},
		{"add subsystem", before, patch(added, &gen.AlertMetadata{SubsystemCounts: map[string]uint32{"Subsystem1": 21}}), &gen.PullAlertsResponse_Change{NewValue: &gen.Alert{Subsystem: "Subsystem1"}}, false},
		{"remove subsystem", before, patch(removed, &gen.AlertMetadata{SubsystemCounts: map[string]uint32{"Subsystem1": 19}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Subsystem: "Subsystem1"}}, false},

		{"update ack (ack->nak)", before, patch(before, &gen.AlertMetadata{
			AcknowledgedCounts:   map[bool]uint32{true: 79, false: 21},
			NeedsAttentionCounts: map[string]uint32{"nack_unresolved": 16, "ack_unresolved": 69},
		}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Acknowledgement: &gen.Alert_Acknowledgement{AcknowledgeTime: timestamppb.Now()}}, NewValue: &gen.Alert{}}, false},
		{"update ack (nak->ack)", before, patch(before, &gen.AlertMetadata{
			AcknowledgedCounts:   map[bool]uint32{true: 81, false: 19},
			NeedsAttentionCounts: map[string]uint32{"nack_unresolved": 14, "ack_unresolved": 71},
		}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{}, NewValue: &gen.Alert{Acknowledgement: &gen.Alert_Acknowledgement{AcknowledgeTime: timestamppb.Now()}}}, false},
		{"update floor", before, patch(before, &gen.AlertMetadata{FloorCounts: map[string]uint32{"Floor1": 19, "Floor2": 31}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Floor: "Floor1"}, NewValue: &gen.Alert{Floor: "Floor2"}}, false},
		{"update floor (zero->)", before, patch(before, &gen.AlertMetadata{FloorCounts: map[string]uint32{"Floor1": 21}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Floor: ""}, NewValue: &gen.Alert{Floor: "Floor1"}}, false},
		{"update floor (->zero)", before, patch(before, &gen.AlertMetadata{FloorCounts: map[string]uint32{"Floor1": 19}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Floor: "Floor1"}, NewValue: &gen.Alert{Floor: ""}}, false},
		{"update zone", before, patch(before, &gen.AlertMetadata{ZoneCounts: map[string]uint32{"Zone1": 19, "Zone2": 31}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Zone: "Zone1"}, NewValue: &gen.Alert{Zone: "Zone2"}}, false},
		{"update zone (zero->)", before, patch(before, &gen.AlertMetadata{ZoneCounts: map[string]uint32{"Zone1": 21}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Zone: ""}, NewValue: &gen.Alert{Zone: "Zone1"}}, false},
		{"update zone (->zero)", before, patch(before, &gen.AlertMetadata{ZoneCounts: map[string]uint32{"Zone1": 19}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Zone: "Zone1"}, NewValue: &gen.Alert{Zone: ""}}, false},
		{"update severity", before, patch(before, &gen.AlertMetadata{SeverityCounts: map[int32]uint32{1: 19, 2: 31}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Severity: 1}, NewValue: &gen.Alert{Severity: 2}}, false},
		{"update severity (zero->)", before, patch(before, &gen.AlertMetadata{SeverityCounts: map[int32]uint32{1: 21}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Severity: 0}, NewValue: &gen.Alert{Severity: 1}}, false},
		{"update severity (->zero)", before, patch(before, &gen.AlertMetadata{SeverityCounts: map[int32]uint32{1: 19}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Severity: 1}, NewValue: &gen.Alert{Severity: 0}}, false},
		{"update subsystem", before, patch(before, &gen.AlertMetadata{SubsystemCounts: map[string]uint32{"Subsystem1": 19, "Subsystem2": 31}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Subsystem: "Subsystem1"}, NewValue: &gen.Alert{Subsystem: "Subsystem2"}}, false},
		{"update subsystem (zero->)", before, patch(before, &gen.AlertMetadata{SubsystemCounts: map[string]uint32{"Subsystem1": 21}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Subsystem: ""}, NewValue: &gen.Alert{Subsystem: "Subsystem1"}}, false},
		{"update subsystem (->zero)", before, patch(before, &gen.AlertMetadata{SubsystemCounts: map[string]uint32{"Subsystem1": 19}}), &gen.PullAlertsResponse_Change{OldValue: &gen.Alert{Subsystem: "Subsystem1"}, NewValue: &gen.Alert{Subsystem: ""}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := resource.NewValue(resource.WithInitialValue(tt.before))
			err := ApplyMdDelta(res, tt.e)
			if (err != nil) != tt.wantErr {
				t.Fatalf("applyMdDelta error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			got := res.Get()
			if diff := cmp.Diff(tt.after, got, protocmp.Transform()); diff != "" {
				t.Fatalf("applyMdDelta md (-want,+got)\n%s", diff)
			}
		})
	}

	t.Run("empty old maps", func(t *testing.T) {
		res := resource.NewValue(resource.WithInitialValue(&gen.AlertMetadata{
			FloorCounts:          make(map[string]uint32),
			ZoneCounts:           make(map[string]uint32),
			SeverityCounts:       make(map[int32]uint32),
			AcknowledgedCounts:   make(map[bool]uint32),
			ResolvedCounts:       make(map[bool]uint32),
			NeedsAttentionCounts: make(map[string]uint32),
		}))
		err := ApplyMdDelta(res, &gen.PullAlertsResponse_Change{
			NewValue: &gen.Alert{Floor: "foo", Zone: "bar", Severity: 1},
		})
		if err != nil {
			t.Fatal(err)
		}
		got := res.Get()
		want := &gen.AlertMetadata{
			TotalCount:           1,
			FloorCounts:          map[string]uint32{"foo": 1},
			ZoneCounts:           map[string]uint32{"bar": 1},
			SeverityCounts:       map[int32]uint32{1: 1},
			AcknowledgedCounts:   map[bool]uint32{false: 1},
			ResolvedCounts:       map[bool]uint32{false: 1},
			NeedsAttentionCounts: map[string]uint32{"nack_unresolved": 1},
		}
		if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
			t.Fatalf("applyMdDelta md (-want,+got)\n%s", diff)
		}
	})
}

func patch(before, change *gen.AlertMetadata) *gen.AlertMetadata {
	dst := &gen.AlertMetadata{}
	proto.Merge(dst, before)
	proto.Merge(dst, change)
	return dst
}
