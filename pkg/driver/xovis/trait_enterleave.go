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

func (e *enterLeaveServer) GetEnterLeaveEvent(ctx context.Context, request *traits.GetEnterLeaveEventRequest) (*traits.EnterLeaveEvent, error) {
	res, err := GetLiveLogic(e.client, e.multiSensor, e.logicID)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	_, forwardCount, fwOK := findCountByName(res.Logic.Counts, "fw")
	_, backwardCount, bwOK := findCountByName(res.Logic.Counts, "bw")
	if !fwOK || !bwOK {
		return nil, status.Error(codes.FailedPrecondition,
			"Counts don't match expected structure; check that this is an InOut logic")
	}

	forwardCount32, backwardCount32 := int32(forwardCount), int32(backwardCount)
	return &traits.EnterLeaveEvent{
		EnterTotal: &forwardCount32,
		LeaveTotal: &backwardCount32,
	}, nil
}

func (e *enterLeaveServer) ResetEnterLeaveTotals(ctx context.Context, request *traits.ResetEnterLeaveTotalsRequest) (*traits.ResetEnterLeaveTotalsResponse, error) {
	return nil, ResetLiveLogic(e.client, e.multiSensor, e.logicID)
}

func (e *enterLeaveServer) PullEnterLeaveEvents(request *traits.PullEnterLeaveEventsRequest, server traits.EnterLeaveSensorApi_PullEnterLeaveEventsServer) error {
	// get the initial value of the logics so we can compare later
	res, err := GetLiveLogic(e.client, e.multiSensor, e.logicID)
	if err != nil {
		return status.Error(codes.Unavailable, err.Error())
	}

	fwID, forwardCount, fwOK := findCountByName(res.Logic.Counts, "fw")
	bwID, backwardCount, bwOK := findCountByName(res.Logic.Counts, "bw")
	if !fwOK || !bwOK {
		return status.Error(codes.FailedPrecondition,
			"Counts don't match expected structure; check that this is an InOut logic")
	}

	if !request.UpdatesOnly {
		enterTotal, leaveTotal := int32(forwardCount), int32(backwardCount)
		err := server.Send(&traits.PullEnterLeaveEventsResponse{Changes: []*traits.PullEnterLeaveEventsResponse_Change{
			{
				Name:       request.Name,
				ChangeTime: timestamppb.Now(),
				EnterLeaveEvent: &traits.EnterLeaveEvent{
					EnterTotal: &enterTotal,
					LeaveTotal: &leaveTotal,
				},
			},
		}})
		if err != nil {
			return err
		}
	}

	// note: the accumulator continues to count totals even if the sensor is reset, for as long as the stream is active.
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

		// note: accumulator totals are updated during consumeRecords. We want the values before this happens
		enterTotal, leaveTotal := int32(accumulator.forwardCountValue), int32(accumulator.backwardCountValue)
		events, err := accumulator.consumeRecords(records)
		if err != nil {
			return err
		}

		if len(events) == 0 {
			continue
		}

		var enterLeaveChanges []*traits.PullEnterLeaveEventsResponse_Change
		for _, event := range events {
			switch event.direction {
			case traits.EnterLeaveEvent_ENTER:
				enterTotal++
			case traits.EnterLeaveEvent_LEAVE:
				leaveTotal++
			}
			enterLeaveChanges = append(enterLeaveChanges, &traits.PullEnterLeaveEventsResponse_Change{
				Name:       request.Name,
				ChangeTime: timestamppb.New(event.time),
				EnterLeaveEvent: &traits.EnterLeaveEvent{
					Direction:  event.direction,
					EnterTotal: &enterTotal,
					LeaveTotal: &leaveTotal,
				},
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
