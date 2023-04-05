package merge

import (
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
