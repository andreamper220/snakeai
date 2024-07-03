package ws_handlers

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"snake_ai/pkg/logger"

	"snake_ai/internal/domain/ws"
)

func PlayerConnection(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "something happened establishing your websocket connection", http.StatusInternalServerError)
		return
	}

	closeChannel := make(chan bool, 1)
	messagesChannel := make(chan []byte, 100)
	go func(conn *websocket.Conn) {
		defer close(closeChannel)
		defer close(messagesChannel)
		for {
			select {
			case message := <-messagesChannel:
				err = conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					logger.Log.Errorf("error writing to websocket: %s", err.Error())
				}
			case <-closeChannel:
				return
			}
		}
	}(c)
	ws.Connections.Add(userId, messagesChannel, closeChannel)
}
