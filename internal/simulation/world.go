package simulation

type World struct {
	Agents []*Agent `json:"agents"`
	Size   float64  `json:"size"`
	Tick   int      `json:"tick"`
}

func NewWorld(agentCount int, size float64) *World {
	agents := make([]*Agent, agentCount)
	for i := range agents {
		agents[i] = NewAgent(i, size)
	}
	return &World{Agents: agents, Size: size}
}

func (w *World) Update() {
	w.Tick++
	for _, a := range w.Agents {
		a.Update(1.0, w.Size)
	}
	for _, a := range w.Agents {
		a.CheckInteraction(w.Agents)
	}
}

func (w *World) Snapshot() []map[string]interface{} {
	snap := make([]map[string]interface{}, len(w.Agents))
	for i, a := range w.Agents {
		snap[i] = map[string]interface{}{
			"id":    a.ID,
			"x":     a.X,
			"y":     a.Y,
			"dx":    a.DX,
			"dy":    a.DY,
			"state": a.State,
		}
	}
	return snap
}
