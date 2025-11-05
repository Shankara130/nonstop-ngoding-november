package simulation

import (
	"math"
	"math/rand"
	"time"
)

type AgentType string

const (
	AgentTypeWorker    AgentType = "worker"
	AgentTypeExplorer  AgentType = "explorer"
	AgentTypeCollector AgentType = "collector"
	AgentTypeGuard     AgentType = "guard"
)

type AgentState string

const (
	StateIdle       AgentState = "idle"
	StateMoving     AgentState = "moving"
	StateCollecting AgentState = "collecting"
	StateResting    AgentState = "resting"
	StatePatrolling AgentState = "patrolling"
)

type Agent struct {
	ID        int
	Type      AgentType
	X, Y      float64
	VX, VY    float64
	Energy    float64
	MaxEnergy float64
	Speed     float64
	State     AgentState
	Target    interface{}
	Inventory float64
	alive     bool
}

type AgentSnapshot struct {
	ID        int        `json:"id"`
	Type      AgentType  `json:"type"`
	X         float64    `json:"x"`
	Y         float64    `json:"y"`
	Energy    float64    `json:"energy"`
	MaxEnergy float64    `json:"maxEnergy"`
	State     AgentState `json:"state"`
	Inventory float64    `json:"inventory"`
}

func NewAgent(id int, worldSize float64, agentType AgentType) *Agent {
	a := &Agent{
		ID:        id,
		Type:      agentType,
		X:         rand.Float64() * worldSize,
		Y:         rand.Float64() * worldSize,
		MaxEnergy: 100,
		Energy:    100,
		State:     StateIdle,
		alive:     true,
	}

	// Set attributes based on agent type
	switch agentType {
	case AgentTypeWorker:
		a.Speed = 8.0
		a.MaxEnergy = 100
	case AgentTypeExplorer:
		a.Speed = 12.0
		a.MaxEnergy = 150
	case AgentTypeCollector:
		a.Speed = 6.0
		a.MaxEnergy = 120
	case AgentTypeGuard:
		a.Speed = 10.0
		a.MaxEnergy = 200
	}

	a.Energy = a.MaxEnergy
	return a
}

func (a *Agent) Run(world *World) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for a.alive {
		<-ticker.C
		a.Update(world, 0.1)

		if a.Energy <= 0 {
			a.alive = false
			world.RemoveAgent(a.ID)
			break
		}
	}
}

func (a *Agent) Update(world *World, dt float64) {
	// consume energy
	energyCost := 0.05
	if a.State == StateMoving {
		energyCost = 0.15
	}
	a.Energy -= energyCost

	// behavior based on agent type
	switch a.Type {
	case AgentTypeWorker:
		a.workerBehavior(world, dt)
	case AgentTypeExplorer:
		a.explorerBehavior(world, dt)
	case AgentTypeCollector:
		a.collectorBehavior(world, dt)
	case AgentTypeGuard:
		a.guardBehavior(world, dt)
	}

	// boundary check

	a.X = math.Max(0, math.Min(world.Size, a.X))
	a.Y = math.Max(0, math.Min(world.Size, a.Y))
}

func (a *Agent) workerBehavior(world *World, dt float64) {
	// worker : seek resources and bring them to home base
	if a.Energy < 30 {
		a.seekRest(world, dt)
		return
	}

	if a.Inventory < 10 {
		// find food resource
		target := world.FindNearestResource(a.X, a.Y, ResourceTypeFood)
		if target != nil {
			a.moveTo(target.X, target.Y, dt)
			if a.reachedPoint(target.X, target.Y, 3) {
				a.State = StateCollecting
				if world.ConsumeResource(target.ID, 5) {
					a.Inventory += 5
				}
			}
		} else {
			a.wander(dt)
		}
	} else {
		// return to home base
		home := world.FindNearestZone(a.X, a.Y, ZoneTypeHome)
		if home != nil {
			a.moveTo(home.X, home.Y, dt)
			if a.reachedPoint(home.X, home.Y, home.Radius) {
				a.Inventory = 0
				a.Energy = math.Min(a.MaxEnergy, a.Energy+10)
			}
		}
	}
}

func (a *Agent) explorerBehavior(world *World, dt float64) {
	// explorer : wander map, find new resources
	if a.Energy < 40 {
		a.seekRest(world, dt)
		return
	}

	a.wander(dt)
	a.Speed = 12.0
}

func (a *Agent) collectorBehavior(world *World, dt float64) {
	// collector : gather water resources
	if a.Energy < 25 {
		a.seekRest(world, dt)
		return
	}

	target := world.FindNearestResource(a.X, a.Y, ResourceTypeWater)
	if target != nil {
		a.moveTo(target.X, target.Y, dt)
		if a.reachedPoint(target.X, target.Y, 3) {
			a.State = StateCollecting
			if world.ConsumeResource(target.ID, 3) {
				a.Energy = math.Min(a.MaxEnergy, a.Energy+5)
			}
		}
	} else {
		a.wander(dt)
	}
}

func (a *Agent) guardBehavior(world *World, dt float64) {
	// guard : patrol around home base
	home := world.FindNearestZone(a.X, a.Y, ZoneTypeHome)
	if home != nil {
		// patrol around home
		dist := distance(a.X, a.Y, home.X, home.Y)
		if dist > home.Radius*2 {
			a.moveTo(home.X, home.Y, dt)
		} else {
			a.State = StatePatrolling
			a.wander(dt)
		}
	}
}

func (a *Agent) seekRest(world *World, dt float64) {
	a.State = StateResting
	home := world.FindNearestZone(a.X, a.Y, ZoneTypeHome)
	if home != nil {
		a.moveTo(home.X, home.Y, dt)
		if a.reachedPoint(home.X, home.Y, home.Radius) {
			a.Energy = math.Min(a.MaxEnergy, a.Energy+1)
		}
	}
}

func (a *Agent) moveTo(tx, ty, dt float64) {
	dx, dy := tx-a.X, ty-a.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist < 0.1 {
		return
	}

	a.VX = (dx / dist) * a.Speed
	a.VY = (dy / dist) * a.Speed
	a.X += a.VX * dt
	a.Y += a.VY * dt
	a.State = StateMoving
}

func (a *Agent) wander(dt float64) {
	if rand.Float64() < 0.05 {
		angle := rand.Float64() * 2 * math.Pi
		a.VX = math.Cos(angle) * a.Speed
		a.VY = math.Sin(angle) * a.Speed
	}
	a.X += a.VX * dt
	a.Y += a.VY * dt
	a.State = StateMoving
}

func (a *Agent) reachedPoint(x, y, threshold float64) bool {
	return distance(a.X, a.Y, x, y) < threshold
}

func (a *Agent) Snapshot() AgentSnapshot {
	return AgentSnapshot{
		ID:        a.ID,
		Type:      a.Type,
		X:         a.X,
		Y:         a.Y,
		Energy:    a.Energy,
		MaxEnergy: a.MaxEnergy,
		State:     a.State,
		Inventory: a.Inventory,
	}
}
