package data

import (
	"golang.org/x/exp/rand"
	"sync"
	"time"
)

var colors = [10]string{"yellow", "blue", "red", "purple", "pink", "orange", "black", "brown", "cyan", "gray"}

type Point struct {
	X int
	Y int
}

type Snake struct {
	mux         sync.Mutex
	Color       string
	Body        []Point
	Direction   Point
	GrowCounter int
	AiFunc      []func(snake *Snake)
	AIFuncNum   int
}

func NewSnake(x, y, xTo, yTo int, aiFunc []func(snake *Snake)) *Snake {
	rand.Seed(uint64(time.Now().Unix()))
	return &Snake{
		Color: colors[rand.Intn(len(colors))],
		Body: []Point{
			{X: x, Y: y},
		},
		Direction: Point{X: xTo, Y: yTo},
		AiFunc:    aiFunc,
	}
}
func (s *Snake) Lock() {
	s.mux.Lock()
}
func (s *Snake) Unlock() {
	s.mux.Unlock()
}
func (s *Snake) Move() {
	newHead := Point{
		X: s.Body[0].X + s.Direction.X,
		Y: s.Body[0].Y + s.Direction.Y,
	}
	s.Body = append([]Point{newHead}, s.Body...)

	if s.GrowCounter > 0 {
		s.GrowCounter--
	} else {
		s.Body = s.Body[:len(s.Body)-1]
	}
}
func (s *Snake) Left() {
	if s.Direction.X == 0 {
		if s.Direction.Y == 1 {
			s.Direction = Point{X: 1, Y: 0}
		} else {
			s.Direction = Point{X: -1, Y: 0}
		}
	} else if s.Direction.Y == 0 {
		if s.Direction.X == 1 {
			s.Direction = Point{X: 0, Y: -1}
		} else {
			s.Direction = Point{X: 0, Y: 1}
		}
	}
}
func (s *Snake) Right() {
	if s.Direction.X == 0 {
		if s.Direction.Y == 1 {
			s.Direction = Point{X: -1, Y: 0}
		} else {
			s.Direction = Point{X: 1, Y: 0}
		}
	} else if s.Direction.Y == 0 {
		if s.Direction.X == 1 {
			s.Direction = Point{X: 0, Y: 1}
		} else {
			s.Direction = Point{X: 0, Y: -1}
		}
	}
}
