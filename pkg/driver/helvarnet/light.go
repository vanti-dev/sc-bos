package helvarnet

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/auto/udmi"
	"github.com/vanti-dev/sc-bos/pkg/driver/helvarnet/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/minibus"
)

// Light represents a single light device within the HelvarNet system.
type Light struct {
	gen.UnimplementedStatusApiServer
	traits.UnimplementedLightApiServer
	gen.UnimplementedEmergencyLightApiServer
	gen.UnimplementedUdmiServiceServer

	brightness      *resource.Value // *traits.Brightness
	client          *tcpClient
	conf            *config.Device
	logger          *zap.Logger
	helvarnetStatus uint32          // The status flags field from the device, unique to Helvarnet protocol. See config.DeviceStatuses
	status          *resource.Value // *gen.StatusLog
	udmiBus         minibus.Bus[*gen.PullExportMessagesResponse]

	testResultSet *resource.Value // *gen.TestResultSet
}

func newLight(client *tcpClient, l *zap.Logger, conf *config.Device) *Light {
	return &Light{
		brightness:    resource.NewValue(resource.WithInitialValue(&traits.Brightness{}), resource.WithNoDuplicates()),
		client:        client,
		conf:          conf,
		logger:        l,
		status:        resource.NewValue(resource.WithInitialValue(&gen.StatusLog{}), resource.WithNoDuplicates()),
		testResultSet: resource.NewValue(resource.WithInitialValue(&gen.TestResultSet{}), resource.WithNoDuplicates()),
	}
}

// setScene sets the lighting scene for the device
func (l *Light) setScene(block string, scene string, constant string) error {
	command := recallDeviceScene(l.conf.Address, block, scene, constant)

	_, err := l.client.sendAndReceive(command, "")
	if err != nil {
		return err
	}
	return nil
}

// setLevel sets the light level for this device
func (l *Light) setLevel(level int) error {
	command := changeDeviceLevel(l.conf.Address, level)

	_, err := l.client.sendAndReceive(command, "")
	if err != nil {
		return err
	}
	return nil
}

// refreshBrightness queries the device's load and updates the brightness value
func (l *Light) refreshBrightness() error {
	command := queryLoadLevel(l.conf.Address)
	want := "?" + command[1:len(command)-1]

	r, err := l.client.sendAndReceive(command, want)
	if err != nil {
		return err
	}

	split := strings.Split(r, "=")
	if len(split) < 2 {
		return fmt.Errorf("invalid response in refreshBrightness: %s", r)
	}

	s := strings.TrimSuffix(split[1], "#")
	brightness, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	_, _ = l.brightness.Set(&traits.Brightness{
		LevelPercent: float32(brightness),
	})
	return nil
}

// refreshDeviceStatus queries the device and updates the status value
func (l *Light) refreshDeviceStatus() error {
	command := queryDeviceState(l.conf.Address)
	want := "?" + command[1:len(command)-1]

	r, err := l.client.sendAndReceive(command, want)
	if err != nil {
		return err
	}

	split := strings.Split(r, "=")
	if len(split) < 2 {
		return fmt.Errorf("invalid response in refreshDeviceStatus: %s", r)
	}

	s := strings.TrimSuffix(split[1], "#")
	statusInt, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	l.helvarnetStatus = uint32(statusInt)
	log := &gen.StatusLog{
		RecordTime: timestamppb.Now(),
	}
	for _, ds := range config.DeviceStatuses {
		if (ds.FlagValue & l.helvarnetStatus) > 0 {
			log.Problems = append(log.Problems, &gen.StatusLog_Problem{
				Level:       ds.Level,
				Name:        ds.State,
				Description: ds.Description,
			})
		}
	}
	_, _ = l.status.Set(log)
	return nil
}

// UpdateBrightness update the brightness level or preset (scene) of the device
// if the request has a present included, this takes precedence and the level percent is ignored
func (l *Light) UpdateBrightness(_ context.Context, req *traits.UpdateBrightnessRequest) (*traits.Brightness, error) {
	if req.Brightness == nil {
		return nil, status.Error(codes.InvalidArgument, "no brightness in request")
	}

	// I am not sure how the scene recall works at the device level, there is a command to set it so we support it
	// but there is no command to query the scene names for devices like here is for groups
	if req.Brightness.Preset != nil {
		// helvarnet scenes are in 8 blocks of 16 scenes, preset name is described in info as <block>:<scene>
		sceneSplit := strings.Split(req.Brightness.Preset.Name, ":")
		if len(sceneSplit) < 2 {
			return nil, status.Error(codes.InvalidArgument, "invalid scene format, must be <block>:<scene>")
		}
		block := sceneSplit[0]
		scene := sceneSplit[1]
		constant := "0"

		if len(sceneSplit) == 3 {
			constant = sceneSplit[2]
		}
		err := l.setScene(block, scene, constant)
		if err != nil {
			return nil, status.Error(codes.DeadlineExceeded, "failed to set scene")
		}
		_, _ = l.brightness.Set(&traits.Brightness{
			Preset: &traits.LightPreset{
				Name: req.Brightness.Preset.Name,
			},
		})
	} else {
		level := req.Brightness.LevelPercent
		err := l.setLevel(int(level))
		if err != nil {
			return nil, status.Error(codes.DeadlineExceeded, "failed to set scene")
		}
		_, _ = l.brightness.Set(&traits.Brightness{
			LevelPercent: level,
		})
	}

	return nil, nil
}

func (l *Light) GetBrightness(_ context.Context, _ *traits.GetBrightnessRequest) (*traits.Brightness, error) {
	err := l.refreshBrightness()
	if err != nil {
		return nil, status.Error(codes.DeadlineExceeded, "failed to get brightness")
	}
	value := l.brightness.Get()
	brightness := value.(*traits.Brightness)
	return brightness, nil
}

func (l *Light) PullBrightness(_ *traits.PullBrightnessRequest, server traits.LightApi_PullBrightnessServer) error {
	for value := range l.brightness.Pull(server.Context()) {
		brightness := value.Value.(*traits.Brightness)
		err := server.Send(&traits.PullBrightnessResponse{Changes: []*traits.PullBrightnessResponse_Change{
			{
				Name:       l.conf.Name,
				ChangeTime: timestamppb.New(value.ChangeTime),
				Brightness: brightness,
			},
		}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Light) GetCurrentStatus(context.Context, *gen.GetCurrentStatusRequest) (*gen.StatusLog, error) {
	value := l.status.Get()
	s := value.(*gen.StatusLog)
	return s, nil
}

func (l *Light) PullCurrentStatus(_ *gen.PullCurrentStatusRequest, server gen.StatusApi_PullCurrentStatusServer) error {
	for value := range l.status.Pull(server.Context()) {
		statusLog := value.Value.(*gen.StatusLog)
		err := server.Send(&gen.PullCurrentStatusResponse{Changes: []*gen.PullCurrentStatusResponse_Change{
			{
				Name:          l.conf.Name,
				ChangeTime:    timestamppb.New(value.ChangeTime),
				CurrentStatus: statusLog,
			},
		}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Light) udmiPointsetFromData() (*gen.MqttMessage, error) {
	points := make(udmi.PointsEvent)
	brightness := l.brightness.Get().(*traits.Brightness)
	points["BrightnessLvl%"] = udmi.PointValue{PresentValue: brightness.LevelPercent}

	if brightness.Preset != nil && brightness.Preset.Title != "" {
		points["Preset"] = udmi.PointValue{PresentValue: brightness.Preset.Title}
	}

	statuses := config.GetStatusListFromFlag(l.helvarnetStatus)
	points["Status"] = udmi.PointValue{PresentValue: strings.Join(statuses, ", ")}

	b, err := json.Marshal(points)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal udmi points: %w", err)
	}
	return &gen.MqttMessage{
		Topic:   l.conf.TopicPrefix + "/event/pointset/points",
		Payload: string(b),
	}, nil
}

func (l *Light) sendUdmiMessage(ctx context.Context) {
	m, err := l.udmiPointsetFromData()
	if err != nil {
		l.logger.Error("failed to create udmi pointset message", zap.Error(err))
		return
	}

	l.udmiBus.Send(ctx, &gen.PullExportMessagesResponse{
		Name:    l.conf.Name,
		Message: m,
	})
}

func (l *Light) refreshData(ctx context.Context) {
	err := l.refreshDeviceStatus()
	if err != nil {
		l.logger.Error("failed to refresh device status, will try again on next run...", zap.Error(err))
	}
	err = l.refreshBrightness()
	if err != nil {
		l.logger.Error("failed to refresh brightness, will try again on next run...", zap.Error(err))
	}

	l.sendUdmiMessage(ctx)
}

// runHealthCheck runs queries on a schedule to check the health of the device.
func (l *Light) runHealthCheck(ctx context.Context, t time.Duration) error {
	ticker := time.NewTicker(t)
	defer ticker.Stop()
	l.refreshData(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			l.refreshData(ctx)
		}
	}
}

// runFunctionTest requests a function test from the device. Does not expect a response.
// To get the result of the function test, you need to call queryEmergencyFunctionTestState
func (l *Light) runFunctionTest() error {
	command := deviceEmergencyFunctionTest(l.conf.Address)

	_, err := l.client.sendAndReceive(command, "")
	if err != nil {
		return err
	}
	return nil
}

// runDurationTest requests a duration test from the device. Does not expect a response.
// To get the result of the duration test, you need to call queryEmergencyDurationTestState
func (l *Light) runDurationTest() error {
	command := deviceEmergencyDurationTest(l.conf.Address)

	_, err := l.client.sendAndReceive(command, "")
	if err != nil {
		return err
	}
	return nil
}

// stopTest stops any running emergency test on the device.
func (l *Light) stopTest() error {
	command := deviceStopEmergencyTests(l.conf.Address)

	_, err := l.client.sendAndReceive(command, "")
	if err != nil {
		return err
	}
	return nil
}

// EmergencyState Values
//
// - Pass 0
//
// - Lamp Failure 1
//
// - Battery Failure 2
//
// - Faulty 4
//
// - Failure 8
//
// - Test Pending 16
//
// - Unknown 32
type EmergencyState int

const (
	Pass           EmergencyState = 0
	LampFailure    EmergencyState = 1
	BatteryFailure EmergencyState = 2
	Faulty         EmergencyState = 4
	Failure        EmergencyState = 8
	TestPending    EmergencyState = 16
	Unknown        EmergencyState = 32
)

func hasTestCompleted(state *EmergencyState) bool {
	if state == nil {
		return false
	}
	switch *state {
	case Pass, LampFailure, BatteryFailure, Faulty, Failure:
		return true
	case TestPending, Unknown:
		return false
	default:
		return false
	}
}

func parseGetResultResponse(r string) (*EmergencyState, error) {
	// example response ?V:1,C:171,@1.1.2.15=16#
	split := strings.Split(r, "=")
	if len(split) < 2 {
		return nil, fmt.Errorf("invalid response in getFunctionTestResult: %s", r)
	}
	state := strings.TrimSuffix(split[1], "#")
	stateInt, err := strconv.Atoi(state)
	if err != nil {
		return nil, fmt.Errorf("failed to parse function test state: %w", err)
	}

	switch EmergencyState(stateInt) {
	case Pass, LampFailure, BatteryFailure, Faulty, Failure, TestPending, Unknown:
		// Do nothing, we have a valid state
	default:
		return nil, fmt.Errorf("unknown emergency state: %d", stateInt)
	}

	e := EmergencyState(stateInt)
	return &e, nil
}

// getFunctionTestResult queries the device for the result of the last function test.
// The result is a valid EmergencyState value as defined by the protocol, else an error is returned.
func (l *Light) getFunctionTestResult() (*EmergencyState, error) {

	command := queryEmergencyFunctionTestState(l.conf.Address)
	want := "?" + command[1:len(command)-1]

	r, err := l.client.sendAndReceive(command, want)
	if err != nil {
		return nil, err
	}

	return parseGetResultResponse(r)
}

// getDurationTestResult queries the device for the result of the last duration test.
// The result is a valid EmergencyState value as defined by the protocol, else an error is returned.
func (l *Light) getDurationTestResult() (*EmergencyState, error) {

	command := queryEmergencyDurationTestState(l.conf.Address)
	want := "?" + command[1:len(command)-1]

	r, err := l.client.sendAndReceive(command, want)
	if err != nil {
		return nil, err
	}

	return parseGetResultResponse(r)
}

// parseGetCompletionTimeResponse parses the response from the device for the completion time of a function/duration test.
func parseGetCompletionTimeResponse(r string) (*time.Time, error) {
	// example response ?V:1,C:170,@10.106.4.40=1754495355#
	// the time is seconds in the Linux epoch
	split := strings.Split(r, "=")
	if len(split) < 2 {
		return nil, fmt.Errorf("invalid response in getFunctionTestCompletionTime: %s", r)
	}
	timeStr := strings.TrimSuffix(split[1], "#")
	epochSeconds, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse function test completion time: %w", err)
	}
	if epochSeconds == 0 {
		// this means the test has never been run
		return nil, fmt.Errorf("function test completion time is zero, test has never been run")
	}
	t := time.Unix(epochSeconds, 0)
	if t.IsZero() {
		return nil, fmt.Errorf("function test completion time is zero")
	}
	return &t, nil
}

// getFunctionTestCompletionTime queries the device for the finish time of the last function test.
func (l *Light) getFunctionTestCompletionTime() (*time.Time, error) {

	command := queryEmergencyFunctionTestTime(l.conf.Address)
	want := "?" + command[1:len(command)-1]

	r, err := l.client.sendAndReceive(command, want)
	if err != nil {
		return nil, err
	}

	return parseGetCompletionTimeResponse(r)
}

// getDurationTestCompletionTime queries the device for the finish time of the last duration test.
func (l *Light) getDurationTestCompletionTime() (*time.Time, error) {

	command := queryEmergencyDurationTestTime(l.conf.Address)
	want := "?" + command[1:len(command)-1]

	r, err := l.client.sendAndReceive(command, want)
	if err != nil {
		return nil, err
	}

	return parseGetCompletionTimeResponse(r)
}

func (l *Light) StartFunctionTest(context.Context, *gen.StartEmergencyTestRequest) (*gen.StartEmergencyTestResponse, error) {
	l.logger.Info("Starting function test for light", zap.String("name", l.conf.Name))
	err := l.runFunctionTest()
	if err != nil {
		l.logger.Error("Failed to start function test", zap.String("name", l.conf.Name), zap.Error(err))
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to start function test"))
	}
	_, _ = l.testResultSet.Set(&gen.TestResultSet{
		FunctionTest: &gen.EmergencyTestResult{
			StartTime: timestamppb.Now(),
		},
	}, resource.WithUpdateMask(&fieldmaskpb.FieldMask{Paths: []string{"function_test.start_time"}}))
	return &gen.StartEmergencyTestResponse{
		StartTime: timestamppb.Now(),
	}, nil
}

func (l *Light) StartDurationTest(context.Context, *gen.StartEmergencyTestRequest) (*gen.StartEmergencyTestResponse, error) {
	l.logger.Info("Starting duration test for light", zap.String("name", l.conf.Name))
	err := l.runDurationTest()
	if err != nil {
		l.logger.Error("Failed to start duration test", zap.String("name", l.conf.Name), zap.Error(err))
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to start duration test"))
	}
	_, _ = l.testResultSet.Set(&gen.TestResultSet{
		DurationTest: &gen.EmergencyTestResult{
			StartTime: timestamppb.Now(),
		},
	}, resource.WithUpdateMask(&fieldmaskpb.FieldMask{Paths: []string{"duration_test.start_time"}}))
	var duration *durationpb.Duration
	if l.conf.DurationTestLength != nil {
		duration = durationpb.New(l.conf.DurationTestLength.Duration)
	}
	return &gen.StartEmergencyTestResponse{
		StartTime: timestamppb.Now(),
		Duration:  duration,
	}, nil
}

func (l *Light) StopEmergencyTest(context.Context, *gen.StopEmergencyTestsRequest) (*gen.StopEmergencyTestsResponse, error) {
	l.logger.Info("Stopping test for light", zap.String("name", l.conf.Name))
	err := l.stopTest()
	if err != nil {
		l.logger.Error("Failed to stop test", zap.String("name", l.conf.Name), zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to stop test")
	}
	// when we stop tests clear any current data
	_, _ = l.testResultSet.Set(&gen.TestResultSet{})
	return &gen.StopEmergencyTestsResponse{}, nil
}

func (l *Light) GetTestResultSet(_ context.Context, req *gen.GetTestResultSetRequest) (*gen.TestResultSet, error) {

	result := &gen.TestResultSet{}

	if req.ReadMask == nil || slices.Contains(req.ReadMask.Paths, "function_test") {
		result.FunctionTest = l.testResultSet.Get().(*gen.TestResultSet).FunctionTest
		if req.QueryDevice {
			fRes, err := getTestResult(l.getFunctionTestResult, l.getFunctionTestCompletionTime)
			if err != nil {
				l.logger.Error("Failed to get function test result", zap.String("name", l.conf.Name), zap.Error(err))
				return nil, err
			}
			result.FunctionTest = fRes
			// update the stored test result set with the new result
			_, _ = l.testResultSet.Set(result, resource.WithUpdateMask(&fieldmaskpb.FieldMask{
				Paths: []string{"function_test"},
			}))
		}
	}
	if req.ReadMask == nil || slices.Contains(req.ReadMask.Paths, "duration_test") {
		result.DurationTest = l.testResultSet.Get().(*gen.TestResultSet).DurationTest
		if req.QueryDevice {
			dRes, err := getTestResult(l.getDurationTestResult, l.getDurationTestCompletionTime)
			if err != nil {
				l.logger.Error("Failed to get duration test result", zap.String("name", l.conf.Name), zap.Error(err))
				return nil, err
			}
			result.DurationTest = dRes

			r, err := l.testResultSet.Set(result, resource.WithUpdateMask(&fieldmaskpb.FieldMask{
				Paths: []string{"duration_test"},
			}))
			if err == nil {
				result = r.(*gen.TestResultSet)
			}
		}
	}
	return result, nil
}

func (l *Light) PullTestResultSets(request *gen.PullTestResultRequest, server grpc.ServerStreamingServer[gen.PullTestResultsResponse]) error {
	for value := range l.testResultSet.Pull(server.Context(), resource.WithReadMask(request.GetReadMask()), resource.WithUpdatesOnly(request.UpdatesOnly)) {
		resultSet := value.Value.(*gen.TestResultSet)
		err := server.Send(&gen.PullTestResultsResponse{Changes: []*gen.PullTestResultsResponse_Change{
			{
				Name:       l.conf.Name,
				ChangeTime: timestamppb.New(value.ChangeTime),
				TestResult: resultSet,
			},
		}})
		if err != nil {
			return err
		}
	}
	return nil
}

func getTestResult(getResult func() (*EmergencyState, error), getTime func() (*time.Time, error)) (*gen.EmergencyTestResult, error) {

	eState, err := getResult()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to getResult test result")
	}
	result := &gen.EmergencyTestResult{}
	result.Result = helvarResultToTrait(eState)
	if hasTestCompleted(eState) {
		t, _ := getTime()
		if t != nil {
			result.EndTime = timestamppb.New(*t)
		}
	}
	return result, nil
}

func helvarResultToTrait(e *EmergencyState) gen.EmergencyTestResult_Result {
	if e == nil {
		return gen.EmergencyTestResult_TEST_RESULT_UNSPECIFIED
	}
	switch *e {
	case Pass:
		return gen.EmergencyTestResult_TEST_PASSED
	case LampFailure:
		return gen.EmergencyTestResult_LAMP_FAILURE
	case BatteryFailure:
		return gen.EmergencyTestResult_BATTERY_FAILURE
	case Faulty:
		return gen.EmergencyTestResult_LIGHT_FAULTY
	case Failure:
		return gen.EmergencyTestResult_TEST_FAILED
	case TestPending:
		return gen.EmergencyTestResult_TEST_RESULT_PENDING
	case Unknown:
		return gen.EmergencyTestResult_TEST_RESULT_UNSPECIFIED
	}
	return gen.EmergencyTestResult_TEST_RESULT_UNSPECIFIED
}

func (l *Light) GetExportMessage(context.Context, *gen.GetExportMessageRequest) (*gen.MqttMessage, error) {
	m, err := l.udmiPointsetFromData()
	if err != nil {
		l.logger.Error("failed to create udmi pointset message", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create udmi pointset message")
	}
	return m, nil
}

func (l *Light) PullExportMessages(_ *gen.PullExportMessagesRequest, server gen.UdmiService_PullExportMessagesServer) error {
	for msg := range l.udmiBus.Listen(server.Context()) {
		err := server.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Light) PullControlTopics(*gen.PullControlTopicsRequest, grpc.ServerStreamingServer[gen.PullControlTopicsResponse]) error {
	return status.Error(codes.Unimplemented, "PullControlTopics is not implemented for Light")
}

func (l *Light) OnMessage(context.Context, *gen.OnMessageRequest) (*gen.OnMessageResponse, error) {
	return nil, status.Error(codes.Unimplemented, "OnMessage is not implemented for Light")
}
