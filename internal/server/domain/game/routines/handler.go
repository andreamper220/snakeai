package routines

import (
	gamedata "github.com/andreamper220/snakeai/internal/server/domain/game/data"
	"github.com/andreamper220/snakeai/internal/server/domain/game/json"
	"github.com/andreamper220/snakeai/internal/server/domain/ws"
	"github.com/google/uuid"
	"time"

	"github.com/andreamper220/snakeai/pkg/logger"
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
			game.RUnlock()

			for _, p := range game.Party.Players {
				err := ws.Connections.WriteJSON(p.Id, gameJson)
				if err != nil {
					ws.Connections.Remove(p.Id)
					logger.Log.Errorf("error writing to websocket: %s", err.Error())
				}
			}
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
