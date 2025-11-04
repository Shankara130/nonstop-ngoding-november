package simulation

import (
	"log"
	"time"
)

func SimulationLoop(world *World, tickRate time.Duration, onUpdate func(map[string]interface{})) {
	ticker := time.NewTicker(tickRate)
	defer ticker.Stop()

	for range ticker.C {
		world.Update(dt)
		state := world.Snapshot()
		onUpdate(state)
		log.Printf("Tick %d selesai", world.Tick)
	}
}
