package domain

import "github.com/google/uuid"

type Map struct {
	Id        uuid.UUID
	Width     int32
	Height    int32
	Obstacles [][2]int32
}
