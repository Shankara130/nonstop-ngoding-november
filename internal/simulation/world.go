package simulation

import (
	"fmt"
	"math/rand"
	"sync"
)

type World struct {
	Width, Height float64
	Tick          int
	Agents        []*Agent
	mu            sync.Mutex
}

func NewWorld(width, height float64, numAgents int) *World {
	world := &World{
		Width:  width,
		Height: height,
		Agents: make([]*Agent, 0, numAgents),
	}

	for i := 0; i < numAgents; i++ {
		a := &Agent{
			ID:   i,
			Name: fmt.Sprintf("Agent-%d", i),
			X:    randFloat(0, width),
			Y:    randFloat(0, height),
		}
		world.Agents = append(world.Agents, a)
	}

	return world
}

func (w *World) Update() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, agent := range w.Agents {
		agent.Move(w.Width, w.Height)
	}

	w.Tick++
}

func (w *World) Snapshot() []map[string]interface{} {
	w.mu.Lock()
	defer w.mu.Unlock()

	result := []map[string]interface{}{}
	for _, a := range w.Agents {
		result = append(result, map[string]interface{}{
			"id":   a.ID,
			"name": a.Name,
			"x":    a.X,
			"y":    a.Y,
		})
	}
	return result
}

func randFloat(min, max float64) float64 {
	return min + (max-min)*float64(rand.Intn(1000))/1000.0
}
