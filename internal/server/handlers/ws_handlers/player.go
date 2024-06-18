package ws_handlers

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"

	"snake_ai/internal/logger"
	"snake_ai/internal/shared/ws"
)

func PlayerConnection(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "something happened establishing your websocket connection", http.StatusInternalServerError)
		return
	}

	defer c.Close()
	for {
		conn, exists := ws.Connections[userId]
		if !exists {
			ws.Connections[userId] = c
			conn = c
			logger.Log.Infof("ws connection added to player with ID = %s", userId.String())
		}
		//w, err := c.conn.NextWriter(websocket.TextMessage)
		//if err != nil {
		//	return
		//}
		//w.Write(message)
		conn.WriteMessage(websocket.TextMessage, []byte("connection established"))
	}
}
