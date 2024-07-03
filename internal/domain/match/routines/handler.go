package match_routines

import (
	"time"

	gamedata "snakeai/internal/domain/game/data"
	gameroutines "snakeai/internal/domain/game/routines"
	matchdata "snakeai/internal/domain/match/data"
	"snakeai/internal/domain/ws"
	"snakeai/pkg/logger"
)

var PartiesChannel = make(chan *matchdata.Party, 100)

func HandlePartyMessages() {
	for {
		pa := <-PartiesChannel

		// TODO change default interval ?
		g := gamedata.NewGame(pa.Width, pa.Height, pa)
		gamedata.CurrentGames.AddGame(g)
		go gameroutines.HandleGames(g, *time.NewTicker(time.Duration(1) * time.Second))

		for _, p := range pa.Players {
			err := ws.Connections.WriteJSON(p.Id, pa)
			if err != nil {
				ws.Connections.Remove(p.Id)
				logger.Log.Errorf("error writing to websocket: %s", err.Error())
			} else {
				logger.Log.Infof("player with ID %s found party with ID %s", p.Id, pa.Id)
			}
		}
	}
}
