package data

import (
	"github.com/google/uuid"
	"sync"

	"snake_ai/internal/shared"
	"snake_ai/internal/shared/match/data"
)

var CurrentGames Games

type Games struct {
	mux   sync.Mutex
	Games []*Game
}

func (games *Games) AddGame(game *Game) {
	games.mux.Lock()
	for _, g := range games.Games {
		if g == game {
			games.RemoveGame(game)
			break
		}
	}
	games.Games = append(games.Games, game)
	games.mux.Unlock()
}
func (games *Games) RemoveGame(game *Game) {
	games.mux.Lock()
	result := make([]*Game, len(games.Games))
	for _, g := range games.Games {
		if g != game {
			result = append(result, g)
		}
	}
	games.Games = result
	games.mux.Unlock()
}

type Game struct {
	mux    sync.Mutex
	Id     string            `json:"-"`
	Width  int               `json:"width"`
	Height int               `json:"height"`
	Party  *data.Party       `json:"-"`
	Snakes []*Snake          `json:"snakes"`
	Scores map[uuid.UUID]int `json:"scores"`
	Food   *Food             `json:"food"`
	Done   chan bool         `json:"-"`
}

func NewGame(width, height int, party *data.Party) *Game {
	game := &Game{
		Id:     shared.RandSeq(10),
		Width:  width,
		Height: height,
		Party:  party,
		Snakes: make([]*Snake, 0),
		Scores: make(map[uuid.UUID]int),
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
func (g *Game) AddSnake(snake *Snake, userId uuid.UUID) {
	g.lock()
	defer g.unlock()
	for _, s := range g.Snakes {
		if s.UserId == userId {
			g.RemoveSnake(s)
			break
		}
	}
	g.Snakes = append(g.Snakes, snake)
}
func (g *Game) RemoveSnake(snake *Snake) {
	result := make([]*Snake, 0)
	g.lock()
	defer g.unlock()
	for _, s := range g.Snakes {
		if s != snake {
			result = append(result, s)
		}
	}
	g.Snakes = result
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
				g.handleCollisions(s)
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
func (g *Game) handleCollisions(snake *Snake) {
	head := snake.Body[0]
	// edge walls collision
	if head.X == 0 || head.Y == 0 || head.X > g.Width || head.Y > g.Height {
		g.RemoveSnake(snake)
		return
	}
	// self-collision
	for _, part := range snake.Body[1:] {
		if head.X == part.X && head.Y == part.Y {
			g.RemoveSnake(snake)
			return
		}
	}
	// another snake collision
	for _, targetSnake := range g.Snakes {
		for _, part := range targetSnake.Body[1:] {
			if head.X == part.X && head.Y == part.Y {
				g.RemoveSnake(snake)
				g.RemoveSnake(targetSnake)
				return
			}
		}
	}
	// food eating
	if head.X == g.Food.Position.X && head.Y == g.Food.Position.Y {
		g.Scores[snake.UserId]++
		snake.GrowCounter += 1
		g.Food = NewFood(g.Width, g.Height)
	}
}
