package xovis

import (
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type countAccumulator struct {
	forwardCountID     int
	backwardCountID    int
	forwardCountValue  int
	backwardCountValue int
}

type countEvent struct {
	time time.Time
	// We associate the Xovis API's "forward" events with "enter" events in Smart Core,
	// and Xovis "backward" events with "leave" events in Smart Core.
	direction traits.EnterLeaveEvent_Direction
}

func (c *countAccumulator) consumeRecords(records []LogicRecord) ([]countEvent, error) {
	var events []countEvent
	for _, record := range records {
		newForwardCount, ok := findCountValueByID(record.Counts, c.forwardCountID)
		if ok {
			delta := newForwardCount - c.forwardCountValue
			if delta < 0 {
				return nil, status.Error(codes.DataLoss, "logic forwards counter desynchronised")
			}

			event := countEvent{
				time:      record.To,
				direction: traits.EnterLeaveEvent_ENTER,
			}
			for i := 0; i < delta; i++ {
				events = append(events, event)
			}

			c.forwardCountValue = newForwardCount
		}

		newBackwardCount, ok := findCountValueByID(record.Counts, c.backwardCountID)
		if ok {
			delta := newBackwardCount - c.backwardCountValue
			if delta < 0 {
				return nil, status.Error(codes.DataLoss, "logic backwards counter desynchronised")
			}

			event := countEvent{
				time:      record.To,
				direction: traits.EnterLeaveEvent_LEAVE,
			}
			for i := 0; i < delta; i++ {
				events = append(events, event)
			}

			c.backwardCountValue = newBackwardCount
		}
	}
	return events, nil
}

func findCountIDByName(counts []Count, name string) (id int, ok bool) {
	for _, count := range counts {
		if count.Name == name {
			return count.ID, true
		}
	}
	return
}

func findCountValueByID(counts []Count, id int) (value int, ok bool) {
	for _, count := range counts {
		if count.ID == id {
			return count.Value, true
		}
	}
	return
}
