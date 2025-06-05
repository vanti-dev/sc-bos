package merge

import (
	"context"
	"encoding/json"
	"errors"
	"math"

	"github.com/pborman/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vanti-dev/gobacnet"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/comm"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/status"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/securityevent"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

type securityEventSource struct {
	*config.ValueSource

	Description string `json:"description"` // description of the security event, must be set to make sense of the event
	// You want to set at least 1 of these OkBound values, else a security event will never get raised.
	OkLowerBound *float64 `json:"okLowerBound,omitempty"` // if the point is equal to or greater than this value, it is ok.
	OkUpperBound *float64 `json:"okUpperBound,omitempty"` // if the point is equal to or less than this value, it is ok.

	Actor     *string `json:"actor,omitempty"`     // Optional. Actor of the security event, e.g. "John Doe"
	EventType *int32  `json:"eventType,omitempty"` // Optional. the type of event, must be one of gen.SecurityEvent_EventType
	Priority  *int32  `json:"priority,omitempty"`  // Optional. Priority of the security event, lower is more important
	Source    *string `json:"source,omitempty"`    // Optional. Source of the security event, e.g. "Door 1"

	IsActive bool
}

type securityEventConfig struct {
	config.Trait
	SecurityEventSources []*securityEventSource `json:"securityEventSources"`
}

func readSecurityEventConfig(raw []byte) (cfg securityEventConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		return
	}

	if cfg.SecurityEventSources == nil || len(cfg.SecurityEventSources) == 0 {
		return cfg, errors.New("no security events configured")
	}

	for _, se := range cfg.SecurityEventSources {
		if se == nil {
			return cfg, errors.New("nil security event source found in config")
		}
		if se.ValueSource == nil {
			return cfg, errors.New("no value source provided for security event")
		}
		if se.Description == "" {
			return cfg, errors.New("no description provided for security event")
		}
		if se.OkLowerBound == nil {
			lowerInf := math.Inf(-1)
			se.OkLowerBound = &lowerInf
		}
		if se.OkUpperBound == nil {
			upperInf := math.Inf(1)
			se.OkUpperBound = &upperInf
		}
	}

	return
}

type securityEventImpl struct {
	client   *gobacnet.Client
	known    known.Context
	statuses *statuspb.Map
	logger   *zap.Logger

	model *securityevent.Model
	*securityevent.ModelServer
	config   securityEventConfig
	pollTask *task.Intermittent
}

func newSecurityEvent(client *gobacnet.Client, devices known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*securityEventImpl, error) {
	cfg, err := readSecurityEventConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	model := securityevent.NewModel()
	t := &securityEventImpl{
		client:      client,
		known:       devices,
		statuses:    statuses,
		logger:      logger,
		model:       model,
		ModelServer: securityevent.NewModelServer(model),
		config:      cfg,
	}
	t.pollTask = task.NewIntermittent(t.startPoll)
	initTraitStatus(statuses, cfg.Name, "SecurityEvent")
	return t, nil
}

func (s *securityEventImpl) startPoll(init context.Context) (stop task.StopFn, err error) {
	return startPoll(init, "securityevent", s.config.PollPeriodDuration(), s.config.PollTimeoutDuration(), s.logger, func(ctx context.Context) error {
		return s.pollPeer(ctx)
	})
}

func (s *securityEventImpl) AnnounceSelf(a node.Announcer) node.Undo {
	return a.Announce(s.config.Name, node.HasTrait(securityevent.TraitName, node.WithClients(gen.WrapSecurityEventApi(s))))
}

func (s *securityEventImpl) ListSecurityEvents(ctx context.Context, request *gen.ListSecurityEventsRequest) (*gen.ListSecurityEventsResponse, error) {
	err := s.pollPeer(ctx)
	if err != nil {
		return nil, err
	}
	return s.ModelServer.ListSecurityEvents(ctx, request)
}

func (s *securityEventImpl) PullSecurityEvents(request *gen.PullSecurityEventsRequest, server gen.SecurityEventApi_PullSecurityEventsServer) error {
	_ = s.pollTask.Attach(server.Context())
	return s.ModelServer.PullSecurityEvents(request, server)
}

func (s *securityEventImpl) pollPeer(ctx context.Context) error {
	var data []*gen.SecurityEvent
	var resProcessors []func(response any) error
	var readValues []config.ValueSource
	var requestNames []string

	for _, se := range s.config.SecurityEventSources {
		se := se
		requestNames = append(requestNames, se.ValueSource.String())
		readValues = append(readValues, *se.ValueSource)
		resProcessors = append(resProcessors, func(response any) error {
			event, err := se.checkResponseForSecurityEvent(response)
			if err != nil {
				return err
			}
			if event != nil {
				data = append(data, event)
			}
			return nil
		})
	}

	responses := comm.ReadProperties(ctx, s.client, s.known, readValues...)
	var errs []error
	for i, response := range responses {
		err := resProcessors[i](response)
		if err != nil {
			errs = append(errs, err)
		}
	}
	status.UpdatePollErrorStatus(s.statuses, s.config.Name, "SecurityEvent", requestNames, errs)
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	for _, se := range data {
		_, err := s.model.AddSecurityEvent(se)
		if err != nil {
			errs = append(errs, err)
			s.logger.Warn("failed to add security event", zap.Error(err), zap.String("eventId", se.Id), zap.String("description", se.Description))
		}
	}
	return errors.Join(errs...)
}

func (cfg *securityEventSource) checkResponseForSecurityEvent(response any) (*gen.SecurityEvent, error) {
	data := &gen.SecurityEvent{}

	value, err := comm.Float64Value(response)
	if err != nil {
		return nil, comm.ErrReadProperty{Prop: "securityEvent", Cause: err}
	}

	if value < *cfg.OkLowerBound ||
		value > *cfg.OkUpperBound {
		// we only want to add a new security event if it wasn't active on the last poll
		if !cfg.IsActive {
			cfg.IsActive = true
			data.Description = cfg.Description
			data.Id = uuid.New()
			data.SecurityEventTime = timestamppb.Now() // not strictly true but as long as the bacnet poll is not too slow, this should be fine

			if cfg.Actor != nil {
				data.Actor = &gen.Actor{
					Name: *cfg.Actor,
				}
			}

			if cfg.EventType != nil {
				data.EventType = gen.SecurityEvent_EventType(*cfg.EventType)
			}

			if cfg.Priority != nil {
				data.Priority = *cfg.Priority
			}

			if cfg.Source != nil {
				data.Source = &gen.SecurityEvent_Source{
					Name: *cfg.Source,
				}
			}
			return data, nil
		}
	} else {
		cfg.IsActive = false
	}
	return nil, nil
}
