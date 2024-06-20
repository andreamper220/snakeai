package data

import "sync"

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
	mux         sync.Mutex
	Body        []Point              `json:"body"`
	Direction   Point                `json:"-"`
	GrowCounter int                  `json:"-"`
	AiFunc      []func(snake *Snake) `json:"-"`
	AIFuncNum   int                  `json:"-"`
}

func NewSnake(x, y, xTo, yTo int, aiFunc []func(snake *Snake)) *Snake {
	return &Snake{
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
func (s *Snake) Down() {
	if s.Direction.Y == 0 {
		s.Direction = Point{X: 0, Y: 1}
	}
}
func (s *Snake) Up() {
	if s.Direction.Y == 0 {
		s.Direction = Point{X: 0, Y: -1}
	}
}
func (s *Snake) Left() {
	if s.Direction.X == 0 {
		s.Direction = Point{X: -1, Y: 0}
	}
}
func (s *Snake) Right() {
	if s.Direction.X == 0 {
		s.Direction = Point{X: 1, Y: 0}
	}
}
