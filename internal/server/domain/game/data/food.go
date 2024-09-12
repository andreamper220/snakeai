package data

import "math/rand"

type Food struct {
	Position Point
}

func NewFood(width, height int) *Food {
	return &Food{
		Position: Point{
			X: rand.Intn(width) + 1,
			Y: rand.Intn(height) + 1,
		},
	}
}
