package hikcentral

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/vanti-dev/sc-bos/pkg/driver/hikcentral/api"
)

func TestANPRController_getMostRelevantAppointment(t *testing.T) {
	type parameters struct {
		now          time.Time
		appointments []*api.Appointment
		expected     *api.Appointment
	}

	tests := map[string]parameters{
		"one appointment": {
			now: time.Date(2025, 6, 18, 12, 0, 0, 0, time.UTC),
			appointments: []*api.Appointment{
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 0, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 13, 0, 0, 0, time.UTC)),
				},
			},
			expected: &api.Appointment{
				AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 0, 0, 0, time.UTC)),
				AppointEndTime:   formatTime(time.Date(2025, 6, 18, 13, 0, 0, 0, time.UTC)),
			},
		},
		" 2 appointments - ignore end time": {
			now: time.Date(2025, 6, 18, 12, 0, 0, 0, time.UTC),
			appointments: []*api.Appointment{
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 0, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 12, 30, 0, 0, time.UTC)),
				},
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 10, 30, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 13, 30, 0, 0, time.UTC)),
				},
			},
			expected: &api.Appointment{
				AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 0, 0, 0, time.UTC)),
				AppointEndTime:   formatTime(time.Date(2025, 6, 18, 12, 30, 0, 0, time.UTC)),
			},
		},
		"2 appointments - closer start time": {
			now: time.Date(2025, 6, 18, 12, 0, 0, 0, time.UTC),
			appointments: []*api.Appointment{
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 0, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 13, 0, 0, 0, time.UTC)),
				},
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 30, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 13, 0, 0, 0, time.UTC)),
				},
			},
			expected: &api.Appointment{
				AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 30, 0, 0, time.UTC)),
				AppointEndTime:   formatTime(time.Date(2025, 6, 18, 13, 0, 0, 0, time.UTC)),
			},
		},
		"2 appointments - further end time": {
			now: time.Date(2025, 6, 18, 12, 0, 0, 0, time.UTC),
			appointments: []*api.Appointment{
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 0, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 13, 0, 0, 0, time.UTC)),
				},
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 0, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 14, 0, 0, 0, time.UTC)),
				},
			},
			expected: &api.Appointment{
				AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 0, 0, 0, time.UTC)),
				AppointEndTime:   formatTime(time.Date(2025, 6, 18, 14, 0, 0, 0, time.UTC)),
			},
		},
		"2 appointments - pick longest appointment": {
			now: time.Date(2025, 6, 18, 12, 0, 0, 0, time.UTC),
			appointments: []*api.Appointment{
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 0, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 12, 30, 0, 0, time.UTC)),
				},
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 30, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 13, 0, 0, 0, time.UTC)),
				},
			},
			expected: &api.Appointment{
				AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 30, 0, 0, time.UTC)),
				AppointEndTime:   formatTime(time.Date(2025, 6, 18, 13, 0, 0, 0, time.UTC)),
			},
		},
		"2 appointments - end times are irrelevant": {
			now: time.Date(2025, 6, 18, 12, 0, 0, 0, time.UTC),
			appointments: []*api.Appointment{
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 0, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 12, 30, 0, 0, time.UTC)),
				},
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 30, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 12, 15, 0, 0, time.UTC)),
				},
			},
			expected: &api.Appointment{
				AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 30, 0, 0, time.UTC)),
				AppointEndTime:   formatTime(time.Date(2025, 6, 18, 12, 15, 0, 0, time.UTC)),
			},
		},
		"3 appointments - pick longest appointment": {
			now: time.Date(2025, 6, 18, 12, 0, 0, 0, time.UTC),
			appointments: []*api.Appointment{
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 0, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 12, 30, 0, 0, time.UTC)),
				},
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 30, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 12, 30, 0, 0, time.UTC)),
				},
				{
					AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 0, 0, 0, time.UTC)),
					AppointEndTime:   formatTime(time.Date(2025, 6, 18, 13, 0, 0, 0, time.UTC)),
				},
			},
			expected: &api.Appointment{
				AppointStartTime: formatTime(time.Date(2025, 6, 18, 11, 0, 0, 0, time.UTC)),
				AppointEndTime:   formatTime(time.Date(2025, 6, 18, 13, 0, 0, 0, time.UTC)),
			},
		},
	}

	for name, params := range tests {
		t.Run(name, func(t *testing.T) {
			ret := getMostRelevantAppointment(params.now, params.appointments)

			if diff := cmp.Diff(params.expected, ret); diff != "" {
				t.Errorf("getMostRelevantAppointment() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestANPRController_formatPlateNo(t *testing.T) {
	tests := map[string]struct {
		plateNo  string
		expected string
	}{
		"empty plate": {
			plateNo:  "",
			expected: "",
		},
		"valid plate": {
			plateNo:  "ABC123",
			expected: "ABC123",
		},
		"plate with spaces": {
			plateNo:  "A B C 1 2 3",
			expected: "ABC123",
		},
		"plate with special characters": {
			plateNo:  "A-B_C@123",
			expected: "ABC123",
		},
		"plate with mixed case": {
			plateNo:  "aBc123",
			expected: "ABC123",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := formatPlateNo(test.plateNo)
			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
}
