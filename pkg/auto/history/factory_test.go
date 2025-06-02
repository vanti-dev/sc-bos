package history

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airqualitysensorpb"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperaturepb"
	"github.com/smart-core-os/sc-golang/pkg/trait/electricpb"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensorpb"
	"github.com/vanti-dev/sc-bos/pkg/auto/history/config"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	meterpb "github.com/vanti-dev/sc-bos/pkg/gentrait/meter"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
	"github.com/vanti-dev/sc-bos/pkg/node"
	"github.com/vanti-dev/sc-bos/pkg/util/jsontypes"
)

func Test_automation_applyConfig(t *testing.T) {
	ctx := context.Background()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	t.Cleanup(cancel)

	logger := zap.NewNop()
	occupancy := occupancysensorpb.NewModel()
	airQuality := airqualitysensorpb.NewModel()
	airTemperature := airtemperaturepb.NewModel()
	electric := electricpb.NewModel()
	meter := meterpb.NewModel()
	status := statuspb.NewModel()

	announcer := node.New("test")

	announcer.Logger = logger

	announcer.Announce("occupancy",
		node.HasTrait(trait.OccupancySensor),
		node.HasServer(
			traits.RegisterOccupancySensorApiServer,
			traits.OccupancySensorApiServer(occupancysensorpb.NewModelServer(occupancy)),
		),
	)

	announcer.Announce("airquality",
		node.HasTrait(trait.AirQualitySensor),
		node.HasServer(
			traits.RegisterAirQualitySensorApiServer,
			traits.AirQualitySensorApiServer(airqualitysensorpb.NewModelServer(airQuality)),
		),
	)

	announcer.Announce("airtemperature",
		node.HasTrait(trait.AirTemperature),
		node.HasServer(
			traits.RegisterAirTemperatureApiServer,
			traits.AirTemperatureApiServer(airtemperaturepb.NewModelServer(airTemperature)),
		),
	)

	announcer.Announce("electric",
		node.HasTrait(trait.Electric),
		node.HasServer(
			traits.RegisterElectricApiServer,
			traits.ElectricApiServer(electricpb.NewModelServer(electric)),
		),
	)

	announcer.Announce("meter",
		node.HasTrait(meterpb.TraitName),
		node.HasServer(
			gen.RegisterMeterApiServer,
			gen.MeterApiServer(meterpb.NewModelServer(meter)),
		),
	)

	announcer.Announce("status",
		node.HasTrait(statuspb.TraitName),
		node.HasServer(
			gen.RegisterStatusApiServer,
			gen.StatusApiServer(statuspb.NewModelServer(status)),
		),
	)

	for _, cfg := range cfgs {
		a := &automation{
			clients:   announcer,
			announcer: node.NewReplaceAnnouncer(announcer),
			logger:    logger,
		}

		err := a.applyConfig(ctx, cfg)

		if err != nil {
			t.Fatal(err)
		}
	}

	// many events to each model server
	for i := 0; i < 10; i++ {
		if _, err := occupancy.SetOccupancy(&traits.Occupancy{
			State:       traits.Occupancy_OCCUPIED,
			PeopleCount: int32(rand.Intn(10)),
		}); err != nil {
			t.Fatal(err)
		}
		if _, err := airQuality.UpdateAirQuality(&traits.AirQuality{
			CarbonDioxideLevel:       ptr(rand.Float32()),
			VolatileOrganicCompounds: ptr(rand.Float32()),
			AirPressure:              ptr(rand.Float32()),
			InfectionRisk:            ptr(rand.Float32()),
			Score:                    ptr(rand.Float32()),
			ParticulateMatter_1:      ptr(rand.Float32()),
			ParticulateMatter_25:     ptr(rand.Float32()),
			ParticulateMatter_10:     ptr(rand.Float32()),
			AirChangePerHour:         ptr(rand.Float32()),
		}); err != nil {
			t.Fatal(err)
		}

		if _, err := airTemperature.UpdateAirTemperature(&traits.AirTemperature{
			AmbientTemperature: &types.Temperature{ValueCelsius: rand.Float64()},
		}); err != nil {
			t.Fatal(err)
		}

		if _, err := electric.UpdateDemand(&traits.ElectricDemand{
			Voltage:       ptr(rand.Float32()),
			Current:       rand.Float32(),
			ReactivePower: ptr(rand.Float32()),
			ApparentPower: ptr(rand.Float32()),
			PowerFactor:   ptr(rand.Float32()),
			RealPower:     ptr(rand.Float32()),
		}); err != nil {
			t.Fatal(err)
		}

		if _, err := meter.UpdateMeterReading(&gen.MeterReading{
			Usage:    rand.Float32(),
			Produced: rand.Float32(),
		}); err != nil {
			t.Fatal(err)
		}

		if _, err := status.UpdateProblem(&gen.StatusLog_Problem{
			Name: randString(16),
		}); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
	}

	aqCli := gen.NewAirQualitySensorHistoryClient(announcer.ClientConn())
	occupancyCli := gen.NewOccupancySensorHistoryClient(announcer.ClientConn())
	airTempCli := gen.NewAirTemperatureHistoryClient(announcer.ClientConn())
	electricCli := gen.NewElectricHistoryClient(announcer.ClientConn())
	meterCli := gen.NewMeterHistoryClient(announcer.ClientConn())
	statusCli := gen.NewStatusHistoryClient(announcer.ClientConn())

	aqHist, err := aqCli.ListAirQualityHistory(ctx, &gen.ListAirQualityHistoryRequest{Name: "airquality", PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(int32(2), aqHist.GetTotalSize()); diff != "" {
		t.Fatal(diff, "airquality")
	}

	occHist, err := occupancyCli.ListOccupancyHistory(ctx, &gen.ListOccupancyHistoryRequest{Name: "occupancy", PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(int32(2), occHist.GetTotalSize()); diff != "" {
		t.Fatal(diff, "occupancy")
	}

	airTempHist, err := airTempCli.ListAirTemperatureHistory(ctx, &gen.ListAirTemperatureHistoryRequest{Name: "airtemperature", PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(int32(2), airTempHist.GetTotalSize()); diff != "" {
		t.Fatal(diff, "airtemperature")
	}

	electricHist, err := electricCli.ListElectricDemandHistory(ctx, &gen.ListElectricDemandHistoryRequest{Name: "electric", PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(int32(2), electricHist.GetTotalSize()); diff != "" {
		t.Fatal(diff, "electric")
	}

	meterHist, err := meterCli.ListMeterReadingHistory(ctx, &gen.ListMeterReadingHistoryRequest{Name: "meter", PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(int32(2), meterHist.GetTotalSize()); diff != "" {
		t.Fatal(diff, "meter")
	}

	statusHist, err := statusCli.ListCurrentStatusHistory(ctx, &gen.ListCurrentStatusHistoryRequest{Name: "status", PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(int32(2), statusHist.GetTotalSize()); diff != "" {
		t.Fatal(diff, "status")
	}

}

func ptr[T any](v T) *T {
	return &v
}

var cfgs = []config.Root{
	{
		Source: &config.Source{
			Name:            "occupancy",
			Trait:           trait.OccupancySensor,
			PollingSchedule: jsontypes.MustParseExtendedSchedule("*/5 * * * * *"),
		},
		Storage: &config.Storage{
			Type: "memory",
			TTL: &config.TTL{
				MaxAge:   jsontypes.Duration{Duration: time.Minute * 3},
				MaxCount: 10,
			},
		},
	},
	{
		Source: &config.Source{
			Name:            "airquality",
			Trait:           trait.AirQualitySensor,
			PollingSchedule: jsontypes.MustParseExtendedSchedule("*/5 * * * * *"),
		},
		Storage: &config.Storage{
			Type: "memory",
			TTL: &config.TTL{
				MaxAge:   jsontypes.Duration{Duration: time.Minute * 3},
				MaxCount: 10,
			},
		},
	},
	{
		Source: &config.Source{
			Name:            "airtemperature",
			Trait:           trait.AirTemperature,
			PollingSchedule: jsontypes.MustParseExtendedSchedule("*/5 * * * * *"),
		},
		Storage: &config.Storage{
			Type: "memory",
			TTL: &config.TTL{
				MaxAge:   jsontypes.Duration{Duration: time.Minute * 3},
				MaxCount: 10,
			},
		},
	},
	{
		Source: &config.Source{
			Name:            "electric",
			Trait:           trait.Electric,
			PollingSchedule: jsontypes.MustParseExtendedSchedule("*/5 * * * * *"),
		},
		Storage: &config.Storage{
			Type: "memory",
			TTL: &config.TTL{
				MaxAge:   jsontypes.Duration{Duration: time.Minute * 3},
				MaxCount: 10,
			},
		},
	},
	{
		Source: &config.Source{
			Name:            "meter",
			Trait:           meterpb.TraitName,
			PollingSchedule: jsontypes.MustParseExtendedSchedule("*/5 * * * * *"),
		},
		Storage: &config.Storage{
			Type: "memory",
			TTL: &config.TTL{
				MaxAge:   jsontypes.Duration{Duration: time.Minute * 3},
				MaxCount: 10,
			},
		},
	},
	{
		Source: &config.Source{
			Name:            "status",
			Trait:           statuspb.TraitName,
			PollingSchedule: jsontypes.MustParseExtendedSchedule("*/5 * * * * *"),
		},
		Storage: &config.Storage{
			Type: "memory",
			TTL: &config.TTL{
				MaxAge:   jsontypes.Duration{Duration: time.Minute * 3},
				MaxCount: 10,
			},
		},
	},
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
