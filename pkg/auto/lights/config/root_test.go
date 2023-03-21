package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name string
		dd   *DaylightDimming
		want []LevelThreshold
	}{
		{
			name: "empty",
			dd:   nil,
			want: nil,
		},
		{
			name: "has thresholds",
			dd: &DaylightDimming{
				// this is invalid, but ignored
				Segments: &ThresholdSegments{
					MaxLevel: 2500,
				},
				Thresholds: []LevelThreshold{
					{BelowLux: 500, LevelPercent: 43},
					{BelowLux: 250, LevelPercent: 50},
				},
			},
			want: []LevelThreshold{
				{BelowLux: 500, LevelPercent: 43},
				{BelowLux: 250, LevelPercent: 50},
			},
		},
		{
			name: "generates thresholds",
			dd: &DaylightDimming{
				// this is invalid, but ignored
				Segments: &ThresholdSegments{
					Steps:  2,
					MinLux: 100,
					MaxLux: 1000,
				},
			},
			want: []LevelThreshold{
				{BelowLux: 1000, LevelPercent: 0},
				{BelowLux: 550, LevelPercent: 100},
			},
		},
		{
			name: "custom levels",
			dd: &DaylightDimming{
				// this is invalid, but ignored
				Segments: &ThresholdSegments{
					Steps:    2,
					MinLux:   100,
					MaxLux:   1000,
					MinLevel: 25,
					MaxLevel: 75,
				},
			},
			want: []LevelThreshold{
				{BelowLux: 1000, LevelPercent: 25},
				{BelowLux: 550, LevelPercent: 75},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dd.process()
			if err != nil {
				t.Fatalf("%s", err)
			}
			var got []LevelThreshold
			if tt.dd != nil {
				got = tt.dd.Thresholds
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("thresholds (-want,+got)\n%s", diff)
			}
		})
	}
}
