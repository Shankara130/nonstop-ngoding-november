package simulation

import (
	"math/rand"
	"time"
)

type State string

const (
	Idle     State = "idle"
	Moving   State = "moving"
	Interact State = "interact"
)

type Agent struct {
	ID    int     `json:"id"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	DX    float64 `json:"dx"`
	DY    float64 `json:"dy"`
	State State   `json:"state"`
}

func (a *Agent) Update(dt float64, worldSize float64) {
	if a.State == Interact {
		return
	}

	a.X += a.DX * dt
	a.Y += a.DY * dt

	if a.X < 0 || a.X > worldSize {
		a.DX *= -1
	}

	if a.Y < 0 || a.Y > worldSize {
		a.DY *= -1
	}

	if rand.Float64() < 0.01 {
		a.DX = rand.Float64()*2 - 1
		a.DY = rand.Float64()*2 - 1
		a.State = Moving
	}
}

func (a *Agent) CheckInteraction(other []*Agent) {
	for _, b := range other {
		if a.ID == b.ID {
			continue
		}
		dx := a.X - b.X
		dy := a.Y - b.Y
		dist := dx*dx + dy*dy
		if dist < 25 {
			a.State = Interact
			b.State = Interact
			return
		}
	}
	a.State = Moving
}

func NewAgent(id int, worldSize float64) *Agent {
	rand.Seed(time.Now().UnixNano())
	return &Agent{
		ID:    id,
		X:     rand.Float64() * worldSize,
		Y:     rand.Float64() * worldSize,
		DX:    rand.Float64()*2 - 1,
		DY:    rand.Float64()*2 - 1,
		State: Idle,
	}
}
