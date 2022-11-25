package testlight

import (
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestGenerateReport(t *testing.T) {
	db, cleanup := prepareEmptyTestDB()
	defer cleanup()

	baseTime := time.Date(2022, time.November, 25, 11, 0, 0, 0, time.UTC)
	// test-1 has a function test pass and 2 duration test passes
	// test-2 has a function test pass only
	// test-3 has no passes
	latestStatuses := []LatestStatusRecord{
		{
			Name:       "test-1",
			LastUpdate: baseTime.Add(time.Hour),
			Faults:     nil,
		},
		{
			Name:       "test-2",
			LastUpdate: baseTime.Add(time.Hour + time.Minute),
			Faults:     []gen.EmergencyLightFault{gen.EmergencyLightFault_LAMP_FAULT},
		},
		{
			Name:       "test-3",
			LastUpdate: baseTime.Add(time.Hour + 2*time.Minute),
			Faults:     []gen.EmergencyLightFault{gen.EmergencyLightFault_FUNCTION_TEST_FAILED, gen.EmergencyLightFault_DURATION_TEST_FAILED},
		},
	}
	events := []EventRecord{
		{
			ID:        1,
			Name:      "test-1",
			Timestamp: baseTime,
			Kind:      FunctionTestPassEvent,
		},
		{
			ID:        2,
			Name:      "test-2",
			Timestamp: baseTime.Add(time.Minute),
			Kind:      FunctionTestPassEvent,
		},
		{
			ID:        3,
			Name:      "test-1",
			Timestamp: baseTime.Add(2 * time.Minute),
			Kind:      DurationTestPassEvent,
			DurationTestPass: &gen.EmergencyLightingEvent_DurationTestPass{
				AchievedDuration: durationpb.New(3 * time.Hour),
			},
		},
		{
			ID:        4,
			Name:      "test-1",
			Timestamp: baseTime.Add(time.Hour),
			Kind:      DurationTestPassEvent,
			DurationTestPass: &gen.EmergencyLightingEvent_DurationTestPass{
				AchievedDuration: durationpb.New(2 * time.Hour),
			},
		},
	}
	addLatestStatus(t, db, latestStatuses)
	addEvents(t, db, events)

	report, err := GenerateReport(db)
	if err != nil {
		t.Fatalf("failed to generate report: %s", err.Error())
	}

	expect := []ReportEntry{
		{
			Name:                   "test-1",
			LastUpdate:             latestStatuses[0].LastUpdate,
			Faults:                 nil,
			LatestFunctionTestPass: events[0].Timestamp,
			LatestDurationTestPass: events[3].Timestamp,
		},
		{
			Name:                   "test-2",
			LastUpdate:             latestStatuses[1].LastUpdate,
			Faults:                 latestStatuses[1].Faults,
			LatestFunctionTestPass: events[1].Timestamp,
		},
		{
			Name:       "test-3",
			LastUpdate: latestStatuses[2].LastUpdate,
			Faults:     latestStatuses[2].Faults,
		},
	}

	diff := cmp.Diff(expect, report, protocmp.Transform(), cmpopts.EquateEmpty())
	if diff != "" {
		t.Errorf("report mismatch (-want +got):\n%s", diff)
	}

	t.Run("Sorted", func(t *testing.T) {
		sorted := sort.SliceIsSorted(report, func(i, j int) bool {
			return report[i].Name < report[j].Name
		})

		if !sorted {
			t.Error("slice is not sorted")
		}
	})
}
