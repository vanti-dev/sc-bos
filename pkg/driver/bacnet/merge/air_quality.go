package merge

import (
	"context"
	"encoding/json"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/smart-core-os/gobacnet"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/cmp"
	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airqualitysensorpb"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/comm"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/config"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/known"
	"github.com/vanti-dev/sc-bos/pkg/driver/bacnet/status"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/task"
)

type airQualityConfig struct {
	config.Trait
	AirPressure *config.ValueSource `json:"airPressure,omitempty"`
	Co2         *config.ValueSource `json:"co2,omitempty"`
	IaqScore    *config.ValueSource `json:"iaqScore,omitempty"`
	// A measure of particles in the air measuring 10 microns or less in size, in micrograms per cubic meter.
	Pm10 *config.ValueSource `json:"pm10,omitempty"`
	// A measure of particles in the air measuring 2.5 microns or less in size, in micrograms per cubic meter.
	Pm25 *config.ValueSource `json:"pm25,omitempty"`
	VOC  *config.ValueSource `json:"voc,omitempty"`
}

func readAirQualitySensorConfig(raw []byte) (cfg airQualityConfig, err error) {
	err = json.Unmarshal(raw, &cfg)
	return
}

type airQualitySensor struct {
	client   *gobacnet.Client
	known    known.Context
	statuses *statuspb.Map
	logger   *zap.Logger

	model *airqualitysensorpb.Model
	*airqualitysensorpb.ModelServer
	config   airQualityConfig
	pollTask *task.Intermittent
}

func newAirQualitySensor(client *gobacnet.Client, devices known.Context, statuses *statuspb.Map, config config.RawTrait, logger *zap.Logger) (*airQualitySensor, error) {
	cfg, err := readAirQualitySensorConfig(config.Raw)
	if err != nil {
		return nil, err
	}
	model := airqualitysensorpb.NewModel(resource.WithMessageEquivalence(cmp.Equal(
		cmp.FloatValueApprox(0, 10), // report co2 changes of 10ppm or more
	)))
	t := &airQualitySensor{
		client:      client,
		known:       devices,
		statuses:    statuses,
		logger:      logger,
		model:       model,
		ModelServer: airqualitysensorpb.NewModelServer(model),
		config:      cfg,
	}
	t.pollTask = task.NewIntermittent(t.startPoll)
	initTraitStatus(statuses, cfg.Name, "AirQualitySensor")
	return t, nil
}

func (aq *airQualitySensor) startPoll(init context.Context) (stop task.StopFn, err error) {
	return startPoll(init, "air quality sensor", aq.config.PollPeriodDuration(), aq.config.PollTimeoutDuration(), aq.logger, func(ctx context.Context) error {
		_, err := aq.pollPeer(ctx)
		return err
	})
}

func (aq *airQualitySensor) AnnounceSelf(a node.Announcer) node.Undo {
	return a.Announce(aq.config.Name, node.HasTrait(trait.AirQualitySensor, node.WithClients(airqualitysensorpb.WrapApi(aq))))
}

func (aq *airQualitySensor) GetAirQuality(ctx context.Context, request *traits.GetAirQualityRequest) (*traits.AirQuality, error) {
	_, err := aq.pollPeer(ctx)
	if err != nil {
		return nil, err
	}
	return aq.ModelServer.GetAirQuality(ctx, request)
}

func (aq *airQualitySensor) PullAirQuality(request *traits.PullAirQualityRequest, server traits.AirQualitySensorApi_PullAirQualityServer) error {
	_ = aq.pollTask.Attach(server.Context())
	return aq.ModelServer.PullAirQuality(request, server)
}

// pollPeer fetches data from the peer device and saves the data locally.
func (aq *airQualitySensor) pollPeer(ctx context.Context) (*traits.AirQuality, error) {
	data := &traits.AirQuality{}
	var resProcessors []func(response any) error
	var readValues []config.ValueSource
	var requestNames []string

	if aq.config.AirPressure != nil {
		readValues = append(readValues, *aq.config.AirPressure)
		requestNames = append(requestNames, "AirPressure")
		resProcessors = append(resProcessors, func(response any) error {
			pressure, err := comm.Float32Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "air pressure", Cause: err}
			}
			data.AirPressure = &pressure
			return nil
		})
	}
	if aq.config.Co2 != nil {
		readValues = append(readValues, *aq.config.Co2)
		requestNames = append(requestNames, "Co2")
		resProcessors = append(resProcessors, func(response any) error {
			co2, err := comm.Float32Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "co2", Cause: err}
			}
			data.CarbonDioxideLevel = &co2
			return nil
		})
	}
	if aq.config.IaqScore != nil {
		readValues = append(readValues, *aq.config.IaqScore)
		requestNames = append(requestNames, "IaqScore")
		resProcessors = append(resProcessors, func(response any) error {
			iaq, err := comm.Float32Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "iaq", Cause: err}
			}
			data.Score = &iaq
			return nil
		})
	}
	if aq.config.Pm10 != nil {
		readValues = append(readValues, *aq.config.Pm10)
		requestNames = append(requestNames, "Pm10")
		resProcessors = append(resProcessors, func(response any) error {
			pm10, err := comm.Float32Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "pm10", Cause: err}
			}
			data.ParticulateMatter_10 = &pm10
			return nil
		})
	}
	if aq.config.Pm25 != nil {
		readValues = append(readValues, *aq.config.Pm25)
		requestNames = append(requestNames, "Pm25")
		resProcessors = append(resProcessors, func(response any) error {
			pm25, err := comm.Float32Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "pm25", Cause: err}
			}
			data.ParticulateMatter_25 = &pm25
			return nil
		})
	}
	if aq.config.VOC != nil {
		readValues = append(readValues, *aq.config.VOC)
		requestNames = append(requestNames, "VOC")
		resProcessors = append(resProcessors, func(response any) error {
			voc, err := comm.Float32Value(response)
			if err != nil {
				return comm.ErrReadProperty{Prop: "voc", Cause: err}
			}
			data.VolatileOrganicCompounds = &voc
			return nil
		})
	}
	responses := comm.ReadProperties(ctx, aq.client, aq.known, readValues...)
	var errs []error
	for i, response := range responses {
		err := resProcessors[i](response)
		if err != nil {
			errs = append(errs, err)
		}
	}
	status.UpdatePollErrorStatus(aq.statuses, aq.config.Name, "AirQualitySensor", requestNames, errs)
	if len(errs) > 0 {
		return nil, multierr.Combine(errs...)
	}
	return aq.model.UpdateAirQuality(data)
}
