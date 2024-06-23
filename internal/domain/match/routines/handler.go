package match_routines

import (
	"time"

	gamedata "snake_ai/internal/domain/game/data"
	gameroutines "snake_ai/internal/domain/game/routines"
	matchdata "snake_ai/internal/domain/match/data"
	"snake_ai/internal/domain/ws"
	"snake_ai/pkg/logger"
)

var PartiesChannel = make(chan *matchdata.Party, 100)

func HandlePartyMessages() {
	for {
		pa := <-PartiesChannel

		// TODO change default interval
		g := gamedata.NewGame(pa.Width, pa.Height, pa)
		gamedata.CurrentGames.AddGame(g)
		go gameroutines.HandleGames(g, *time.NewTicker(time.Duration(1) * time.Second))

		for _, p := range pa.Players {
			conn, exists := ws.Connections.Get(p.Id)
			if exists {
				err := conn.WriteJSON(pa)
				if err != nil {
					conn.Close()
					ws.Connections.Remove(p.Id)
					logger.Log.Errorf("error writing to websocket: %s", err.Error())
				}
				logger.Log.Infof("player with ID %s found party with ID %s", p.Id, pa.Id)
			}
		}
	}
}
