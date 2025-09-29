package merge

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/comm"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/status"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/accesspb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

type accessConfig struct {
	config.Trait

	IngressPermitted     *config.ValueSource `json:"ingressPermitted,omitempty"`
	IngressPermittedType *string             `json:"ingressPermittedType,omitempty"`
	IngressDenied        *config.ValueSource `json:"ingressDenied,omitempty"`
	IngressDeniedType    *string             `json:"ingressDeniedType,omitempty"`
}

func readAccessConfig(raw []byte) (cfg accessConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	return
}

type access struct {
	client   *gobacnet.Client
	known    known.Context
	statuses *statuspb.Map
	logger   *zap.Logger

	model *accesspb.Model
	*accesspb.ModelServer
	cfg      accessConfig
	pollTask *task.Intermittent
}

func newAccess(client *gobacnet.Client, devices known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*access, error) {
	cfg, err := readAccessConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	model := accesspb.NewModel()
	a := &access{
		client:      client,
		known:       devices,
		statuses:    statuses,
		logger:      logger,
		model:       model,
		ModelServer: accesspb.NewModelServer(model),
		cfg:         cfg,
	}
	a.pollTask = task.NewIntermittent(a.startPoll)
	initTraitStatus(statuses, cfg.Name, "Access")
	return a, nil
}

func (a *access) startPoll(init context.Context) (stop task.StopFn, err error) {
	return startPoll(init, "access", a.cfg.PollPeriodDuration(), a.cfg.PollTimeoutDuration(), a.logger, func(ctx context.Context) error {
		_, err := a.pollPeer(ctx)
		return err
	})
}

func (a *access) AnnounceSelf(ann node.Announcer) node.Undo {
	return ann.Announce(a.cfg.Name, node.HasTrait(accesspb.TraitName, node.WithClients(gen.WrapAccessApi(a))))
}

func (a *access) GetLastAccessAttempt(ctx context.Context, request *gen.GetLastAccessAttemptRequest) (*gen.AccessAttempt, error) {
	_, err := a.pollPeer(ctx)
	if err != nil {
		return nil, err
	}
	return a.ModelServer.GetLastAccessAttempt(ctx, request)
}

func (a *access) PullAccessAttempts(request *gen.PullAccessAttemptsRequest, server gen.AccessApi_PullAccessAttemptsServer) error {
	_ = a.pollTask.Attach(server.Context())
	return a.ModelServer.PullAccessAttempts(request, server)
}

func (a *access) pollPeer(ctx context.Context) (*gen.AccessAttempt, error) {
	data := &gen.AccessAttempt{}

	var resProcessors []func(response any, data *gen.AccessAttempt) error
	var readValues []config.ValueSource
	var requestNames []string

	if a.cfg.IngressPermitted != nil {
		requestNames = append(requestNames, "ingressPermitted")
		readValues = append(readValues, *a.cfg.IngressPermitted)
		resProcessors = append(resProcessors, a.processIngressPermitted)
	}

	if a.cfg.IngressDenied != nil {
		requestNames = append(requestNames, "ingressDenied")
		readValues = append(readValues, *a.cfg.IngressDenied)
		resProcessors = append(resProcessors, a.processIngressDenied)
	}

	responses := comm.ReadProperties(ctx, a.client, a.known, readValues...)
	var errs []error
	for i, response := range responses {
		err := resProcessors[i](response, data)
		if err != nil {
			errs = append(errs, err)
		}
	}

	status.UpdatePollErrorStatus(a.statuses, a.cfg.Name, "Access", requestNames, errs)
	if len(errs) > 0 {
		return nil, multierr.Combine(errs...)
	}

	return a.model.UpdateLastAccessAttempt(data)
}

func (a *access) processIngressPermitted(response any, data *gen.AccessAttempt) error {
	value, ok := response.(string)
	if !ok {
		return comm.ErrReadProperty{Prop: "ingressPermitted", Cause: fmt.Errorf("converting to string")}
	}
	data.Grant = gen.AccessAttempt_GRANTED
	data.AccessAttemptTime = timestamppb.Now()
	data.Actor = &gen.Actor{
		LastGrantTime: timestamppb.Now(),
	}
	if a.cfg.IngressDeniedType != nil {
		data.Actor.Ids = map[string]string{
			*a.cfg.IngressDeniedType: value,
		}
	}
	return nil
}

func (a *access) processIngressDenied(response any, data *gen.AccessAttempt) error {
	value, ok := response.(string)
	if !ok {
		return comm.ErrReadProperty{Prop: "ingressDenied", Cause: fmt.Errorf("converting to string")}
	}
	data.Grant = gen.AccessAttempt_DENIED
	data.AccessAttemptTime = timestamppb.Now()

	if a.cfg.IngressDeniedType != nil {
		data.Actor = &gen.Actor{
			Ids: map[string]string{
				*a.cfg.IngressDeniedType: value,
			},
		}
	}
	return nil
}
