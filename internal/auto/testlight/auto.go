package testlight

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/timshannon/bolthold"
	"github.com/vanti-dev/bsp-ew/internal/auto"
	"github.com/vanti-dev/bsp-ew/internal/driver/tc3dali/rpc"
	"github.com/vanti-dev/bsp-ew/internal/node"
	"github.com/vanti-dev/bsp-ew/internal/task"
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
	return &EmergencyTestAutomation{
		logger:  logger,
		db:      db,
		clients: clients,
		config:  DefaultConfig(),
	}
}

type EmergencyTestAutomation struct {
	logger  *zap.Logger
	db      *bolthold.Store
	clients node.Clienter
	config  Config

	stateM  sync.Mutex
	running bool
	cancel  context.CancelFunc
}

func (a *EmergencyTestAutomation) Start(_ context.Context) error {
	a.stateM.Lock()
	defer a.stateM.Unlock()
	if a.running {
		return errors.New("already running")
	}

	runctx, cancel := context.WithCancel(context.Background())
	a.running = true
	a.cancel = cancel
	go func() {
		a.run(runctx)

		a.stateM.Lock()
		defer a.stateM.Unlock()
		a.running = false
		a.cancel = nil
	}()
	return nil
}

func (a *EmergencyTestAutomation) Stop() error {
	a.stateM.Lock()
	defer a.stateM.Unlock()
	if !a.running {
		return errors.New("not running")
	}
	a.cancel()
	return nil
}

func (a *EmergencyTestAutomation) Configure(raw []byte) error {
	config, err := DecodeConfig(raw)
	if err != nil {
		return err
	}
	a.config = config
	return nil
}

func (a *EmergencyTestAutomation) process(ctx context.Context, name string, test rpc.Test) error {
	var client rpc.DaliApiClient
	err := a.clients.Client(&client)
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
	lastScan, _, err := updateScanTime(a.db, name, test, scanTime)
	if err != nil {
		a.logger.Warn("failed to update last scan time - test CompletedAfter times may be inaccurate",
			zap.Error(err), zap.String("deviceName", name), zap.Any("test", test))
	}

	if result != nil {
		// store the test in the database
		var achievedDuration time.Duration
		if result.Duration != nil {
			achievedDuration = result.Duration.AsDuration()
		}
		err = a.db.Insert(bolthold.NextSequence(), TestResultRecord{
			Name:             name,
			Kind:             test,
			CompletedAfter:   lastScan,
			CompletedBefore:  scanTime,
			Success:          result.Pass,
			AchievedDuration: achievedDuration,
		})
		if err != nil {
			return fmt.Errorf("store test result: %w", err)
		}

		// clear the test result from the emergency light now we have safely captured it
		_, err = client.DeleteTestResult(ctx, &rpc.DeleteTestResultRequest{
			Name: name,
			Test: test,
		})
		if err != nil {
			return fmt.Errorf("DeleteTestResult: %w", err)
		}
	}

	return nil
}

func (a *EmergencyTestAutomation) runOneLoop(ctx context.Context, test rpc.Test) error {
	var errs error

	ticker := time.NewTicker(a.config.Interval)
	defer ticker.Stop()

	for _, name := range a.config.Devices {
		err := a.process(ctx, name, test)
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

func (a *EmergencyTestAutomation) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		err := a.runOneLoop(ctx, rpc.Test_DURATION_TEST)
		if err != nil {
			a.logger.Error("errors retrieving duration test results", zap.Error(err))
		}

		err = a.runOneLoop(ctx, rpc.Test_FUNCTION_TEST)
		if err != nil {
			a.logger.Error("errors retrieving function test results", zap.Error(err))
		}
	}
}
