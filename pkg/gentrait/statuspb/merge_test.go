package statuspb

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
)

func TestProblemMerger(t *testing.T) {
	t.Run("highest level (nothing nominal)", func(t *testing.T) {
		now := time.Unix(0, 0)
		pm := &ProblemMerger{}
		pm.AddProblem(&gen.StatusLog_Problem{
			Name:        "p2",
			Level:       gen.StatusLog_NOTICE,
			Description: "p2 is a notice",
			RecordTime:  timestamppb.New(now.Add(-2 * time.Second)),
		})
		pm.AddProblem(&gen.StatusLog_Problem{
			Name:        "p1",
			Level:       gen.StatusLog_NON_FUNCTIONAL,
			Description: "p1 is non-functional",
			RecordTime:  timestamppb.New(now.Add(-time.Second)),
		})
		sl := pm.Build()
		want := &gen.StatusLog{
			Level:       gen.StatusLog_NON_FUNCTIONAL,
			Description: "p1 is non-functional",
			RecordTime:  timestamppb.New(now.Add(-time.Second)),
			Problems: []*gen.StatusLog_Problem{
				{
					Name:        "p2",
					Level:       gen.StatusLog_NOTICE,
					Description: "p2 is a notice",
					RecordTime:  timestamppb.New(now.Add(-2 * time.Second)),
				},
				{
					Name:        "p1",
					Level:       gen.StatusLog_NON_FUNCTIONAL,
					Description: "p1 is non-functional",
					RecordTime:  timestamppb.New(now.Add(-time.Second)),
				},
			},
		}
		if diff := cmp.Diff(want, sl, protocmp.Transform()); diff != "" {
			t.Errorf("unexpected diff (-want +got):\n%s", diff)
		}
	})

	t.Run("some offline", func(t *testing.T) {
		now := time.Unix(0, 0)
		pm := &ProblemMerger{}
		pm.AddProblem(&gen.StatusLog_Problem{
			Name:        "p2",
			Level:       gen.StatusLog_NOMINAL,
			Description: "p2 nominal",
			RecordTime:  timestamppb.New(now.Add(-2 * time.Second)),
		})
		pm.AddProblem(&gen.StatusLog_Problem{
			Name:        "p1",
			Level:       gen.StatusLog_OFFLINE,
			Description: "p1 is offline",
			RecordTime:  timestamppb.New(now.Add(-time.Second)),
		})
		sl := pm.Build()
		want := &gen.StatusLog{
			Level:       gen.StatusLog_REDUCED_FUNCTION,
			Description: "p1 is offline",
			RecordTime:  timestamppb.New(now.Add(-time.Second)),
			Problems: []*gen.StatusLog_Problem{
				{
					Name:        "p1",
					Level:       gen.StatusLog_OFFLINE,
					Description: "p1 is offline",
					RecordTime:  timestamppb.New(now.Add(-time.Second)),
				},
			},
		}
		if diff := cmp.Diff(want, sl, protocmp.Transform()); diff != "" {
			t.Errorf("unexpected diff (-want +got):\n%s", diff)
		}
	})

	t.Run("all offline", func(t *testing.T) {
		now := time.Unix(0, 0)
		pm := &ProblemMerger{}
		pm.AddProblem(&gen.StatusLog_Problem{
			Name:        "p2",
			Level:       gen.StatusLog_OFFLINE,
			Description: "p2 is offline",
			RecordTime:  timestamppb.New(now.Add(-2 * time.Second)),
		})
		pm.AddProblem(&gen.StatusLog_Problem{
			Name:        "p1",
			Level:       gen.StatusLog_OFFLINE,
			Description: "p1 is offline",
			RecordTime:  timestamppb.New(now.Add(-time.Second)),
		})
		sl := pm.Build()
		want := &gen.StatusLog{
			Level:       gen.StatusLog_OFFLINE,
			Description: "p2 is offline",
			RecordTime:  timestamppb.New(now.Add(-2 * time.Second)),
			Problems: []*gen.StatusLog_Problem{
				{
					Name:        "p2",
					Level:       gen.StatusLog_OFFLINE,
					Description: "p2 is offline",
					RecordTime:  timestamppb.New(now.Add(-2 * time.Second)),
				},
				{
					Name:        "p1",
					Level:       gen.StatusLog_OFFLINE,
					Description: "p1 is offline",
					RecordTime:  timestamppb.New(now.Add(-time.Second)),
				},
			},
		}
		if diff := cmp.Diff(want, sl, protocmp.Transform()); diff != "" {
			t.Errorf("unexpected diff (-want +got):\n%s", diff)
		}
	})

	t.Run("nominal", func(t *testing.T) {
		now := time.Unix(0, 0)
		pm := &ProblemMerger{}
		pm.AddStatusLog(&gen.StatusLog{
			Level:       gen.StatusLog_NOMINAL,
			Description: "p2 is nominal",
			RecordTime:  timestamppb.New(now.Add(-2 * time.Second)),
		})
		pm.AddStatusLog(&gen.StatusLog{
			Level:       gen.StatusLog_NOMINAL,
			Description: "p1 is nominal",
			RecordTime:  timestamppb.New(now.Add(-time.Second)),
		})
		sl := pm.Build()
		want := &gen.StatusLog{
			Level:       gen.StatusLog_NOMINAL,
			Description: "p1 is nominal",
			RecordTime:  timestamppb.New(now.Add(-time.Second)),
		}
		if diff := cmp.Diff(want, sl, protocmp.Transform()); diff != "" {
			t.Errorf("unexpected diff (-want +got):\n%s", diff)
		}
	})
}
