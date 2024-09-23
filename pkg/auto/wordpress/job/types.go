package job

import (
	"time"
)

// Meta shared metadata
type Meta struct {
	Timestamp time.Time `json:"timestamp"`
	Site      string    `json:"siteIdentifier"`
}

type IntMeasure struct {
	Value int32  `json:"value"`
	Units string `json:"units"`
}

type Float32Measure struct {
	Value float32 `json:"value"`
	Units string  `json:"units"`
}

type Float64Measure struct {
	Value float64 `json:"value"`
	Units string  `json:"units"`
}

type TotalOccupancy struct {
	Meta
	TotalOccupancy IntMeasure `json:"totalOccupancy"`
}

type AverageTemperature struct {
	Meta
	AverageTemperature Float64Measure `json:"averageTemperature"`
}

type EnergyConsumption struct {
	Meta
	TodaysEnergyConsumption Float32Measure `json:"todaysEnergyConsumption"`
}

type AverageCo2 struct {
	Meta
	AverageCo2 Float32Measure `json:"averageCo2"`
}

type WaterConsumption struct {
	Meta
	TodaysWaterConsumption Float32Measure `json:"todaysWaterConsumption"`
}
