package routines

import (
	"snake_ai/internal/logger"
	"snake_ai/internal/shared/match/data"
)

var PartiesChannel = make(chan *data.Party, 100)

func HandlePartyMessages() {
	for {
		pa := <-PartiesChannel
		logger.Log.Infof("found party: %v", pa)
		//for _, player := range pa.Players {
		//	playerName := player.name
		//	conn := server.connectionsMap[playerName]
		//	if conn != nil {
		//		msg := fmt.Sprintf("found party : %v", pa)
		//		err := conn.WriteJSON(msg)
		//		if err != nil {
		//			fmt.Println(err)
		//			conn.Close()
		//			delete(server.connectionsMap, playerName)
		//		}
		//	}
		//}
	}
}
