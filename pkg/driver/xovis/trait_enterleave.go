package xovis

import (
	"context"

	"github.com/smart-core-os/sc-api/go/traits"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vanti-dev/sc-bos/pkg/minibus"
)

type enterLeaveServer struct {
	traits.UnimplementedEnterLeaveSensorApiServer
	client      *Client
	logicID     int
	multiSensor bool
	bus         *minibus.Bus[LogicsPushData]
}

func (e *enterLeaveServer) PullEnterLeaveEvents(
	request *traits.PullEnterLeaveEventsRequest, server traits.EnterLeaveSensorApi_PullEnterLeaveEventsServer,
) error {
	// get the initial value of the logics so we can compare later
	res, err := GetLiveLogic(e.client, e.multiSensor, e.logicID)
	if err != nil {
		return status.Error(codes.Unavailable, err.Error())
	}

	fwID, fwOK := findCountIDByName(res.Logic.Counts, "fw")
	bwID, bwOK := findCountIDByName(res.Logic.Counts, "bw")
	if !fwOK || !bwOK {
		return status.Error(codes.FailedPrecondition,
			"Counts don't match expected structure; check that this is an InOut logic")
	}
	// these will definitely work because we found these ID in the same slice above
	fwAcc, _ := findCountValueByID(res.Logic.Counts, fwID)
	bwAcc, _ := findCountValueByID(res.Logic.Counts, bwID)

	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()
	for pushData := range e.bus.Listen(ctx) {
		var events []*traits.PullEnterLeaveEventsResponse_Change
		for _, logic := range pushData.Logics {
			if logic.ID != e.logicID {
				continue
			}
			events = append(events, logicRecordsToChanges(logic.Records)...)
		}

		err = server.Send(&traits.PullEnterLeaveEventsResponse{
			Changes: events,
		})
	}

	return status.Error(codes.Unimplemented, "PullEnterLeaveEvents unimplemented for this device")
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

func decodeInOutCounts(counts []Count) (fw, bw int, ok bool) {
	var foundFw, foundBw bool
	for _, count := range counts {
		if count.Name == "fw" {
			fw = count.Value
			foundFw = true
		} else if count.Name == "bw" {
			bw = count.Value
			foundBw = true
		}
	}
	ok = foundFw && foundBw
	return
}

func logicRecordsToChanges(records []LogicRecord, fwID, bwID int) []*traits.PullEnterLeaveEventsResponse_Change {
	var changes []*traits.PullEnterLeaveEventsResponse_Change
	for _, record := range records {

	}
	return changes
}
