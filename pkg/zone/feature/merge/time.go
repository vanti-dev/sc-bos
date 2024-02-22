package merge

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EarliestTimestamp[E any](items []E, f func(E) *timestamppb.Timestamp) *timestamppb.Timestamp {
	var res *timestamppb.Timestamp
	for _, item := range items {
		if v := f(item); v != nil {
			if res == nil || v.AsTime().Before(res.AsTime()) {
				res = v
			}
		}
	}
	return res
}

func LatestTimestamp[E any](items []E, f func(E) *timestamppb.Timestamp) *timestamppb.Timestamp {
	var res *timestamppb.Timestamp
	for _, item := range items {
		if v := f(item); v != nil {
			if res == nil || v.AsTime().After(res.AsTime()) {
				res = v
			}
		}
	}
	return res
}

// MaxDuration returns the maximum duration from the given items.
func MaxDuration[E any](items []E, f func(E) *durationpb.Duration) *durationpb.Duration {
	var res *durationpb.Duration
	for _, item := range items {
		if v := f(item); v != nil {
			if res == nil || v.AsDuration() > res.AsDuration() {
				res = v
			}
		}
	}
	return res
}
