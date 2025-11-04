package simulation

type ZoneType string

const (
	ZoneMeet    ZoneType = "meet"
	ZoneCollect ZoneType = "collect"
	ZoneRest    ZoneType = "rest"
)

type Zone struct {
	ID     int      `json:"id"`
	X      float64  `json:"x"`
	Y      float64  `json:"y"`
	Radius float64  `json:"radius"`
	Type   ZoneType `json:"type"`
}

func NewZone(id int, x, y, radius float64, t ZoneType) *Zone {
	return &Zone{
		ID: id, X: x, Y: y, Radius: radius, Type: t,
	}
}
