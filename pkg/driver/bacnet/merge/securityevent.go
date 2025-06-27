package merge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	// description of the security event, must be set to make sense of the event
	// if IsString is true, this description acts as a prefix to the value source.
	Description string `json:"description"`
	// You want to set at least 1 of these OkBound values, else a security event will never get raised.
	OkLowerBound *float64 `json:"okLowerBound,omitempty"` // if the point is equal to or greater than this value, it is ok.
	OkUpperBound *float64 `json:"okUpperBound,omitempty"` // if the point is equal to or less than this value, it is ok.

	Actor     *string                      `json:"actor,omitempty"`     // Optional. Actor of the security event, e.g. "John Doe"
	EventType *gen.SecurityEvent_EventType `json:"eventType,omitempty"` // Optional. the type of event, must be one of gen.SecurityEvent_EventType
	Priority  *int32                       `json:"priority,omitempty"`  // Optional. Priority of the security event, lower is more important
	Source    *string                      `json:"source,omitempty"`    // Optional. Source of the security event, e.g. "Door 1"

	IsString      bool   `json:"isString,omitempty"`      // Optional. If true, the value source is a string, otherwise it is a float64.
	OkStringValue string `json:"okStringValue,omitempty"` // if IsString is true, this is the value that is considered OK i.e. will deactivate the security event.
}

type securityEvent struct {
	cfg      securityEventSource
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
	events   []*securityEvent
	pollTask *task.Intermittent
}

func newSecurityEvent(client *gobacnet.Client, devices known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*securityEventImpl, error) {
	cfg, err := readSecurityEventConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	model := securityevent.NewModel()
	var events []*securityEvent
	for _, se := range cfg.SecurityEventSources {
		e := &securityEvent{
			cfg:      *se,
			IsActive: false,
		}
		events = append(events, e)
	}
	t := &securityEventImpl{
		client:      client,
		known:       devices,
		statuses:    statuses,
		logger:      logger,
		model:       model,
		ModelServer: securityevent.NewModelServer(model),
		config:      cfg,
		events:      events,
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

	for _, se := range s.events {
		se := se
		requestNames = append(requestNames, se.cfg.ValueSource.String())
		readValues = append(readValues, *se.cfg.ValueSource)
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

func (se *securityEvent) checkResponseForSecurityEvent(response any) (*gen.SecurityEvent, error) {
	data := &gen.SecurityEvent{}

	if se.cfg.IsString {
		value, err := comm.StringValue(response)
		if err != nil {
			return nil, comm.ErrReadProperty{Prop: "securityEvent", Cause: err}
		}

		if value == se.cfg.OkStringValue {
			se.IsActive = false
			return nil, nil // no event if the value is empty
		}

		if se.IsActive {
			// if the security event is already active, we don't want to create a new one
			// TODO: this also means that a new event on the existing security event source will be ignored
			// e.g. a different card used to scan that is rejected won't reveal a new card ID
			return nil, nil
		}

		se.IsActive = true

		se.cfgDefaults(data)

		if se.cfg.Description != "" {
			data.Description = fmt.Sprintf("%s: %s", se.cfg.Description, value)
		} else {
			data.Description = value
		}

		return data, nil
	}

	value, err := comm.Float64Value(response)
	if err != nil {
		return nil, comm.ErrReadProperty{Prop: "securityEvent", Cause: err}
	}

	if value < *se.cfg.OkLowerBound ||
		value > *se.cfg.OkUpperBound {
		// we only want to add a new security event if it wasn't active on the last poll
		if !se.IsActive {
			se.IsActive = true
			se.cfgDefaults(data)
			data.Description = se.cfg.Description
			return data, nil
		}
	} else {
		se.IsActive = false
	}
	return nil, nil
}

func (se *securityEvent) cfgDefaults(data *gen.SecurityEvent) {
	data.Id = uuid.New()
	data.SecurityEventTime = timestamppb.Now() // not strictly true but as long as the bacnet poll is not too slow, this should be fine

	if se.cfg.Actor != nil {
		data.Actor = &gen.Actor{
			Name: *se.cfg.Actor,
		}
	}

	if se.cfg.EventType != nil {
		data.EventType = *se.cfg.EventType
	}

	if se.cfg.Priority != nil {
		data.Priority = *se.cfg.Priority
	}

	if se.cfg.Source != nil {
		data.Source = &gen.SecurityEvent_Source{
			Name: *se.cfg.Source,
		}
	}
}
