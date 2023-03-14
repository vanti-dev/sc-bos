package xovis

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/vanti-dev/sc-bos/pkg/minibus"
)

type enterLeaveServer struct {
	traits.UnimplementedEnterLeaveSensorApiServer
	client      *Client
	logicID     int
	multiSensor bool
	bus         *minibus.Bus[PushData]
}

func (e *enterLeaveServer) PullEnterLeaveEvents(request *traits.PullEnterLeaveEventsRequest, server traits.EnterLeaveSensorApi_PullEnterLeaveEventsServer) error {
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
	forwardCount, _ := findCountValueByID(res.Logic.Counts, fwID)
	backwardCount, _ := findCountValueByID(res.Logic.Counts, bwID)

	accumulator := countAccumulator{
		forwardCountID:     fwID,
		backwardCountID:    bwID,
		forwardCountValue:  forwardCount,
		backwardCountValue: backwardCount,
	}

	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()
	for pushData := range e.bus.Listen(ctx) {
		if pushData.LogicsData == nil {
			continue
		}
		records, ok := findLogicRecords(pushData.LogicsData, e.logicID)
		if !ok {
			continue
		}

		events, err := accumulator.consumeRecords(records)
		if err != nil {
			return err
		}

		var enterLeaveChanges []*traits.PullEnterLeaveEventsResponse_Change
		for _, event := range events {
			enterLeaveChanges = append(enterLeaveChanges, &traits.PullEnterLeaveEventsResponse_Change{
				Name:            request.Name,
				ChangeTime:      timestamppb.New(event.time),
				EnterLeaveEvent: &traits.EnterLeaveEvent{Direction: event.direction},
			})
		}

		err = server.Send(&traits.PullEnterLeaveEventsResponse{
			Changes: enterLeaveChanges,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func findLogicRecords(data *LogicsPushData, logicID int) (records []LogicRecord, ok bool) {
	for _, logic := range data.Logics {
		if logic.ID == logicID {
			records = logic.Records
			ok = true
			return
		}
	}
	return
}
