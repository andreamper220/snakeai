package data

import (
	"github.com/google/uuid"
	"sync"

	"snake_ai/internal/shared"
	"snake_ai/internal/shared/match/data"
)

var Games []*Game

type Game struct {
	mux    sync.Mutex
	Id     string               `json:"-"`
	Width  int                  `json:"width"`
	Height int                  `json:"height"`
	Party  *data.Party          `json:"-"`
	Snakes map[uuid.UUID]*Snake `json:"snakes"`
	Food   *Food                `json:"food"`
	Done   chan bool            `json:"-"`
}

func NewGame(width, height int, party *data.Party) *Game {
	game := &Game{
		Id:     shared.RandSeq(10),
		Width:  width,
		Height: height,
		Party:  party,
		Snakes: make(map[uuid.UUID]*Snake),
		Food:   NewFood(width, height),
		Done:   make(chan bool),
	}

	return game
}
func (g *Game) lock() {
	g.mux.Lock()
}
func (g *Game) unlock() {
	g.mux.Unlock()
}
func (g *Game) Update() {
	g.lock()
	defer g.unlock()

	snakeJobsChannel := make(chan *Snake, 10)
	numSnakeWorkers := 8
	for w := 0; w < numSnakeWorkers; w++ {
		go func() {
			for s := range snakeJobsChannel {
				s.Lock()
				s.AiFunc[s.AIFuncNum](s)
				// TODO handle collisions + food eating
				s.AIFuncNum++
				if len(s.AiFunc) == s.AIFuncNum {
					s.AIFuncNum = 0
				}
				s.Unlock()
			}
		}()
	}

	for _, snake := range g.Snakes {
		snakeJobsChannel <- snake
	}
	close(snakeJobsChannel)
}
