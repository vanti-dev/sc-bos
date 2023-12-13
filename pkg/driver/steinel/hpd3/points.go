package hpd3

import (
	"encoding/json"
)

const (
	PointPresence1           = "Presence1"
	PointMotion1             = "Motion1"
	PointTemperature         = "Temperature"
	PointHumidity            = "Humidity"
	PointNumberOfPeopleTotal = "NumberOfPeopleTotal"
	PointCO2                 = "CO2"
	PointVOC                 = "VOC"
)

type PointData struct {
	Presence1           bool    `json:"Presence1"`
	Motion1             bool    `json:"Motion1"`
	Temperature         float64 `json:"Temperature"` // degrees C
	Humidity            float64 `json:"Humidity"`    // percent relative humidity
	NumberOfPeopleTotal int     `json:"NumberOfPeopleTotal"`
	CO2                 float64 `json:"CO2"` // ppm
	VOC                 float64 `json:"VOC"` // ppb
}

func (d *PointData) AsMap() map[string]any {
	// none of the data types should return errors when (de)serialising so nothing here should error
	// if it does, it must be a bug in this function

	asJson, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}

	var asMap map[string]any
	err = json.Unmarshal(asJson, &asMap)
	if err != nil {
		panic(err)
	}

	return asMap
}
