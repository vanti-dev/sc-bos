package helvarnet

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/driver/helvarnet/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

// Light represents a single light device within the HelvarNet system.
type Light struct {
	gen.UnimplementedStatusApiServer
	traits.UnimplementedLightApiServer
	// Helvarnet is not really DALI, but it is really just a wrapper around DALI devices
	// the DALI API has the stuff in for emergency lighting tests, so just use it for now.
	gen.UnimplementedDaliApiServer

	brightness *resource.Value // *traits.Brightness
	client     *tcpClient
	conf       *config.Device
	logger     *zap.Logger
	status     *resource.Value // *gen.StatusLog
}

func newLight(client *tcpClient, l *zap.Logger, conf *config.Device) *Light {
	return &Light{
		brightness: resource.NewValue(resource.WithInitialValue(&traits.Brightness{}), resource.WithNoDuplicates()),
		client:     client,
		conf:       conf,
		logger:     l,
		status:     resource.NewValue(resource.WithInitialValue(&gen.StatusLog{}), resource.WithNoDuplicates()),
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

	log := &gen.StatusLog{
		RecordTime: timestamppb.Now(),
	}
	for _, ds := range config.DeviceStatuses {
		if (ds.FlagValue & statusInt) > 0 {
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
	err := l.refreshDeviceStatus()
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

// runHealthCheck runs queries on a schedule to check the health of the device.
func (l *Light) runHealthCheck(ctx context.Context, t time.Duration) error {
	ticker := time.NewTicker(t)
	defer ticker.Stop()
	err := l.refreshDeviceStatus()
	if err != nil {
		l.logger.Error("failed to refresh device status, will try again on next run...", zap.Error(err))
	}
	err = l.refreshBrightness()
	if err != nil {
		l.logger.Error("failed to refresh brightness, will try again on next run...", zap.Error(err))
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := l.refreshDeviceStatus()
			if err != nil {
				l.logger.Error("failed to refresh device status, will try again on next run...", zap.Error(err))
			}
			err = l.refreshBrightness()
			if err != nil {
				l.logger.Error("failed to refresh brightness, will try again on next run...", zap.Error(err))
			}
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
	l.logger.Debug("getFunctionTestResult", zap.String("response", r), zap.String("command", command))
	if err != nil {
		return nil, err
	}

	// todo implement GetEmergencyStatus, if the response is bad we want it to show up
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

// StartTest Attempt to start a function or duration test.
func (l *Light) StartTest(_ context.Context, req *gen.StartTestRequest) (*gen.StartTestResponse, error) {

	switch req.Test {
	case gen.EmergencyStatus_FUNCTION_TEST:
		l.logger.Info("Starting function test for light", zap.String("name", l.conf.Name))
		err := l.runFunctionTest()
		if err != nil {
			l.logger.Error("Failed to start function test", zap.String("name", l.conf.Name), zap.Error(err))
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to start function test"))
		}
	case gen.EmergencyStatus_DURATION_TEST:
		l.logger.Info("Starting duration test for light", zap.String("name", l.conf.Name))
		err := l.runDurationTest()
		if err != nil {
			l.logger.Error("Failed to start duration test", zap.String("name", l.conf.Name), zap.Error(err))
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to start duration test"))
		}
	default:
		l.logger.Error("Unsupported test type requested", zap.String("name", l.conf.Name), zap.String("test", req.Test.String()))
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unsupported test type %s", req.Test.String()))
	}
	return &gen.StartTestResponse{}, nil
}

func (l *Light) StopTest(context.Context, *gen.StopTestRequest) (*gen.StopTestResponse, error) {
	l.logger.Info("Stopping test for light", zap.String("name", l.conf.Name))
	err := l.stopTest()
	if err != nil {
		l.logger.Error("Failed to stop test", zap.String("name", l.conf.Name), zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to stop test")
	}
	return &gen.StopTestResponse{}, nil
}

func (l *Light) GetTestResult(_ context.Context, req *gen.GetTestResultRequest) (*gen.TestResult, error) {

	var eState *EmergencyState
	var err error
	switch req.Test {
	case gen.EmergencyStatus_FUNCTION_TEST:
		eState, err = l.getFunctionTestResult()
		if err != nil {
			l.logger.Error("Failed to get function test result", zap.String("name", l.conf.Name), zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to get function test result")
		}
	case gen.EmergencyStatus_DURATION_TEST:
		eState, err = l.getDurationTestResult()
		if err != nil {
			l.logger.Error("Failed to get duration test result", zap.String("name", l.conf.Name), zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to get duration test result")
		}
	default:
		l.logger.Error("Unsupported test type requested", zap.String("name", l.conf.Name), zap.String("test", req.Test.String()))
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unsupported test type %s", req.Test.String()))
	}

	pass := false
	if eState != nil && *eState == Pass {
		pass = true
	}
	return &gen.TestResult{
		Test: req.Test,
		Pass: pass,
	}, nil
}
