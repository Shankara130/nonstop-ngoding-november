package simulation

import (
	"math/rand"
)

type Agent struct {
	ID   int
	Name string
	X, Y float64
}

func (a *Agent) Move(boundsX, boundsY float64) {
	dx := (rand.Float64() - 0.5) * 2
	dy := (rand.Float64() - 0.5) * 2

	a.X += dx
	a.Y += dy

	if a.X < 0 {
		a.X = 0
	} else if a.X > boundsX {
		a.X = boundsX
	}

	if a.Y < 0 {
		a.Y = 0
	} else if a.Y > boundsY {
		a.Y = boundsY
	}
}
