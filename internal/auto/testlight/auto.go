package testlight

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
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const AutoType = "testlights"

var Factory = auto.FactoryFunc(func(services auto.Services) task.Starter {
	return ScanEmergencyTests(services.Node, services.Database, services.Logger)
})

func ScanEmergencyTests(clients node.Clienter, db *bolthold.Store, logger *zap.Logger) *EmergencyTestAutomation {
	configSend := make(chan Config)
	return &EmergencyTestAutomation{
		logger:  logger,
		db:      db,
		clients: clients,
		state:   state.NewManager(runstate.Idle),

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

func (r *runner) process(ctx context.Context, name string, test rpc.Test) error {
	var client rpc.DaliApiClient
	err := r.clients.Client(&client)
	if err != nil {
		return err
	}

	scanTime := time.Now()
	result, err := client.GetTestResult(ctx, &rpc.GetTestResultRequest{
		Name: name,
		Test: test,
	})
	if code := status.Code(err); code == codes.NotFound {
		result = nil
	} else if err != nil {
		return fmt.Errorf("GetTestResult: %w", err)
	}

	// we have successfully scanned for new data, so update last scan time
	lastScan, _, err := updateScanTime(r.db, name, test, scanTime)
	if err != nil {
		r.logger.Warn("failed to update last scan time - test After times may be inaccurate",
			zap.Error(err), zap.String("deviceName", name), zap.Any("test", test))
	}

	if result != nil {
		r.logger.Debug("emergency light has some test data",
			zap.String("deviceName", name),
			zap.String("test", test.String()),
			zap.Bool("pass", result.Pass))

		// store the test in the database
		record := TestResultRecord{
			Name:     name,
			Kind:     test,
			After:    lastScan,
			Before:   scanTime,
			TestPass: result.Pass,
		}
		if result.Duration != nil {
			record.AchievedDuration = result.Duration.AsDuration()
		}
		err = saveTestResult(r.db, record)
		if err != nil {
			return fmt.Errorf("store test result: %w", err)
		}

		if record.TestPass {
			// clear the test result from the emergency light now we have safely captured it
			// this can only be done for passes; failure status can only be cleared by a subsequent pass
			_, err = client.DeleteTestResult(ctx, &rpc.DeleteTestResultRequest{
				Name: name,
				Test: test,
			})
			if err != nil {
				return fmt.Errorf("DeleteTestResult: %w", err)
			}
		}
	}

	return nil
}

func (r *runner) runOneLoop(ctx context.Context, test rpc.Test) error {
	var errs error

	ticker := time.NewTicker(r.config.PollInterval.Duration)
	defer ticker.Stop()

	for _, name := range r.config.Devices {
		err := r.process(ctx, name, test)
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

func (r *runner) run(ctx context.Context) {
	ticker := time.NewTicker(r.config.CycleInterval.Duration)
	defer ticker.Stop()
	for {
		err := r.runOneLoop(ctx, rpc.Test_DURATION_TEST)
		if err != nil {
			r.logger.Error("errors retrieving duration test results", zap.Error(err))
		}

		err = r.runOneLoop(ctx, rpc.Test_FUNCTION_TEST)
		if err != nil {
			r.logger.Error("errors retrieving function test results", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			r.logger.Info("EmergencyTestAutomation stopping because its context was cancelled")
			return
		case newConfig := <-r.configUpdates:
			r.logger.Debug("EmergencyTestAutomation got a config update")
			r.config = newConfig
			ticker.Reset(r.config.CycleInterval.Duration)
		case <-ticker.C:
		}

	}
}
