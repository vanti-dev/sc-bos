package xovis

import (
	"time"
)

type PushData struct {
	LogicsData *LogicsPushData `json:"logics_data,omitempty"`
}

type LogicsPushData struct {
	PackageInfo *PackageInfo `json:"package_info,omitempty"`
	SensorInfo  *SensorInfo  `json:"sensor_info,omitempty"`
	Logics      []Logic      `json:"logics"`
}

type Logic struct {
	ID         int           `json:"id"`
	Name       string        `json:"name"`
	Info       string        `json:"info"`
	Geometries []Geometry    `json:"geometries"`
	Records    []LogicRecord `json:"records"`
}

type LogicRecord struct {
	From            time.Time `json:"from"`
	To              time.Time `json:"to"`
	Samples         int       `json:"samples"`
	SamplesExpected int       `json:"samples_expected"`
	Counts          []Count   `json:"counts"`
}

type PackageInfo struct {
	Version   string    `json:"version"`
	ID        int       `json:"id"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	AgentID   int       `json:"agent_id"`
	AgentName string    `json:"agent_name"`
	AgentType string    `json:"agent_type"`
}

type SensorInfo struct {
	Type         string    `json:"type"`
	SerialNumber string    `json:"serial_number"`
	Name         string    `json:"name"`
	Group        string    `json:"group"`
	DeviceType   string    `json:"device_type"`
	HWVersion    string    `json:"hw_version"`
	SWVersion    string    `json:"sw_version"`
	Time         time.Time `json:"time"`
	Timezone     string    `json:"timezone"`
}
