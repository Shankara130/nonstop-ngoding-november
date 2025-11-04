package simulation

import (
	"math"
	"math/rand"
	"time"
)

type Goal string
type State string

const (
	GoalExplore Goal = "explore"
	GoalMeet    Goal = "meet"
	GoalCollect Goal = "collect"
	GoalRest    Goal = "rest"
)

const (
	Idle     State = "idle"
	Moving   State = "moving"
	Interact State = "interact"
)

type Agent struct {
	ID     int     `json:"id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	DX     float64 `json:"dx"`
	DY     float64 `json:"dy"`
	Energy float64 `json:"energy"`
	Goal   Goal    `json:"goal"`
	State  State   `json:"state"`
	Target *Zone   `json:"target,omitempty"`
}

func (a *Agent) distanceTo(x, y float64) float64 {
	dx := a.X - x
	dy := a.Y - y
	return math.Sqrt(dx*dx + dy*dy)
}

func (a *Agent) moveTo(target *Zone, dt float64) {
	if target == nil {
		return
	}
	dx := target.X - a.X
	dy := target.Y - a.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist == 0 {
		return
	}
	speed := 10.0
	step := speed * dt
	a.DX = dx / dist * step
	a.DY = dy / dist * step
	a.X += a.DX
	a.Y += a.DY
	a.State = Moving
}

func (a *Agent) wander() {
	a.X += a.DX
	a.Y += a.DY
	if rand.Float64() < 0.02 {
		a.DX = rand.Float64()*2 - 1
		a.DY = rand.Float64()*2 - 1
	}
}

func (a *Agent) consumeEnergy() {
	switch a.State {
	case Moving:
		a.Energy -= 0.5
	case Interact:
		a.Energy -= 0.2
	default:
		a.Energy -= 0.1
	}
	if a.Energy < 0 {
		a.Energy = 0
	}
}

func (a *Agent) reached(z *Zone) bool {
	if z == nil {
		return false
	}
	return a.distanceTo(z.X, z.Y) < z.Radius
}

func (a *Agent) Update(world *World) {
	if a.State == Interact && rand.Float64() < 0.01 {
		a.State = Idle
		a.Goal = GoalExplore
	}

	a.consumeEnergy()

	if a.Energy <= 10 {
		a.Goal = GoalRest
		a.Target = world.FindNearestZone(a.X, a.Y, ZoneRest)
	}

	switch a.Goal {
	case GoalExplore:
		a.wander()
		if rand.Float64() < 0.005 {
			a.Goal = GoalMeet
			a.Target = world.FindNearestZone(a.X, a.Y, ZoneMeet)
		}
	case GoalMeet, GoalCollect, GoalRest:
		a.moveTo(a.Target)
	}

	if a.reached(a.Target) {
		switch a.Target.Type {
		case ZoneRest:
			a.Energy = 100
			a.Goal = GoalExplore
		case ZoneCollect:
			a.Energy += 10
			if a.Energy > 100 {
				a.Energy = 100
			}
		case ZoneMeet:
			a.State = Interact
		}
	}

	if a.X < 0 {
		a.X = 0
		a.DX *= -1
	}
	if a.Y < 0 {
		a.Y = 0
		a.DY *= -1
	}
	if a.X > world.Size {
		a.X = world.Size
		a.DX *= -1
	}
	if a.Y > world.Size {
		a.Y = world.Size
		a.DY *= -1
	}
}

func NewAgent(id int, worldSize float64) *Agent {
	return &Agent{
		ID:     id,
		X:      rand.Float64() * worldSize,
		Y:      rand.Float64() * worldSize,
		DX:     rand.Float64()*2 - 1,
		DY:     rand.Float64()*2 - 1,
		Energy: 100,
		Goal:   GoalExplore,
		State:  Idle,
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
