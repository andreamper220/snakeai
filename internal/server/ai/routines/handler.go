package routines

import (
	"snake_ai/internal/logger"
	"snake_ai/internal/server/ai/data"
	"snake_ai/internal/shared/ws"
	"time"
)

func HandleGames(game *data.Game, gameTicker time.Ticker) {
	for {
		select {
		case <-gameTicker.C: // update game field via ws
			for _, p := range game.Party.Players {
				conn, exists := ws.Connections.Get(p.Id)
				if exists {
					err := conn.WriteJSON(game)
					if err != nil {
						conn.Close()
						ws.Connections.Remove(p.Id)
						logger.Log.Errorf("error writing to websocket: %s", err.Error())
					}
				}
			}
			GameJobsChannel <- game
		case <-game.Done:
			gameTicker.Stop()
			logger.Log.Infof("party with ID %s disbanded", game.Party.Id)
			return
		}
	}
}
