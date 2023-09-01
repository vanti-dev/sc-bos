package hpd3

const (
	PointPresence1           = "Presence1"
	PointMotion1             = "Motion1"
	PointTemperature         = "Temperature"
	PointHumidity            = "Humidity"
	PointNumberOfPeopleTotal = "NumberOfPeopleTotal"
	PointDewPoint            = "DewPoint"
)

type PointData struct {
	Presence1           bool    `json:"Presence1"`
	Motion1             bool    `json:"Motion1"`
	Temperature         float64 `json:"Temperature"` // degrees C
	Humidity            float64 `json:"Humidity"`    // percent relative humidity
	NumberOfPeopleTotal int     `json:"NumberOfPeopleTotal"`
	DewPoint            float64 `json:"DewPoint"` // degrees C
}
