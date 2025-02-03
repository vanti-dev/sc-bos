package history

import (
	"testing"
	"time"
)

func TestRecord_CompareFrom(t *testing.T) {
	t1, t2 := time.Unix(100, 0), time.Unix(200, 0)
	tests := []struct {
		name string
		a, b Record
		want int
	}{
		{"empty==empty", Record{}, Record{}, 0},
		{"empty<id", Record{}, Record{ID: "a"}, -1},
		{"id>empty", Record{ID: "a"}, Record{}, 1},
		{"id==id", Record{ID: "a"}, Record{ID: "a"}, 0},
		{"id<id", Record{ID: "a"}, Record{ID: "b"}, -1},
		{"id>id", Record{ID: "b"}, Record{ID: "a"}, 1},
		{"empty<time", Record{}, Record{CreateTime: t1}, -1},
		{"time>empty", Record{CreateTime: t1}, Record{}, 1},
		{"time==time", Record{CreateTime: t1}, Record{CreateTime: t1}, 0},
		{"time<time", Record{CreateTime: t1}, Record{CreateTime: t2}, -1},
		{"time>time", Record{CreateTime: t2}, Record{CreateTime: t1}, 1},
		{"time>id", Record{CreateTime: t1}, Record{ID: "a"}, 1}, // zero time trumps zero id
		{"id<time", Record{ID: "a"}, Record{CreateTime: t1}, -1},
		{"time==time,id>empty", Record{CreateTime: t1, ID: "a"}, Record{CreateTime: t1}, 1},
		{"time==time,empty<id", Record{CreateTime: t1}, Record{CreateTime: t1, ID: "a"}, -1},
		{"time==time,id<id", Record{CreateTime: t1, ID: "a"}, Record{CreateTime: t1, ID: "b"}, -1},
		{"time==time,id>id", Record{CreateTime: t1, ID: "b"}, Record{CreateTime: t1, ID: "a"}, 1},
		{"time==time,id==id", Record{CreateTime: t1, ID: "a"}, Record{CreateTime: t1, ID: "a"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Compare(tt.b); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_CompareTo(t *testing.T) {
	t1, t2 := time.Unix(100, 0), time.Unix(200, 0)
	tests := []struct {
		name string
		a, b Record
		want int
	}{
		{"empty==empty", Record{}, Record{}, 0},
		{"empty>id", Record{}, Record{ID: "a"}, 1},
		{"id<empty", Record{ID: "a"}, Record{}, -1},
		{"id==id", Record{ID: "a"}, Record{ID: "a"}, 0},
		{"id<id", Record{ID: "a"}, Record{ID: "b"}, -1},
		{"id>id", Record{ID: "b"}, Record{ID: "a"}, 1},
		{"empty>time", Record{}, Record{CreateTime: t1}, 1},
		{"time<empty", Record{CreateTime: t1}, Record{}, -1},
		{"time==time", Record{CreateTime: t1}, Record{CreateTime: t1}, 0},
		{"time<time", Record{CreateTime: t1}, Record{CreateTime: t2}, -1},
		{"time>time", Record{CreateTime: t2}, Record{CreateTime: t1}, 1},
		{"time<id", Record{CreateTime: t1}, Record{ID: "a"}, -1}, // zero time trumps zero id
		{"id>time", Record{ID: "a"}, Record{CreateTime: t1}, 1},
		{"time==time,id<empty", Record{CreateTime: t1, ID: "a"}, Record{CreateTime: t1}, -1},
		{"time==time,empty>id", Record{CreateTime: t1}, Record{CreateTime: t1, ID: "a"}, 1},
		{"time==time,id<id", Record{CreateTime: t1, ID: "a"}, Record{CreateTime: t1, ID: "b"}, -1},
		{"time==time,id>id", Record{CreateTime: t1, ID: "b"}, Record{CreateTime: t1, ID: "a"}, 1},
		{"time==time,id==id", Record{CreateTime: t1, ID: "a"}, Record{CreateTime: t1, ID: "a"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.CompareZeroAfter(tt.b); got != tt.want {
				t.Errorf("CompareZeroAfter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersectRecords(t *testing.T) {
	tests := []struct {
		name                   string
		f1, t1, f2, t2, wf, wt Record
	}{
		{"[-,-],[-,-]=>[-,-]", Record{}, Record{}, Record{}, Record{}, Record{}, Record{}},
		{"[-,-],[-,1]=>[-,1]", Record{}, Record{}, Record{}, Record{ID: "1"}, Record{}, Record{ID: "1"}},
		{"[-,1],[-,1]=>[-,1]", Record{}, Record{ID: "1"}, Record{}, Record{ID: "1"}, Record{}, Record{ID: "1"}},
		{"[-,1],[-,2]=>[-,1]", Record{}, Record{ID: "1"}, Record{}, Record{ID: "2"}, Record{}, Record{ID: "1"}},
		{"[1,-],[-,-]=>[1,-]", Record{ID: "1"}, Record{}, Record{}, Record{}, Record{ID: "1"}, Record{}},
		{"[1,-],[1,-]=>[1,-]", Record{ID: "1"}, Record{}, Record{ID: "1"}, Record{}, Record{ID: "1"}, Record{}},
		{"[1,-],[2,-]=>[2,-]", Record{ID: "1"}, Record{}, Record{ID: "2"}, Record{}, Record{ID: "2"}, Record{}},
		{"[1,2],[0,1]=>[1,1]", Record{ID: "1"}, Record{ID: "2"}, Record{ID: "0"}, Record{ID: "1"}, Record{ID: "1"}, Record{ID: "1"}},
		{"[1,2],[1,2]=>[1,2]", Record{ID: "1"}, Record{ID: "2"}, Record{ID: "1"}, Record{ID: "2"}, Record{ID: "1"}, Record{ID: "2"}},
		{"[2,4],[1,2]=>[2,2]", Record{ID: "2"}, Record{ID: "4"}, Record{ID: "1"}, Record{ID: "2"}, Record{ID: "2"}, Record{ID: "2"}},
		{"[2,4],[2,3]=>[2,3]", Record{ID: "2"}, Record{ID: "4"}, Record{ID: "2"}, Record{ID: "3"}, Record{ID: "2"}, Record{ID: "3"}},
		{"[2,4],[3,4]=>[3,4]", Record{ID: "2"}, Record{ID: "4"}, Record{ID: "3"}, Record{ID: "4"}, Record{ID: "3"}, Record{ID: "4"}},
		{"[2,4],[4,5]=>[4,4]", Record{ID: "2"}, Record{ID: "4"}, Record{ID: "4"}, Record{ID: "5"}, Record{ID: "4"}, Record{ID: "4"}},
		{"[2,4],[5,6]=>[5,4]", Record{ID: "2"}, Record{ID: "4"}, Record{ID: "5"}, Record{ID: "6"}, Record{ID: "5"}, Record{ID: "4"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gf, gt := IntersectRecords(tt.f1, tt.t1, tt.f2, tt.t2); gf.Compare(tt.wf) != 0 || gt.CompareZeroAfter(tt.wt) != 0 {
				t.Errorf("IntersectRecords() = %v,%v, want %v,%v", gf, gt, tt.wf, tt.wt)
			}
			// same again but testing commutativity
			if gf, gt := IntersectRecords(tt.f2, tt.t2, tt.f1, tt.t1); gf.Compare(tt.wf) != 0 || gt.CompareZeroAfter(tt.wt) != 0 {
				t.Errorf("IntersectRecords() = %v,%v, want %v,%v", gf, gt, tt.wf, tt.wt)
			}
		})
	}
}
