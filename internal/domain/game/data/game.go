package data

import (
	"github.com/google/uuid"
	"sync"

	"snake_ai/internal/domain"
	matchdata "snake_ai/internal/domain/match/data"
	"snake_ai/internal/infrastructure/storages"
	"snake_ai/pkg/logger"
)

var CurrentGames Games

type Games struct {
	mux   sync.Mutex
	Games []*Game
}

func (games *Games) AddGame(game *Game) {
	games.mux.Lock()
	defer games.mux.Unlock()
	for _, g := range games.Games {
		if g == game {
			games.RemoveGame(game)
			break
		}
	}
	games.Games = append(games.Games, game)
}
func (games *Games) RemoveGame(game *Game) {
	games.mux.Lock()
	defer games.mux.Unlock()
	result := make([]*Game, 0)
	for _, g := range games.Games {
		if g != game {
			result = append(result, g)
		}
	}
	games.Games = result
}

type Snakes struct {
	mu   sync.RWMutex
	Data map[uuid.UUID]*Snake `json:"data"`
}

type Game struct {
	mux    sync.Mutex
	Id     string            `json:"id"`
	Width  int               `json:"width"`
	Height int               `json:"height"`
	Party  *matchdata.Party  `json:"-"`
	Snakes Snakes            `json:"snakes"`
	Scores map[uuid.UUID]int `json:"scores"`
	Food   *Food             `json:"food"`
	Done   chan bool         `json:"-"`
}

func NewGame(width, height int, party *matchdata.Party) *Game {
	game := &Game{
		Id:     domain.RandSeq(10),
		Width:  width,
		Height: height,
		Party:  party,
		Snakes: Snakes{
			Data: make(map[uuid.UUID]*Snake),
		},
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
	g.Snakes.mu.Lock()
	defer g.Snakes.mu.Unlock()
	g.Snakes.Data[userId] = snake
}
func (g *Game) RemoveSnake(userId uuid.UUID) {
	g.Snakes.mu.Lock()
	defer g.Snakes.mu.Unlock()
	_, exists := g.Snakes.Data[userId]
	if exists {
		delete(g.Snakes.Data, userId)
	}
}
func (g *Game) GetUserIdBySnake(snake *Snake) uuid.UUID {
	g.Snakes.mu.RLock()
	defer g.Snakes.mu.RUnlock()
	for userId, s := range g.Snakes.Data {
		if s == snake {
			return userId
		}
	}
	return uuid.Nil
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

	for _, snake := range g.Snakes.Data {
		snakeJobsChannel <- snake
	}
	close(snakeJobsChannel)
}
func (g *Game) handleCollisions(snake *Snake) {
	userId := g.GetUserIdBySnake(snake)
	if userId == uuid.Nil {
		return
	}

	head := snake.Body[0]
	// edge walls collision
	if head.X == 0 || head.Y == 0 || head.X > g.Width || head.Y > g.Height {
		g.RemoveSnake(userId)
		return
	}
	// self-collision
	for _, part := range snake.Body[1:] {
		if head.X == part.X && head.Y == part.Y {
			g.RemoveSnake(userId)
			return
		}
	}
	// another snake collision
	for targetUserId, targetSnake := range g.Snakes.Data {
		if snake == targetSnake {
			continue
		}
		for _, part := range targetSnake.Body {
			if head.X == part.X && head.Y == part.Y {
				g.RemoveSnake(userId)
				g.RemoveSnake(targetUserId)
				return
			}
		}
	}
	// food eating
	if head.X == g.Food.Position.X && head.Y == g.Food.Position.Y {
		g.Food = NewFood(g.Width, g.Height)
		snake.GrowCounter += 1
		g.Scores[userId]++
		if err := storages.Storage.IncreasePlayerScore(userId); err != nil {
			logger.Log.Error(err.Error())
		}
	}
}

func RemovePlayer(userId uuid.UUID) {
out:
	for _, g := range CurrentGames.Games {
		for _, p := range g.Party.Players {
			if p.Id == userId {
				g.RemoveSnake(userId)
				g.Party.RemovePlayer(p)
				logger.Log.Infof("user with ID %s exited from party with ID %s", userId.String(), g.Party.Id)
				if len(g.Party.Players) == 0 {
					g.Done <- true
					CurrentGames.RemoveGame(g)
				}
				break out
			}
		}
	}
}
