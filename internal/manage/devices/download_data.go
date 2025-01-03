package devices

import (
	"context"
	"fmt"

	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/vanti-dev/sc-bos/pkg/gen"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/accesspb"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/meter"
	"github.com/vanti-dev/sc-bos/pkg/gentrait/statuspb"
)

// getTraitInfo returns the data we support exporting via the DevicesApi download endpoint.
func (s *Server) getTraitInfo() map[string]traitInfo {
	// Note: I want this to be dynamic, however, there are a few issues preventing this:
	// 1. The protoregistry.Files and Types does not include gateway discovered types so we can't use that
	// 2. Not all data for traits is useful, so we need to filter out the useful data
	// 3. Property names aren't that useful to users so some transformation is needed
	// 4. Not all traits include the relevant info in their Api aspect. Check the meter one for where the unit comes from.

	return map[string]traitInfo{
		string(accesspb.TraitName): {
			headers: []string{"access.grant", "access.reason", "access.actor.name", "access.actor.title", "access.actor.displayname", "access.actor.email"},
			get: func(ctx context.Context, name string) (map[string]string, error) {
				vals := make(map[string]string)
				c := gen.NewAccessApiClient(s.node.ClientConn())
				data, err := c.GetLastAccessAttempt(ctx, &gen.GetLastAccessAttemptRequest{Name: name})
				if err != nil {
					return nil, err
				}
				vals["access.grant"] = data.GetGrant().String()
				vals["access.reason"] = data.GetReason()
				if actor := data.GetActor(); actor != nil {
					vals["access.actor.name"] = actor.Name
					vals["access.actor.title"] = actor.Title
					vals["access.actor.displayname"] = actor.DisplayName
					vals["access.actor.email"] = actor.Email
				}
				return vals, nil
			},
		},
		string(meter.TraitName): {
			headers: []string{"meter.usage", "meter.unit"},
			get: func(ctx context.Context, name string) (map[string]string, error) {
				vals := make(map[string]string)
				c := gen.NewMeterApiClient(s.node.ClientConn())
				data, err := c.GetMeterReading(ctx, &gen.GetMeterReadingRequest{Name: name})
				if err != nil {
					return nil, err
				}
				vals["meter.usage"] = fmt.Sprintf("%.3f", data.Usage)

				ci := gen.NewMeterInfoClient(s.node.ClientConn())
				info, err := ci.DescribeMeterReading(ctx, &gen.DescribeMeterReadingRequest{Name: name})
				if err == nil {
					vals["meter.unit"] = info.Unit
				}
				return vals, nil
			},
		},
		string(statuspb.TraitName): {
			headers: []string{"status.level", "status.description", "status.recordtime"},
			get: func(ctx context.Context, name string) (map[string]string, error) {
				vals := make(map[string]string)
				c := gen.NewStatusApiClient(s.node.ClientConn())
				data, err := c.GetCurrentStatus(ctx, &gen.GetCurrentStatusRequest{Name: name})
				if err != nil {
					return nil, err
				}
				vals["status.level"] = data.GetLevel().String()
				vals["status.description"] = data.GetDescription()
				if data.RecordTime != nil {
					vals["status.recordtime"] = data.GetRecordTime().AsTime().String()
				}
				return vals, nil
			},
		},
		string(trait.AirQualitySensor): {
			headers: []string{"iaq.co2", "iaq.voc", "iaq.pressure", "iaq.comfort", "iaq.infectionrisk", "iaq.score", "iaq.pm1", "iaq.pm25", "iaq.pm10", "iaq.airchange"},
			get: func(ctx context.Context, name string) (map[string]string, error) {
				vals := make(map[string]string)
				c := traits.NewAirQualitySensorApiClient(s.node.ClientConn())
				data, err := c.GetAirQuality(ctx, &traits.GetAirQualityRequest{Name: name})
				if err != nil {
					return nil, err
				}
				if data.CarbonDioxideLevel != nil {
					vals["iaq.co2"] = fmt.Sprintf("%.3f", *data.CarbonDioxideLevel)
				}
				if data.VolatileOrganicCompounds != nil {
					vals["iaq.voc"] = fmt.Sprintf("%.3f", *data.VolatileOrganicCompounds)
				}
				if data.AirPressure != nil {
					vals["iaq.pressure"] = fmt.Sprintf("%.3f", *data.AirPressure)
				}
				if data.Comfort != traits.AirQuality_COMFORT_UNSPECIFIED {
					vals["iaq.comfort"] = data.Comfort.String()
				}
				if data.InfectionRisk != nil {
					vals["iaq.infectionrisk"] = fmt.Sprintf("%.3f", *data.InfectionRisk)
				}
				if data.Score != nil {
					vals["iaq.score"] = fmt.Sprintf("%.3f", *data.Score)
				}
				if data.ParticulateMatter_1 != nil {
					vals["iaq.pm1"] = fmt.Sprintf("%.3f", *data.ParticulateMatter_1)
				}
				if data.ParticulateMatter_25 != nil {
					vals["iaq.pm25"] = fmt.Sprintf("%.3f", *data.ParticulateMatter_25)
				}
				if data.ParticulateMatter_10 != nil {
					vals["iaq.pm10"] = fmt.Sprintf("%.3f", *data.ParticulateMatter_10)
				}
				if data.AirChangePerHour != nil {
					vals["iaq.airchange"] = fmt.Sprintf("%.3f", *data.AirChangePerHour)
				}
				return vals, nil
			},
		},
		string(trait.AirTemperature): {
			headers: []string{"airtemperature.temperature", "airtemperature.humidity", "airtemperature.setpoint"},
			get: func(ctx context.Context, name string) (map[string]string, error) {
				vals := make(map[string]string)
				c := traits.NewAirTemperatureApiClient(s.node.ClientConn())
				data, err := c.GetAirTemperature(ctx, &traits.GetAirTemperatureRequest{Name: name})
				if err != nil {
					return nil, err
				}
				if data.AmbientTemperature != nil {
					vals["airtemperature.temperature"] = fmt.Sprintf("%.1f", data.AmbientTemperature.ValueCelsius)
				}
				if data.AmbientHumidity != nil {
					vals["airtemperature.humidity"] = fmt.Sprintf("%.1f", *data.AmbientHumidity)
				}
				if data.GetTemperatureSetPoint() != nil {
					vals["airtemperature.setpoint"] = fmt.Sprintf("%.1f", data.GetTemperatureSetPoint().ValueCelsius)
				}
				return vals, nil
			},
		},
		string(trait.Electric): {
			headers: []string{"electric.current", "electric.voltage", "electric.powerfactor", "electric.realpower", "electric.apparentpower", "electric.reactivepower", "electric.reactivepower"},
			get: func(ctx context.Context, name string) (map[string]string, error) {
				vals := make(map[string]string)
				c := traits.NewElectricApiClient(s.node.ClientConn())
				data, err := c.GetDemand(ctx, &traits.GetDemandRequest{Name: name})
				if err != nil {
					return nil, err
				}
				if data.Current != 0 {
					vals["electric.current"] = fmt.Sprintf("%.3f", data.Current)
				}
				if data.Voltage != nil {
					vals["electric.voltage"] = fmt.Sprintf("%.3f", *data.Voltage)
				}
				if data.PowerFactor != nil {
					vals["electric.powerfactor"] = fmt.Sprintf("%.3f", *data.PowerFactor)
				}
				if data.RealPower != nil {
					vals["electric.realpower"] = fmt.Sprintf("%.3f", *data.RealPower)
				}
				if data.ApparentPower != nil {
					vals["electric.apparentpower"] = fmt.Sprintf("%.3f", *data.ApparentPower)
				}
				if data.ReactivePower != nil {
					vals["electric.reactivepower"] = fmt.Sprintf("%.3f", *data.ReactivePower)
				}
				return vals, nil
			},
		},
		string(trait.EnterLeaveSensor): {
			headers: []string{"enterleave.entertotal", "enterleave.leavetotal"},
			get: func(ctx context.Context, name string) (map[string]string, error) {
				vals := make(map[string]string)
				c := traits.NewEnterLeaveSensorApiClient(s.node.ClientConn())
				data, err := c.GetEnterLeaveEvent(ctx, &traits.GetEnterLeaveEventRequest{Name: name})
				if err != nil {
					return nil, err
				}
				vals["enterleave.entertotal"] = fmt.Sprintf("%d", data.EnterTotal)
				vals["enterleave.leavetotal"] = fmt.Sprintf("%d", data.LeaveTotal)
				return vals, nil
			},
		},
		string(trait.FanSpeed): {
			headers: []string{"fanspeed.percentage", "fanspeed.preset"},
			get: func(ctx context.Context, name string) (map[string]string, error) {
				vals := make(map[string]string)
				c := traits.NewFanSpeedApiClient(s.node.ClientConn())
				data, err := c.GetFanSpeed(ctx, &traits.GetFanSpeedRequest{Name: name})
				if err != nil {
					return nil, err
				}
				vals["fanspeed.percentage"] = fmt.Sprintf("%.1f", data.Percentage)
				vals["fanspeed.preset"] = data.Preset
				return vals, nil
			},
		},
		string(trait.Light): {
			headers: []string{"light.level", "light.preset"},
			get: func(ctx context.Context, name string) (map[string]string, error) {
				vals := make(map[string]string)
				c := traits.NewLightApiClient(s.node.ClientConn())
				data, err := c.GetBrightness(ctx, &traits.GetBrightnessRequest{Name: name})
				if err != nil {
					return nil, err
				}
				vals["light.level"] = fmt.Sprintf("%.1f", data.LevelPercent)
				if data.Preset != nil {
					vals["light.preset"] = data.Preset.Name
				}
				return vals, nil
			},
		},
		string(trait.OccupancySensor): {
			headers: []string{"occupancy.state", "occupancy.peoplecount", "occupancy.changetime"},
			get: func(ctx context.Context, name string) (map[string]string, error) {
				vals := make(map[string]string)
				c := traits.NewOccupancySensorApiClient(s.node.ClientConn())
				data, err := c.GetOccupancy(ctx, &traits.GetOccupancyRequest{Name: name})
				if err != nil {
					return nil, err
				}
				vals["occupancy.state"] = data.State.String()
				vals["occupancy.peoplecount"] = fmt.Sprintf("%d", data.PeopleCount)
				if data.StateChangeTime != nil {
					vals["occupancy.changetime"] = data.StateChangeTime.AsTime().String()
				}
				return vals, nil
			},
		},
		string(trait.OpenClose): {
			headers: []string{"openclose.openpercent"},
			get: func(ctx context.Context, name string) (map[string]string, error) {
				vals := make(map[string]string)
				c := traits.NewOpenCloseApiClient(s.node.ClientConn())
				data, err := c.GetPositions(ctx, &traits.GetOpenClosePositionsRequest{Name: name})
				if err != nil {
					return nil, err
				}
				if len(data.States) != 1 {
					return vals, nil
				}
				pos := data.States[0]
				vals["openclose.openpercent"] = fmt.Sprintf("%.1f", pos.OpenPercent)
				return vals, nil
			},
		},
	}
}
