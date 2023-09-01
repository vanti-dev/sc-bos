package hpd3

const (
	PointPresence1           = "Presence1"
	PointMotion1             = "Motion1"
	PointTemperature         = "Temperature"
	PointNumberOfPeopleTotal = "NumberOfPeopleTotal"
)

type PointData struct {
	Presence1           bool    `json:"Presence1"`
	Motion1             bool    `json:"Motion1"`
	Temperature         float64 `json:"Temperature"`
	NumberOfPeopleTotal int     `json:"NumberOfPeopleTotal"`
}
