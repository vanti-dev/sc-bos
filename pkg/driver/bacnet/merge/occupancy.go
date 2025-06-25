package merge

import (
	"context"
	"encoding/json"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensorpb"
	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/comm"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/status"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

type occupancyCfg struct {
	config.Trait
	OccupancyStatus *config.ValueSource `json:"occupancyStatus,omitempty"` // the point to read occupancy from
}

func readOccupancyConfig(raw []byte) (cfg occupancyCfg, err error) {
	err = json.Unmarshal(raw, &cfg)
	return
}

var _ traits.OccupancySensorApiServer = (*occupancy)(nil)

type occupancy struct {
	traits.UnimplementedOccupancySensorApiServer

	client   *gobacnet.Client
	known    known.Context
	statuses *statuspb.Map
	logger   *zap.Logger

	model *occupancysensorpb.Model
	*occupancysensorpb.ModelServer
	config   occupancyCfg
	pollTask *task.Intermittent
}

func newOccupancy(client *gobacnet.Client, known known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*occupancy, error) {
	cfg, err := readOccupancyConfig(config.Raw)
	if err != nil {
		return nil, err
	}

	model := occupancysensorpb.NewModel()

	o := &occupancy{
		client:      client,
		known:       known,
		statuses:    statuses,
		logger:      logger,
		model:       model,
		ModelServer: occupancysensorpb.NewModelServer(model),
		config:      cfg,
	}

	o.pollTask = task.NewIntermittent(o.startPoll)

	initTraitStatus(statuses, cfg.Name, "OccupancySensor")

	return o, nil
}

func (o *occupancy) AnnounceSelf(a node.Announcer) node.Undo {
	return a.Announce(o.config.Name, node.HasTrait(trait.OccupancySensor, node.WithClients(occupancysensorpb.WrapApi(o))))
}

func (o *occupancy) GetOccupancy(ctx context.Context, request *traits.GetOccupancyRequest) (*traits.Occupancy, error) {
	_, err := o.pollPeer(ctx)
	if err != nil {
		return nil, err
	}
	return o.ModelServer.GetOccupancy(ctx, request)
}

func (o *occupancy) PullOccupancy(request *traits.PullOccupancyRequest, server traits.OccupancySensorApi_PullOccupancyServer) error {
	_ = o.pollTask.Attach(server.Context())
	return o.ModelServer.PullOccupancy(request, server)
}

func (o *occupancy) startPoll(init context.Context) (stop task.StopFn, err error) {
	return startPoll(init, "occupancy", o.config.PollPeriodDuration(), o.config.PollTimeoutDuration(), o.logger, func(ctx context.Context) error {
		_, err := o.pollPeer(ctx)
		return err
	})
}

func (o *occupancy) pollPeer(ctx context.Context) (*traits.Occupancy, error) {
	data := &traits.Occupancy{}

	var resProcessors []func(response any) error
	var readValues []config.ValueSource
	var requestNames []string

	if o.config.OccupancyStatus != nil {
		requestNames = append(requestNames, "occupancy")
		readValues = append(readValues, *o.config.OccupancyStatus)
		resProcessors = append(resProcessors, func(response any) error {
			value, err := comm.BoolValue(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "occupancy", Cause: err}
			}

			data.State = traits.Occupancy_UNOCCUPIED

			if value {
				data.State = traits.Occupancy_OCCUPIED
			}

			return nil
		})
	}
	responses := comm.ReadProperties(ctx, o.client, o.known, readValues...)
	var errs []error
	for i, response := range responses {
		err := resProcessors[i](response)
		if err != nil {
			errs = append(errs, err)
		}
	}

	status.UpdatePollErrorStatus(o.statuses, o.config.Name, "Occupancy", requestNames, errs)
	if len(errs) > 0 {
		return nil, multierr.Combine(errs...)
	}

	return o.model.SetOccupancy(data)
}
