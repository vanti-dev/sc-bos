package merge

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/comm"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/status"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	transportpb "github.com/vanti-dev/sc-bos/pkg/gentrait/transport"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

type transportCfg struct {
	config.Trait
	AssignedLandingCalls *config.ValueSource `json:"assignedLandingCalls,omitempty"`
	MakingCarCall        *config.ValueSource `json:"makingCarCall,omitempty"`
	CarPosition          *config.ValueSource `json:"carPosition,omitempty"`
	CarAssignedDirection *config.ValueSource `json:"carAssignedDirection,omitempty"`
	CarDoorStatus        *config.ValueSource `json:"carDoorStatus,omitempty"`
	CarMode              *config.ValueSource `json:"carMode,omitempty"`
	CarLoad              *config.ValueSource `json:"carLoad,omitempty"`
	CarLoadUnits         *config.ValueSource `json:"carLoadUnits,omitempty"`
	NextStoppingFloor    *config.ValueSource `json:"nextStoppingFloor,omitempty"`
	PassengerAlarm       *config.ValueSource `json:"passengerAlarm,omitempty"`
	CarDriveStatus       *config.ValueSource `json:"carDriveStatus,omitempty"`
	FaultSignals         *config.ValueSource `json:"faultSignals,omitempty"`
}

func readTransportConfig(raw []byte) (cfg transportCfg, err error) {
	err = json.Unmarshal(raw, &cfg)
	return
}

type transport struct {
	gen.UnimplementedTransportApiServer
	gen.UnimplementedTransportInfoServer

	client   *gobacnet.Client
	known    known.Context
	statuses *statuspb.Map
	logger   *zap.Logger

	model *transportpb.Model
	*transportpb.ModelServer
	config   transportCfg
	pollTask *task.Intermittent

	units atomic.Value
}

func newTransport(client *gobacnet.Client, known known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*transport, error) {
	cfg, err := readTransportConfig(config.Raw)
	if err != nil {
		return nil, err
	}

	model := transportpb.NewModel()

	t := &transport{
		client:      client,
		known:       known,
		statuses:    statuses,
		logger:      logger,
		model:       model,
		ModelServer: transportpb.NewModelServer(model),
		config:      cfg,
	}

	t.pollTask = task.NewIntermittent(t.startPoll)

	initTraitStatus(statuses, cfg.Name, "Transport")

	return t, nil
}

func (t *transport) AnnounceSelf(a node.Announcer) node.Undo {
	return a.Announce(t.config.Name, node.HasTrait(transportpb.TraitName, node.WithClients(gen.WrapTransportApi(t), gen.WrapTransportInfo(t))))
}

func (t *transport) startPoll(init context.Context) (stop task.StopFn, err error) {
	return startPoll(init, "transport", t.config.PollTimeoutDuration(), t.config.PollTimeoutDuration(), t.logger, func(ctx context.Context) error {
		_, err := t.pollPeer(ctx)
		return err
	})
}

func (t *transport) pollPeer(ctx context.Context) (*gen.Transport, error) {
	data := &gen.Transport{}

	var resProcessors []func(response any, data *gen.Transport) error
	var readValues []config.ValueSource
	var requestNames []string

	if t.config.CarLoadUnits != nil {
		requestNames = append(requestNames, "loadUnits")
		readValues = append(readValues, *t.config.CarLoadUnits)
		resProcessors = append(resProcessors, t.processCarLoadUnits)
	}

	if t.config.AssignedLandingCalls != nil {
		requestNames = append(requestNames, "assignedLandingCalls")
		readValues = append(readValues, *t.config.AssignedLandingCalls)
		resProcessors = append(resProcessors, t.processAssignedLandingCalls)
	}

	if t.config.MakingCarCall != nil {
		requestNames = append(requestNames, "makingCarCall")
		readValues = append(readValues, *t.config.MakingCarCall)
		resProcessors = append(resProcessors, t.processMakingCarCall)
	}

	if t.config.CarPosition != nil {
		requestNames = append(requestNames, "carPosition")
		readValues = append(readValues, *t.config.CarPosition)
		resProcessors = append(resProcessors, t.processCarPosition)
	}

	if t.config.CarAssignedDirection != nil {
		requestNames = append(requestNames, "carAssignedDirection")
		readValues = append(readValues, *t.config.CarAssignedDirection)
		resProcessors = append(resProcessors, t.processCarAssignedDirection)
	}

	if t.config.CarDoorStatus != nil {
		requestNames = append(requestNames, "carDoorStatus")
		readValues = append(readValues, *t.config.CarDoorStatus)
		resProcessors = append(resProcessors, t.processCarDoorStatus)
	}

	if t.config.CarMode != nil {
		requestNames = append(requestNames, "carMode")
		readValues = append(readValues, *t.config.CarMode)
		resProcessors = append(resProcessors, t.processCarMode)
	}

	if t.config.CarLoad != nil {
		requestNames = append(requestNames, "carLoad")
		readValues = append(readValues, *t.config.CarLoad)
		resProcessors = append(resProcessors, t.processCarLoad)
	}

	if t.config.NextStoppingFloor != nil {
		requestNames = append(requestNames, "nextStoppingFloor")
		readValues = append(readValues, *t.config.NextStoppingFloor)
		resProcessors = append(resProcessors, t.processNextStoppingFloor)
	}

	if t.config.FaultSignals != nil {
		requestNames = append(requestNames, "faultSignals")
		readValues = append(readValues, *t.config.FaultSignals)
		resProcessors = append(resProcessors, t.processFaultSignals)
	}

	if t.config.PassengerAlarm != nil {
		requestNames = append(requestNames, "passengerAlarm")
		readValues = append(readValues, *t.config.PassengerAlarm)
		resProcessors = append(resProcessors, t.processPassengerAlarm)
	}

	responses := comm.ReadProperties(ctx, t.client, t.known, readValues...)
	var errs []error
	for i, response := range responses {
		err := resProcessors[i](response, data)
		if err != nil {
			errs = append(errs, err)
		}
	}

	status.UpdatePollErrorStatus(t.statuses, t.config.Name, "Transport", requestNames, errs)
	if len(errs) > 0 {
		return nil, multierr.Combine(errs...)
	}

	return t.model.UpdateTransport(data)
}

func (t *transport) processCarLoadUnits(response any, _ *gen.Transport) error {
	value, ok := response.(comm.EngineeringUnits)
	if !ok {
		return comm.ErrReadProperty{Prop: "loadUnits", Cause: fmt.Errorf("converting to EngineeringUnits")}
	}

	t.units.Store(value.String())

	return nil

}

func (t *transport) processAssignedLandingCalls(response any, data *gen.Transport) error {
	value, ok := response.([]comm.LandingCall)
	if !ok {
		return comm.ErrReadProperty{Prop: "assignedLandingCalls", Cause: fmt.Errorf("converting to AssignedLandingCalls")}
	}

	for _, v := range value {
		data.NextDestinations = append(data.NextDestinations, &gen.Transport_Location{
			Id:    fmt.Sprintf("%d", v.Floor),
			Title: fmt.Sprintf("Floor %d", v.Floor),
		})
	}
	return nil
}

func (t *transport) processMakingCarCall(response any, data *gen.Transport) error {
	if response == nil {
		return nil
	}

	switch response.(type) {
	case []interface{}:
		return nil
	}

	value, ok := response.([]uint8)

	if !ok {
		return comm.ErrReadProperty{Prop: "makingCarCall", Cause: fmt.Errorf("converting to []int8")}
	}

	for _, v := range value {
		data.NextDestinations = append(data.NextDestinations, &gen.Transport_Location{Id: fmt.Sprintf("%d", v), Title: fmt.Sprintf("Floor %d", v)})
	}
	return nil
}

func (t *transport) processCarPosition(response any, data *gen.Transport) error {
	value, err := comm.IntValue(response)

	if err != nil {
		return comm.ErrReadProperty{Prop: "carPosition", Cause: err}
	}

	data.ActualPosition = &gen.Transport_Location{Id: fmt.Sprintf("%d", value), Title: fmt.Sprintf("Floor %d", value)}
	return nil
}

func (t *transport) processCarAssignedDirection(response any, data *gen.Transport) error {
	value, ok := response.(comm.LiftCarDirection)
	if !ok {
		return comm.ErrReadProperty{Prop: "carAssignedDirection", Cause: fmt.Errorf("converting to CarAssignedDirection")}
	}

	switch value {
	case comm.DirectionUpAndDown:
		fallthrough
	case comm.DirectionUp:
		data.MovingDirection = gen.Transport_UP
	case comm.DirectionDown:
		data.MovingDirection = gen.Transport_DOWN
	default:
		data.MovingDirection = gen.Transport_NO_DIRECTION
	}

	return nil
}

func (t *transport) processCarDoorStatus(response any, data *gen.Transport) error {
	value, err := comm.IntValue(response)
	if err != nil {
		return comm.ErrReadProperty{Prop: "carDoorStatus", Cause: err}
	}

	switch comm.DoorStatus(value) {
	case comm.DoorClosed:
		data.Doors = append(data.Doors, &gen.Transport_Door{Status: gen.Transport_Door_CLOSED})
	case comm.DoorClosing:
		data.Doors = append(data.Doors, &gen.Transport_Door{Status: gen.Transport_Door_CLOSING})
	case comm.DoorOpened:
		data.Doors = append(data.Doors, &gen.Transport_Door{Status: gen.Transport_Door_OPEN})
	case comm.DoorOpening:
		data.Doors = append(data.Doors, &gen.Transport_Door{Status: gen.Transport_Door_OPENING})
	case comm.DoorSafetyLocked:
		data.Doors = append(data.Doors, &gen.Transport_Door{Status: gen.Transport_Door_SAFETY_LOCKED})
	case comm.DoorLimitedOpened:
		data.Doors = append(data.Doors, &gen.Transport_Door{Status: gen.Transport_Door_LIMITED_OPENED})
	default:
		data.Doors = append(data.Doors, &gen.Transport_Door{Status: gen.Transport_Door_DOOR_STATUS_UNSPECIFIED})
	}

	return nil
}

func (t *transport) processCarMode(response any, data *gen.Transport) error {
	value, err := comm.IntValue(response)

	if err != nil {
		return comm.ErrReadProperty{Prop: "carMode", Cause: err}
	}

	switch comm.LiftCarMode(value) {
	case comm.LiftCarModeNormal:
		data.OperatingMode = gen.Transport_NORMAL
	case comm.LiftCarModeVIP:
		data.OperatingMode = gen.Transport_VIP_CONTROL
	case comm.LiftCarModeInspection:
		data.OperatingMode = gen.Transport_SERVICE_CONTROL
	case comm.LiftCarModeFirefighterControl:
		data.OperatingMode = gen.Transport_FIRE_OPERATION
	case comm.LiftCarModeEmergencyPower:
		data.OperatingMode = gen.Transport_EMERGENCY_POWER
	case comm.LiftCarModeEarthquakeOp:
		data.OperatingMode = gen.Transport_EARTHQUAKE_OPERATION
	case comm.LiftCarModeOccupantEvac:
		data.OperatingMode = gen.Transport_OCCUPANT_EVACUATION
	case comm.LiftCarModeHoming:
		data.OperatingMode = gen.Transport_HOMING
	case comm.LiftCarModeParking:
		data.OperatingMode = gen.Transport_PARKING
	case comm.LiftCarModeAttendantControl:
		data.OperatingMode = gen.Transport_ATTENDANT_CONTROL
	case comm.LiftCarModeCabinetRecall:
		data.OperatingMode = gen.Transport_CABINET_RECALL
	case comm.LiftCarModeOutOfService:
		data.OperatingMode = gen.Transport_OUT_OF_SERVICE
	default:
		data.OperatingMode = gen.Transport_OPERATING_MODE_UNSPECIFIED
	}

	return nil
}

func (t *transport) processCarLoad(response any, data *gen.Transport) error {
	value, err := comm.Float32Value(response)
	if err != nil {
		return comm.ErrReadProperty{Prop: "carLoad", Cause: err}
	}

	data.Load = &value
	return nil
}

func (t *transport) processNextStoppingFloor(response any, data *gen.Transport) error {
	value, err := comm.IntValue(response)
	if err != nil {
		return comm.ErrReadProperty{Prop: "nextStoppingFloor", Cause: err}
	}

	data.Payloads = []*gen.Transport_Payload{
		{
			PayloadId: fmt.Sprintf("%d", value),
			IntendedJourney: &gen.Transport_Journey{
				Destinations: []*gen.Transport_Location{
					{
						Title: fmt.Sprintf("Floor %d", value),
					},
				},
			},
		},
	}
	return nil
}

func (t *transport) processFaultSignals(response any, data *gen.Transport) error {
	res, ok := response.([]interface{})

	if !ok {
		return comm.ErrReadProperty{Prop: "faultSignals", Cause: fmt.Errorf("converting to []interface{}")}
	}

	var value []comm.LiftFault
	for _, r := range res {
		v, err := comm.IntValue(r)
		if err != nil {
			return comm.ErrReadProperty{Prop: "faultSignals", Cause: fmt.Errorf("converting to int: %w", err)}
		}
		value = append(value, comm.LiftFault(v))
	}

	for _, v := range value {
		switch v {
		case comm.LiftFaultControllerFault:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_CONTROLLER_FAULT})
		case comm.LiftFaultDriveAndMotorFault:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_DRIVE_AND_MOTOR_FAULT})
		case comm.LiftFaultGovernorAndSafetyGearFault:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_GOVERNOR_AND_SAFETY_GEAR_FAULT})
		case comm.LiftFaultLiftShaftDeviceFault:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_LIFT_SHAFT_DEVICE_FAULT})
		case comm.LiftFaultPowerSupplyFault:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_POWER_SUPPLY_FAULT})
		case comm.LiftFaultSafetyInterlockFault:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_SAFETY_DEVICE_FAULT})
		case comm.LiftFaultDoorClosingFault:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_DOOR_NOT_CLOSING})
		case comm.LiftFaultDoorOpeningFault:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_DOOR_NOT_OPENING})
		case comm.LiftFaultCarStoppedOutsideLanding:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_CAR_STOPPED_OUTSIDE_LANDING_ZONE})
		case comm.LiftFaultCallButtonStuck:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_CALL_BUTTON_STUCK})
		case comm.LiftFaultStartFailure:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_FAIL_TO_START})
		case comm.LiftFaultControllerSupplyFault:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_CONTROLLER_SUPPLY_FAULT})
		case comm.LiftFaultSelfTestFailure:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_SELF_TEST_FAILURE})
		case comm.LiftFaultRuntimeLimitExceeded:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_RUNTIME_LIMIT_EXCEEDED})
		case comm.LiftFaultPositionLost:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_POSITION_LOST})
		case comm.LiftFaultDriveTempExceeded:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_DRIVE_AND_MOTOR_FAULT})
		case comm.LiftFaultLoadMeasurementFault:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_LOAD_MEASUREMENT_FAULT})
		default:
			data.Faults = append(data.Faults, &gen.Transport_Fault{FaultType: gen.Transport_Fault_FAULT_TYPE_UNSPECIFIED})
		}
	}

	return nil
}

func (t *transport) processPassengerAlarm(response any, data *gen.Transport) error {
	value, err := comm.BoolValue(response)
	if err != nil {
		return comm.ErrReadProperty{Prop: "passengerAlarm", Cause: err}
	}

	if value {
		data.PassengerAlarm = &gen.Transport_Alarm{State: gen.Transport_Alarm_ACTIVATED, Time: timestamppb.Now()}
	} else {
		data.PassengerAlarm = &gen.Transport_Alarm{State: gen.Transport_Alarm_UNACTIVATED, Time: timestamppb.Now()}
	}

	return nil
}

func (t *transport) DescribeTransport(ctx context.Context, _ *gen.DescribeTransportRequest) (*gen.TransportSupport, error) {
	err := t.pollTask.Attach(ctx)

	if err != nil {
		return nil, grpcStatus.New(codes.Internal, err.Error()).Err()
	}

	unit, _ := t.units.Load().(string)
	if unit == "" {
		unit = "unspecified"
	}
	return &gen.TransportSupport{
		LoadUnit: unit,
	}, nil
}

func (t *transport) GetTransport(ctx context.Context, request *gen.GetTransportRequest) (*gen.Transport, error) {
	err := t.pollTask.Attach(ctx)
	if err != nil {
		return nil, grpcStatus.New(codes.Internal, err.Error()).Err()
	}
	return t.ModelServer.GetTransport(ctx, request)
}

func (t *transport) PullTransport(request *gen.PullTransportRequest, server gen.TransportApi_PullTransportServer) error {
	err := t.pollTask.Attach(server.Context())
	if err != nil {
		return grpcStatus.New(codes.Internal, err.Error()).Err()
	}
	return t.ModelServer.PullTransport(request, server)
}
