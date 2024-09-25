package routines

import (
	"context"
	"github.com/google/uuid"
	"time"

	gamedata "github.com/andreamper220/snakeai/internal/server/domain/game/data"
	"github.com/andreamper220/snakeai/internal/server/domain/game/json"
	"github.com/andreamper220/snakeai/internal/server/domain/ws"
	grpcclients "github.com/andreamper220/snakeai/internal/server/infrastructure/grpc"
	"github.com/andreamper220/snakeai/pkg/logger"
	pb "github.com/andreamper220/snakeai/proto"
)

func HandleGames(game *gamedata.Game, gameTicker time.Ticker) {
	for {
		select {
		case <-gameTicker.C: // update game field via ws
			game.RLock()
			snakes := json.SnakesJson{
				Data: make(map[uuid.UUID]*json.SnakeJson),
			}
			for userId, snake := range game.GetSnakes() {
				game.Snakes.RLock()
				body := make([]json.PointJson, len(snake.Body))
				for i, point := range snake.Body {
					body[i] = json.PointJson{X: point.X, Y: point.Y}
				}
				snakes.Data[userId] = &json.SnakeJson{
					Color: snake.Color,
					Body:  body,
				}
				game.Snakes.RUnlock()
			}
			scores := game.Scores

			gameJson := json.GameJson{
				Id:     game.Id,
				Width:  game.Width,
				Height: game.Height,
				Snakes: snakes,
				Scores: scores,
				Food: json.FoodJson{
					Position: json.PointJson{
						X: game.Food.Position.X,
						Y: game.Food.Position.Y,
					},
				},
			}

			if game.Party.MapId != "" {
				gameMap, err := grpcclients.EditorClient.GetMap(context.Background(), &pb.GetMapRequest{
					Id: game.Party.MapId,
				})

				if err == nil {
					mapObstacles := gameMap.GetMap().GetStruct().GetObstacles()
					obstacles := make([][2]int32, len(mapObstacles))
					for i := 0; i < len(mapObstacles); i++ {
						obstacles[i][0] = mapObstacles[i].GetCx()
						obstacles[i][1] = mapObstacles[i].GetCy()
					}
					gameJson.Obstacles = obstacles
				}
			}
			game.RUnlock()

			broadcastGameState(game, gameJson)
			game.Lock()
			GameJobsChannel <- game
			game.Unlock()
		case <-game.Done:
			gameTicker.Stop()
			logger.Log.Infof("party with ID %s disbanded", game.Party.Id)
			return
		}
	}
}

func broadcastGameState(game *gamedata.Game, gameJson json.GameJson) {
	for _, p := range game.Party.Players {
		err := ws.Connections.WriteJSON(p.Id, gameJson)
		if err != nil {
			ws.Connections.Remove(p.Id)
			logger.Log.Errorf("error writing to websocket: %s", err.Error())
		}
	}
}
