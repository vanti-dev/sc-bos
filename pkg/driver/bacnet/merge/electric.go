package merge

import (
	"context"
	"encoding/json"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/time/clock"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/electric"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/comm"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

type electricConfig struct {
	config.Trait
	Demand *electricDemandConfig `json:"demand,omitempty"`
}

type electricDemandConfig struct {
	*ElectricPhaseConfig                        // single phase
	Phases               [3]ElectricPhaseConfig `json:"phases,omitempty"`
}

type ElectricPhaseConfig struct {
	Current *config.ValueSource `json:"current,omitempty"`
	Voltage *config.ValueSource `json:"voltage,omitempty"`
	Rating  *config.ValueSource `json:"rating,omitempty"`

	PowerFactor   *config.ValueSource `json:"powerFactor,omitempty"`
	RealPower     *config.ValueSource `json:"realPower,omitempty"`
	ApparentPower *config.ValueSource `json:"apparentPower,omitempty"`
	ReactivePower *config.ValueSource `json:"reactivePower,omitempty"`
}

func readElectricConfig(raw []byte) (cfg electricConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	return
}

type electricTrait struct {
	client   *gobacnet.Client
	known    known.Context
	statuses *statuspb.Map
	logger   *zap.Logger

	model *electric.Model
	*electric.ModelServer
	config   electricConfig
	pollTask *task.Intermittent
}

func newElectric(client *gobacnet.Client, devices known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*electricTrait, error) {
	cfg, err := readElectricConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	model := electric.NewModel(clock.Real())
	_, _ = model.UpdateDemand(&traits.ElectricDemand{}) // reset defaults
	t := &electricTrait{
		client:      client,
		known:       devices,
		statuses:    statuses,
		logger:      logger,
		model:       model,
		ModelServer: electric.NewModelServer(model),
		config:      cfg,
	}
	t.pollTask = task.NewIntermittent(t.startPoll)
	return t, nil
}

func (t *electricTrait) AnnounceSelf(a node.Announcer) node.Undo {
	return a.Announce(t.config.Name, node.HasTrait(trait.Electric, node.WithClients(electric.WrapApi(t))))
}

func (t *electricTrait) GetDemand(ctx context.Context, request *traits.GetDemandRequest) (*traits.ElectricDemand, error) {
	_, err := t.pollPeer(ctx)
	if err != nil {
		return nil, err
	}
	return t.ModelServer.GetDemand(ctx, request)
}

func (t *electricTrait) PullDemand(request *traits.PullDemandRequest, server traits.ElectricApi_PullDemandServer) error {
	err := t.pollTask.Attach(server.Context())
	if err != nil {
		return err
	}

	// avoid returning the zero value if we are the first to attach since reboot
	timeoutCtx, cleanup := context.WithTimeout(server.Context(), t.config.PollTimeoutDuration())
	defer cleanup()
	for change := range t.model.PullDemand(timeoutCtx) {
		if !proto.Equal(change.Value, &traits.ElectricDemand{}) { // skip zero value
			break
		}
	}

	return t.ModelServer.PullDemand(request, server)
}

func (t *electricTrait) startPoll(init context.Context) (stop task.StopFn, err error) {
	return startPoll(init, "electric", t.config.PollPeriodDuration(), t.config.PollTimeoutDuration(), t.logger, func(ctx context.Context) error {
		_, err := t.pollPeer(ctx)
		return err
	})
}

func (t *electricTrait) pollPeer(ctx context.Context) (*traits.ElectricDemand, error) {
	var toRead []config.ValueSource
	var toWrite []func(v any) error // already scaled
	dst := &traits.ElectricDemand{}
	var phaseDemand [3]*traits.ElectricDemand

	if cfg := t.config.Demand; cfg != nil {
		if cfg.Current != nil {
			toRead = append(toRead, *cfg.Current)
			toWrite = append(toWrite, func(v any) (err error) {
				dst.Current, err = comm.Float32Value(v)
				return
			})
		}
		if cfg.Voltage != nil {
			toRead = append(toRead, *cfg.Voltage)
			toWrite = append(toWrite, func(v any) (err error) {
				dst.Voltage, err = ptr(comm.Float32Value(v))
				return
			})
		}
		if cfg.Rating != nil {
			toRead = append(toRead, *cfg.Rating)
			toWrite = append(toWrite, func(v any) (err error) {
				dst.Rating, err = comm.Float32Value(v)
				return
			})
		}

		readPhase := func(phase ElectricPhaseConfig, dst *traits.ElectricDemand) {
			if phase.PowerFactor != nil {
				toRead = append(toRead, *phase.PowerFactor)
				toWrite = append(toWrite, func(v any) (err error) {
					dst.PowerFactor, err = ptr(comm.Float32Value(v))
					return
				})
			}
			if phase.RealPower != nil {
				toRead = append(toRead, *phase.RealPower)
				toWrite = append(toWrite, func(v any) (err error) {
					dst.RealPower, err = ptr(comm.Float32Value(v))
					return
				})
			}
			if phase.ApparentPower != nil {
				toRead = append(toRead, *phase.ApparentPower)
				toWrite = append(toWrite, func(v any) (err error) {
					dst.ApparentPower, err = ptr(comm.Float32Value(v))
					return
				})
			}
			if phase.ReactivePower != nil {
				toRead = append(toRead, *phase.ReactivePower)
				toWrite = append(toWrite, func(v any) (err error) {
					dst.ReactivePower, err = ptr(comm.Float32Value(v))
					return
				})
			}
		}

		if cfg.ElectricPhaseConfig != nil {
			readPhase(*cfg.ElectricPhaseConfig, dst)
		}

		for i, phase := range cfg.Phases {
			dst := &traits.ElectricDemand{}
			phaseDemand[i] = dst
			readPhase(phase, dst)
		}
	}

	var errs []error
	for i, response := range comm.ReadPropertiesChunked(ctx, t.client, t.known, t.config.ChunkSize, toRead...) {
		err := toWrite[i](response)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// process phased electric devices.
	// If there isn't a dedicated power value source for the combination we add each phase to produce the total
	var power, apparent, reactive float32
	for _, demand := range phaseDemand {
		if demand == nil {
			continue
		}
		// note: don't do this for power factor as that is complicated!
		if demand.RealPower != nil {
			power += *demand.RealPower
		}
		if demand.ApparentPower != nil {
			apparent += *demand.ApparentPower
		}
		if demand.ReactivePower != nil {
			reactive += *demand.ReactivePower
		}
	}
	if dst.RealPower == nil && power != 0 {
		dst.RealPower = &power
	}
	if dst.ApparentPower == nil && apparent != 0 {
		dst.ApparentPower = &apparent
	}
	if dst.ReactivePower == nil && reactive != 0 {
		dst.ReactivePower = &reactive
	}

	comm.UpdatePollErrorStatus(t.statuses, t.config.Name, "Electric", len(toRead), errs...)
	if len(errs) > 0 {
		return nil, multierr.Combine(errs...)
	}

	return t.model.UpdateDemand(dst)
}
