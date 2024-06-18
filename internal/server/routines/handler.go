package routines

import (
	"github.com/gorilla/websocket"

	"snake_ai/internal/logger"
	"snake_ai/internal/shared/match/data"
	"snake_ai/internal/shared/ws"
)

var PartiesChannel = make(chan *data.Party, 100)

func HandlePartyMessages() {
	for {
		pa := <-PartiesChannel
		for _, p := range pa.Players {
			conn := ws.Connections[p.Id]
			if conn != nil {
				logger.Log.Infof("found party: %v", pa)
				err := conn.WriteMessage(websocket.TextMessage, []byte("found party"))
				if err != nil {
					conn.Close()
					delete(ws.Connections, p.Id)
					logger.Log.Errorf("error writing to websocket: %s", err.Error())
				}
			}
		}
	}
}
