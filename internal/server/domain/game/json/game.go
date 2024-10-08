package json

import (
	"github.com/google/uuid"
)

type SnakesJson struct {
	Data map[uuid.UUID]*SnakeJson `json:"data"`
}

type GameJson struct {
	Id        string            `json:"id"`
	Width     int               `json:"width"`
	Height    int               `json:"height"`
	Obstacles [][2]int32        `json:"obstacles"`
	Snakes    SnakesJson        `json:"snakes"`
	Scores    map[uuid.UUID]int `json:"scores"`
	Food      FoodJson          `json:"food"`
}
