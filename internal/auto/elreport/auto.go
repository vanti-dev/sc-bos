package elreport

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/timshannon/bolthold"
	"github.com/vanti-dev/bsp-ew/internal/auto"
	"github.com/vanti-dev/bsp-ew/internal/auto/runstate"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/rpc"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/task"
	"github.com/vanti-dev/bsp-ew/internal/util/minibus"
	"github.com/vanti-dev/bsp-ew/internal/util/state"
	"github.com/vanti-dev/bsp-ew/pkg/gen"
	"go.etcd.io/bbolt"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const AutoType = "testlights"

var Factory = auto.FactoryFunc(func(services auto.Services) task.Starter {
	return ScanEmergencyTests(services.Node, services.Database, services.GRPCServices, services.Logger)
})

func ScanEmergencyTests(clients node.Clienter, db *bolthold.Store, registrar grpc.ServiceRegistrar, logger *zap.Logger) *EmergencyTestAutomation {
	configSend := make(chan Config)
	return &EmergencyTestAutomation{
		logger:  logger,
		db:      db,
		clients: clients,
		state:   state.NewManager(runstate.Idle),

		registrar: registrar,
		server: &server{
			db:     db,
			logger: logger.Named("server"),
		},

		config:            DefaultConfig(),
		stop:              func() {},
		configChangesSend: configSend,
		configChangesRecv: minibus.DropExcess(configSend),
	}
}

type EmergencyTestAutomation struct {
	logger  *zap.Logger
	db      *bolthold.Store
	clients node.Clienter
	state   *state.Manager[runstate.RunState]

	server       *server
	registerOnce sync.Once
	registrar    grpc.ServiceRegistrar

	m                 sync.Mutex
	config            Config
	stop              context.CancelFunc
	configChangesSend chan<- Config
	configChangesRecv <-chan Config
}

func (a *EmergencyTestAutomation) CurrentState() runstate.RunState {
	return a.state.CurrentState()
}

func (a *EmergencyTestAutomation) WaitForStateChange(ctx context.Context, source runstate.RunState) error {
	return a.state.WaitForStateChange(ctx, source)
}

func (a *EmergencyTestAutomation) Start(_ context.Context) error {
	a.registerOnce.Do(func() {
		a.server.register(a.registrar)
	})

	a.m.Lock()
	defer a.m.Unlock()
	if s := a.state.CurrentState(); s == runstate.Running || s == runstate.Starting {
		return errors.New("already running")
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.stop = cancel
	initialConfig := a.config
	_ = a.state.Update(runstate.Starting)
	go func() {
		defer a.state.Update(runstate.Stopped)

		r := &runner{
			logger:        a.logger,
			db:            a.db,
			clients:       a.clients,
			config:        initialConfig,
			configUpdates: a.configChangesRecv,
		}
		a.state.Update(runstate.Running)
		r.run(ctx)
	}()

	return nil
}

func (a *EmergencyTestAutomation) Stop() error {
	a.m.Lock()
	defer a.m.Unlock()
	if s := a.state.CurrentState(); s != runstate.Running && s != runstate.Starting {
		return errors.New("not running")
	}
	a.stop()
	return nil
}

func (a *EmergencyTestAutomation) Configure(raw []byte) error {
	parsed, err := DecodeConfig(raw)
	if err != nil {
		return err
	}

	a.m.Lock()
	a.config = parsed
	a.m.Unlock()

	// due to using minibus.DropExcess, this will never block a long time
	a.configChangesSend <- parsed

	return nil
}

type runner struct {
	logger        *zap.Logger
	db            *bolthold.Store
	clients       node.Clienter
	config        Config
	configUpdates <-chan Config
}

// process polls a single emergency light for new data, and saves that data to the database where required.
func (r *runner) process(ctx context.Context, name string) error {
	var client rpc.DaliApiClient
	err := r.clients.Client(&client)
	if err != nil {
		return err
	}
	scanTime := time.Now()
	logger := r.logger.With(zap.String("deviceName", name))

	result := pollDaliEmergency(ctx, client, name)

	// update the database with the test results
	err = r.db.Bolt().Update(func(tx *bbolt.Tx) error {
		// update the status for the light
		statusChanged, err := updateLatestStatus(r.db, tx, name, scanTime, result.faults)
		if err != nil {
			logger.Error("error updating emergency light latest status", zap.Error(err),
				zap.Any("faults", result.faults))
			return err
		}

		// if the status has changed, record it in the log
		if statusChanged {
			err = saveStatusReport(r.db, tx, name, scanTime, result.faults)
			if err != nil {
				logger.Error("failed to save emergency light status report", zap.Error(err))
				return err
			}
		}

		// record any test passes
		if result.functionTestPass {
			err = saveFunctionTestPass(r.db, tx, name, scanTime)
			if err != nil {
				logger.Error("failed to save emergency light function test pass", zap.Error(err))
				return err
			}
		}
		if result.durationTestPass {
			err = saveDurationTestPass(r.db, tx, name, scanTime, result.durationTestResult)
			if err != nil {
				logger.Error("failed to save emergency light duration test pass", zap.Error(err),
					zap.Duration("result", result.durationTestResult))
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Error("database transaction error", zap.Error(err))
		return err
	}

	// now we have saved test passes in the database, we can delete them from the light without data loss
	var errs error
	if result.functionTestPass {
		_, err = client.DeleteTestResult(ctx, &rpc.DeleteTestResultRequest{
			Name: name,
			Test: rpc.Test_FUNCTION_TEST,
		})
		if err != nil {
			logger.Warn("couldn't delete function test pass from light; duplicate records may occur", zap.Error(err))
			errs = multierr.Append(errs, err)
		}
	}
	if result.durationTestPass {
		_, err = client.DeleteTestResult(ctx, &rpc.DeleteTestResultRequest{
			Name: name,
			Test: rpc.Test_DURATION_TEST,
		})
		if err != nil {
			logger.Warn("couldn't delete duration test pass from light; duplicate records may occur", zap.Error(err))
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

// runOneLoop will poll every emergency light in r.config.Devices for updates status & test data.
// Lights are polled at intervals r.config.PollInterval, to leave clear time in between where other commands can get
// through, as this polling can wait a little while without issue.
func (r *runner) runOneLoop(ctx context.Context) error {
	var errs error

	ticker := time.NewTicker(r.config.PollInterval.Duration)
	defer ticker.Stop()

	for _, name := range r.config.Devices {
		err := r.process(ctx, name)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("%q: %w", name, err))
		}

		select {
		case <-ctx.Done():
			return errs
		case <-ticker.C:
		}
	}

	return errs
}

// run repeatedly calls r.runOneLoop, with an interval of r.Config.CycleInterval between each call
// stops once ctx is cancelled.
func (r *runner) run(ctx context.Context) {
	ticker := time.NewTicker(r.config.CycleInterval.Duration)
	defer ticker.Stop()
	for {
		err := r.runOneLoop(ctx)
		if err != nil {
			r.logger.Warn("errors polling for emergency test results", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			return
		case newConfig := <-r.configUpdates:
			r.logger.Debug("EmergencyTestAutomation got a config update")
			r.config = newConfig
			ticker.Reset(r.config.CycleInterval.Duration)
		case <-ticker.C:
		}

	}
}

type pollResult struct {
	faults             []gen.EmergencyLightFault
	errors             error
	functionTestPass   bool
	durationTestPass   bool
	durationTestResult time.Duration // non-zero if durationTestPass is true
}

// pollDaliEmergency gets new data from an emergency light
func pollDaliEmergency(ctx context.Context, client rpc.DaliApiClient, name string) pollResult {
	var result pollResult

	// get the emergency light status
	emStatus, err := client.GetEmergencyStatus(ctx, &rpc.GetEmergencyStatusRequest{Name: name})
	if err != nil {
		result.faults = []gen.EmergencyLightFault{gen.EmergencyLightFault_COMMUNICATION_FAILURE}
		result.errors = multierr.Append(result.errors, err)
		return result
	} else {
		result.faults = translateFaults(emStatus.Failures)
	}

	// check for any test passes
	funcRes, err := client.GetTestResult(ctx, &rpc.GetTestResultRequest{
		Name: name,
		Test: rpc.Test_FUNCTION_TEST,
	})
	if err == nil {
		result.functionTestPass = funcRes.Pass
	} else if status.Code(err) != codes.NotFound {
		result.errors = multierr.Append(result.errors, err)
	}

	durRes, err := client.GetTestResult(ctx, &rpc.GetTestResultRequest{
		Name: name,
		Test: rpc.Test_DURATION_TEST,
	})
	if err == nil {
		result.durationTestPass = durRes.Pass
		result.durationTestResult = durRes.Duration.AsDuration()
	} else if status.Code(err) != codes.NotFound {
		result.errors = multierr.Append(result.errors, err)
	}

	return result
}

func translateFaults(daliFailures []rpc.EmergencyStatus_Failure) (faults []gen.EmergencyLightFault) {
	for _, failure := range daliFailures {
		switch failure {
		case rpc.EmergencyStatus_DURATION_TEST_FAILED:
			faults = append(faults, gen.EmergencyLightFault_DURATION_TEST_FAILED)
		case rpc.EmergencyStatus_FUNCTION_TEST_FAILED:
			faults = append(faults, gen.EmergencyLightFault_FUNCTION_TEST_FAILED)
		case rpc.EmergencyStatus_BATTERY_DURATION_FAILURE, rpc.EmergencyStatus_BATTERY_FAILURE:
			faults = append(faults, gen.EmergencyLightFault_BATTERY_FAULT)
		case rpc.EmergencyStatus_LAMP_FAILURE:
			faults = append(faults, gen.EmergencyLightFault_LAMP_FAULT)
		default:
			faults = append(faults, gen.EmergencyLightFault_OTHER_FAULT)
		}
	}

	faults = sortDeduplicateFaults(faults)
	return
}
