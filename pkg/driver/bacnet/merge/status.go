package merge

import (
	"context"
	"encoding/json"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/smart-core-os/gobacnet"
	"github.com/smart-core-os/gobacnet/property"
	"github.com/smart-core-os/sc-bos/pkg/driver/bacnet/config"
	"github.com/smart-core-os/sc-bos/pkg/driver/bacnet/known"
	"github.com/smart-core-os/sc-bos/pkg/driver/bacnet/status"
	"github.com/smart-core-os/sc-bos/pkg/gentrait/statuspb"
	"github.com/smart-core-os/sc-bos/pkg/node"
	"github.com/smart-core-os/sc-bos/pkg/task"
)

type statusConfig struct {
	config.Trait
	Tests []statusTest `json:"tests,omitempty"`
}

type statusTest struct {
	*status.Object
	*config.ValueSource
	HasLimits bool `json:"hasLimits,omitempty"`
}

func (t statusTest) ToObject() *status.Object {
	o := &status.Object{}
	if vs := t.ValueSource; vs != nil {
		ptr := func(p property.ID) *config.PropertyID {
			cp := config.PropertyID(p)
			return &cp
		}
		o.Test = "BACnetObject"
		o.EventState = &config.ValueSource{Object: vs.Object, Device: vs.Device, Property: ptr(property.EventState)}
		o.Reliability = &config.ValueSource{Object: vs.Object, Device: vs.Device, Property: ptr(property.Reliability)}
		o.OutOfService = &config.ValueSource{Object: vs.Object, Device: vs.Device, Property: ptr(property.OutOfService)}
		if t.HasLimits {
			o.HighLimit = &config.ValueSource{Object: vs.Object, Device: vs.Device, Property: ptr(property.HighLimit)}
			o.LowLimit = &config.ValueSource{Object: vs.Object, Device: vs.Device, Property: ptr(property.LowLimit)}
		}
	}
	if so := t.Object; so != nil {
		if v := so.Name; v != "" {
			o.Name = v
		}
		if v := so.Test; v != "" {
			o.Test = v
		}
		if v := so.Level; v != 0 {
			o.Level = v
		}
		if v := so.EventState; v != nil {
			o.EventState = v
		}
		if v := so.Reliability; v != nil {
			o.Reliability = v
		}
		if v := so.OutOfService; v != nil {
			o.OutOfService = v
		}
		if v := so.HighLimit; v != nil {
			o.HighLimit = v
		}
		if v := so.LowLimit; v != nil {
			o.LowLimit = v
		}
		o.Value = so.Value
		o.NominalValue = so.NominalValue
	}
	return o
}

func readStatusConfig(raw []byte) (cfg statusConfig, err error) {
	err = json.Unmarshal(raw, &cfg)

	for i, test := range cfg.Tests {
		if test.Object == nil {
			test.Object = &status.Object{}
		}
		test.Object.Name = cfg.Name
		cfg.Tests[i] = test
	}
	return
}

type statusImpl struct {
	cfg      statusConfig
	devices  known.Context
	statuses *statuspb.Map
	logger   *zap.Logger
	pollTask *task.Intermittent

	setup   atomic.Bool
	monitor *status.Monitor
}

func newStatus(client *gobacnet.Client, devices known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*statusImpl, error) {
	cfg, err := readStatusConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	monitor := status.NewMonitor(client, devices, statuses)
	monitor.Logger = logger
	impl := &statusImpl{
		cfg:      cfg,
		devices:  devices,
		statuses: statuses,
		logger:   logger,
		monitor:  monitor,
	}
	impl.pollTask = task.NewIntermittent(impl.startPoll)
	initTraitStatus(statuses, cfg.Name, "Status")
	return impl, nil
}

// AnnounceSelf implements node.SelfAnnouncer but does not announce anything.
// Announcement is handled by the statuspb.Map passed during construction.
// Instead we start the monitoring of the status map and rely on that to know when to poll our device.
func (impl *statusImpl) AnnounceSelf(_ node.Announcer) node.Undo {
	ctx, stop := context.WithCancel(context.Background())
	go func() {
		for event := range impl.statuses.WatchEvents(ctx) {
			if event.Name == impl.cfg.Name {
				err := impl.pollTask.Attach(event.Ctx)
				if err != nil {
					impl.logger.Error("failed to start poll task", zap.String("name", impl.cfg.Name), zap.Error(err))
				}
			}
		}
	}()
	return func() {
		stop()
	}
}

func (impl *statusImpl) setupMonitor() error {
	if !impl.setup.CompareAndSwap(false, true) {
		return nil // already setup
	}

	for _, test := range impl.cfg.Tests {
		impl.monitor.AddObject(test.ToObject())
	}

	return nil
}

func (impl *statusImpl) startPoll(init context.Context) (stop task.StopFn, err error) {
	return startPoll(init, "status", impl.cfg.PollPeriodDuration(), impl.cfg.PollTimeoutDuration(), impl.logger, func(ctx context.Context) error {
		return impl.pollPeer(ctx)
	})
}

func (impl *statusImpl) pollPeer(ctx context.Context) error {
	if err := impl.setupMonitor(); err != nil {
		return err
	}
	return impl.monitor.Poll(ctx)
}
