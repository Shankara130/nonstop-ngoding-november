package simulation

type ZoneType string

const (
	ZoneTypeHome   ZoneType = "home"
	ZoneTypeFood   ZoneType = "food"
	ZoneTypeWater  ZoneType = "water"
	ZoneTypeDanger ZoneType = "danger"
)

type Zone struct {
	ID     int
	X, Y   float64
	Radius float64
	Type   ZoneType
}

type ZoneSnapshot struct {
	ID     int      `json:"id"`
	X      float64  `json:"x"`
	Y      float64  `json:"y"`
	Radius float64  `json:"radius"`
	Type   ZoneType `json:"type"`
}
