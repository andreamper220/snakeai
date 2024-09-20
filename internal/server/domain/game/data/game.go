package data

import (
	"context"
	"github.com/andreamper220/snakeai/internal/server/domain"
	matchdata "github.com/andreamper220/snakeai/internal/server/domain/match/data"
	grpcclients "github.com/andreamper220/snakeai/internal/server/infrastructure/grpc"
	"github.com/andreamper220/snakeai/internal/server/infrastructure/storages"
	pb "github.com/andreamper220/snakeai/proto"
	"github.com/google/uuid"
	"sync"

	"github.com/andreamper220/snakeai/pkg/logger"
)

var CurrentGames Games

type Games struct {
	mux   sync.RWMutex
	Games []*Game
}

func (games *Games) GetGames() []*Game {
	games.mux.RLock()
	defer games.mux.RUnlock()
	return games.Games
}
func (games *Games) AddGame(game *Game) {
	gg := games.GetGames()
	for _, g := range gg {
		if g == game {
			games.RemoveGame(game)
			break
		}
	}
	games.Games = append(games.Games, game)
}
func (games *Games) RemoveGame(game *Game) {
	result := make([]*Game, 0)
	gg := games.GetGames()
	for _, g := range gg {
		if g != game {
			result = append(result, g)
		}
	}
	games.Games = result
}

type Snakes struct {
	sync.RWMutex
	Data map[uuid.UUID]*Snake `json:"data"`
}

type Game struct {
	sync.RWMutex
	Id     string
	Width  int
	Height int
	Party  *matchdata.Party
	Snakes Snakes
	Scores map[uuid.UUID]int
	Food   *Food
	Done   chan bool
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
func (g *Game) AddSnake(snake *Snake, userId uuid.UUID) {
	g.RLock()
	snake.Game = g
	g.RUnlock()
	g.Snakes.Lock()
	g.Snakes.Data[userId] = snake
	g.Snakes.Unlock()
}
func (g *Game) RemoveSnake(userId uuid.UUID) {
	g.Snakes.Lock()
	defer g.Snakes.Unlock()
	_, exists := g.Snakes.Data[userId]
	if exists {
		delete(g.Snakes.Data, userId)
	}
}
func (g *Game) GetSnakes() map[uuid.UUID]*Snake {
	g.Snakes.RLock()
	defer g.Snakes.RUnlock()
	return g.Snakes.Data
}
func (g *Game) GetUserIdBySnake(snake *Snake) uuid.UUID {
	g.Snakes.RLock()
	defer g.Snakes.RUnlock()
	for userId, s := range g.Snakes.Data {
		if s == snake {
			return userId
		}
	}
	return uuid.Nil
}
func (g *Game) Update() {
	snakeJobsChannel := make(chan *Snake, 10)
	numSnakeWorkers := 8
	for w := 0; w < numSnakeWorkers; w++ {
		go func() {
			for s := range snakeJobsChannel {
				if len(s.AiFunc) > 0 {
					//s.Lock()
					s.AiFunc[s.AiFuncNum](s)
					s.AiFuncNum++
					if len(s.AiFunc) == s.AiFuncNum {
						s.AiFuncNum = 0
					}
					//s.Unlock()
					g.handleCollisions(s)
				}
			}
		}()
	}

	snakes := g.GetSnakes()
	for _, snake := range snakes {
		snake.RLock()
		snakeJobsChannel <- snake
		snake.RUnlock()
	}
	close(snakeJobsChannel)
}
func (g *Game) handleCollisions(snake *Snake) {
	snake.RLock()
	defer snake.RUnlock()
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
	// edge custom walls collision
	gameMap, err := grpcclients.EditorClient.GetMap(context.Background(), &pb.GetMapRequest{
		Id: g.Party.MapId,
	})
	if err == nil {
		for _, obstacle := range gameMap.Map.Struct.Obstacles {
			if head.X-1 == int(obstacle.Cx) && head.Y-1 == int(obstacle.Cy) {
				g.RemoveSnake(userId)
				return
			}
		}
	}
	// self-collision
	for _, part := range snake.Body[1:] {
		if head.X == part.X && head.Y == part.Y {
			g.RemoveSnake(userId)
			return
		}
	}
	// another snake collision
	for targetUserId, targetSnake := range g.GetSnakes() {
		if snake == targetSnake {
			continue
		}
		targetSnake.RLock()
		for _, part := range targetSnake.Body {
			if head.X == part.X && head.Y == part.Y {
				g.RemoveSnake(userId)
				g.RemoveSnake(targetUserId)
				targetSnake.RUnlock()
				return
			}
		}
		targetSnake.RUnlock()
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

func GetGameByPlayer(userId uuid.UUID) *Game {
	for _, g := range CurrentGames.Games {
		for _, p := range g.Party.Players {
			if p.Id == userId {
				return g
			}
		}
	}
	return nil
}
