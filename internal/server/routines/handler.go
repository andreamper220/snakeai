package routines

import (
	"snake_ai/internal/logger"
	"snake_ai/internal/shared/match/data"
	"snake_ai/internal/shared/ws"
)

var PartiesChannel = make(chan *data.Party, 100)

func HandlePartyMessages() {
	for {
		pa := <-PartiesChannel
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
