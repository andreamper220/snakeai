package data

import "math/rand"

type Food struct {
	Position Point `json:"position"`
}

func NewFood(width, height int) *Food {
	return &Food{
		Position: Point{
			X: rand.Intn(width),
			Y: rand.Intn(height),
		},
	}
}
