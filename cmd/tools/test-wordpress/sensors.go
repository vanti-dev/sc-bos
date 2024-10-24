package main

import (
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-api/go/types"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"github.com/smart-core-os/sc-golang/pkg/trait/airqualitysensor"
	"github.com/smart-core-os/sc-golang/pkg/trait/airtemperature"
	"github.com/smart-core-os/sc-golang/pkg/trait/occupancysensor"

	"github.com/vanti-dev/sc-bos/pkg/node"
)

func announceOccupancy(root node.Announcer, name string, val int32) error {
	model := occupancysensor.NewModel()
	_, err := model.SetOccupancy(&traits.Occupancy{PeopleCount: val})
	if err != nil {
		return err
	}
	client := node.WithClients(occupancysensor.WrapApi(occupancysensor.NewModelServer(model)))
	root.Announce(name, node.HasTrait(trait.OccupancySensor, client))
	return nil
}

func announceTemperature(root node.Announcer, name string, celsius float64) error {
	model := airtemperature.NewModel()
	_, err := model.UpdateAirTemperature(&traits.AirTemperature{AmbientTemperature: &types.Temperature{ValueCelsius: celsius}})
	if err != nil {
		return err
	}
	client := node.WithClients(airtemperature.WrapApi(airtemperature.NewModelServer(model)))
	root.Announce(name, node.HasTrait(trait.AirTemperature, client))
	return nil
}

func announceAirQuality(root node.Announcer, name string, val float32) error {
	model := airqualitysensor.NewModel()
	_, err := model.UpdateAirQuality(&traits.AirQuality{CarbonDioxideLevel: &val})
	if err != nil {
		return err
	}
	client := node.WithClients(airqualitysensor.WrapApi(airqualitysensor.NewModelServer(model)))
	root.Announce(name, node.HasTrait(trait.AirQualitySensor, client))
	return nil
}
