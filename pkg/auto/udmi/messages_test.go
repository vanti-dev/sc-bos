package udmi

import "testing"

func TestPointsEvent_Equal(t *testing.T) {
	tests := []struct {
		name  string
		e1    PointsEvent
		e2    PointsEvent
		equal bool
	}{
		{
			name:  "empty",
			e1:    nil,
			e2:    nil,
			equal: true,
		},
		{
			name:  "one is nil",
			e1:    nil,
			e2:    map[string]PointValue{},
			equal: false,
		},
		{
			name: "matching values",
			e1: map[string]PointValue{
				"p1": {PresentValue: true},
				"p2": {PresentValue: "my-string"},
				"p3": {PresentValue: 10.67},
			},
			e2: map[string]PointValue{
				"p2": {PresentValue: "my-string"},
				"p1": {PresentValue: true},
				"p3": {PresentValue: 10.67},
			},
			equal: true,
		},
		{
			name: "mismatched values",
			e1: map[string]PointValue{
				"p1": {PresentValue: true},
				"p2": {PresentValue: "my-string"},
				"p3": {PresentValue: 10.67},
			},
			e2: map[string]PointValue{
				"p2": {PresentValue: "my-string2"},
				"p1": {PresentValue: true},
				"p3": {PresentValue: 10.67},
			},
			equal: false,
		},
		{
			name: "subset",
			e1: map[string]PointValue{
				"p1": {PresentValue: true},
				"p2": {PresentValue: "my-string"},
				"p3": {PresentValue: 10.67},
			},
			e2: map[string]PointValue{
				"p2": {PresentValue: "my-string"},
				"p1": {PresentValue: true},
				"p3": {PresentValue: 10.67},
				"p4": {PresentValue: false},
			},
			equal: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e1.Equal(tt.e2)
			if got != tt.equal {
				t.Fatalf("expected %t, got %t", tt.equal, got)
			}
			if tt.e1.Equal(tt.e2) != tt.e2.Equal(tt.e1) {
				t.Fatalf("commutability broken e1==e2:%v != e2==e1:%v", tt.e1.Equal(tt.e2), tt.e2.Equal(tt.e1))
			}
		})
	}
}
