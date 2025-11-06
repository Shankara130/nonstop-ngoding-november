package simulation

type ResourceType int

const (
	ResourceTypeFood ResourceType = iota
	ResourceTypeWater
)

type Resource struct {
	ID     int
	X, Y   float64
	Amount float64
	Type   ResourceType
}

type ResourceSnapshot struct {
	ID     int          `json:"id"`
	X      float64      `json:"x"`
	Y      float64      `json:"y"`
	Amount float64      `json:"amount"`
	Type   ResourceType `json:"type"`
}

func (r *Resource) Snapshot() ResourceSnapshot {
	return ResourceSnapshot{
		ID:     r.ID,
		X:      r.X,
		Y:      r.Y,
		Amount: r.Amount,
		Type:   r.Type,
	}
}
