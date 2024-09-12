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
	sync.RWMutex
	Color         string
	Body          []Point
	Direction     Point
	GrowCounter   int
	AiFunc        []func(snake *Snake)
	AiFuncNum     int
	AiFuncNumPrev int
	Game          *Game
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
func (s *Snake) GetGame() *Game {
	s.Game.RLock()
	defer s.Game.RUnlock()
	return s.Game
}
func (s *Snake) Move() {
	s.AiFuncNumPrev = s.AiFuncNum
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
	s.AiFuncNumPrev = s.AiFuncNum
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
	s.AiFuncNumPrev = s.AiFuncNum
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
func (s *Snake) DoIf(condition AiCondition, actionsLength int) {
	// if condition true - to run next action (by default)
	// else - rewind to end of condition actions
	s.AiFuncNumPrev = s.AiFuncNum
	if !condition.Check(s, s.GetGame()) {
		s.AiFuncNum += actionsLength
	}
}
func (s *Snake) DoElseIf(condition AiCondition, actionsLength int) {
	// if only from previous if / elseif
	if s.AiFuncNum-s.AiFuncNumPrev <= 1 {
		s.AiFuncNumPrev += actionsLength + 1
		s.AiFuncNum += actionsLength
	} else {
		s.DoIf(condition, actionsLength)
	}
}
func (s *Snake) DoElse(actionsLength int) {
	// if only from previous if / elseif
	// else - rewind to end of condition actions
	if s.AiFuncNum-s.AiFuncNumPrev <= 1 {
		s.AiFuncNumPrev = s.AiFuncNum
		s.AiFuncNum += actionsLength
	}
}
