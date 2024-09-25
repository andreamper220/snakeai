package ws_handlers

import (
	"github.com/andreamper220/snakeai/internal/server/domain/ws"
	"github.com/andreamper220/snakeai/pkg/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
)

func PlayerConnection(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	if ws.Connections.Exists(userId) {
		ws.Connections.Remove(userId)
	}

	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "something happened establishing your websocket connection", http.StatusInternalServerError)
		return
	}

	// TODO unbuffer ?
	closeChannel := make(chan bool, 1)
	messagesChannel := make(chan []byte, 10)
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
