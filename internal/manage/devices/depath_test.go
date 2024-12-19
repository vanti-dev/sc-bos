package devices

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_depath(t *testing.T) {
	tests := []struct {
		path              string
		deconstructedPath dePath
	}{
		{
			path: "something[0]",
			deconstructedPath: dePath{
				Before: "something",
				After:  "[0].",
				Found:  true,
				Index:  0,
				Next:   "",
			},
		}, {
			path: "something[10].else",
			deconstructedPath: dePath{
				Before: "something",
				After:  "[10].else",
				Found:  true,
				Index:  10,
				Next:   "else",
			},
		}, {
			path: "something[-1].else",
			deconstructedPath: dePath{
				Before: "something",
				After:  "else",
				Found:  true,
				Index:  -1,
				Next:   "else",
			},
		}, {
			path: "something[x20].else",
			deconstructedPath: dePath{
				Before: "something[x20]",
				After:  "else",
				Found:  false,
				Index:  -1,
				Next:   "else",
			},
		}, {
			path: "real[0].life",
			deconstructedPath: dePath{
				Before: "real",
				After:  "[0].life",
				Found:  true,
				Index:  0,
				Next:   "life",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := depath(tt.path)

			if diff := cmp.Diff(tt.deconstructedPath, got); diff != "" {
				t.Errorf("depath(%q): -want +got:\n%s", tt.path, diff)
			}
		})
	}
}
