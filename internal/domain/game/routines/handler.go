package routines

import (
	"github.com/google/uuid"
	"time"

	gamedata "github.com/andreamper220/snakeai.git/internal/domain/game/data"
	gamejson "github.com/andreamper220/snakeai.git/internal/domain/game/json"
	"github.com/andreamper220/snakeai.git/internal/domain/ws"
	"github.com/andreamper220/snakeai.git/pkg/logger"
)

func HandleGames(game *gamedata.Game, gameTicker time.Ticker) {
	for {
		select {
		case <-gameTicker.C: // update game field via ws
			game.RLock()
			snakes := gamejson.SnakesJson{
				Data: make(map[uuid.UUID]*gamejson.SnakeJson),
			}
			for userId, snake := range game.GetSnakes() {
				game.Snakes.RLock()
				body := make([]gamejson.PointJson, len(snake.Body))
				for i, point := range snake.Body {
					body[i] = gamejson.PointJson{X: point.X, Y: point.Y}
				}
				snakes.Data[userId] = &gamejson.SnakeJson{
					Color: snake.Color,
					Body:  body,
				}
				game.Snakes.RUnlock()
			}
			scores := game.Scores
			gameJson := gamejson.GameJson{
				Id:     game.Id,
				Width:  game.Width,
				Height: game.Height,
				Snakes: snakes,
				Scores: scores,
				Food: gamejson.FoodJson{
					Position: gamejson.PointJson{
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
