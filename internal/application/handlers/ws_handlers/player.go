package ws_handlers

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"

	"snake_ai/internal/domain/ws"
	"snake_ai/pkg/logger"
)

func PlayerConnection(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "something happened establishing your websocket connection", http.StatusInternalServerError)
		return
	}

	_, exists := ws.Connections.Get(userId)
	if !exists {
		ws.Connections.Add(userId, c)
		logger.Log.Infof("ws connection added to player with ID = %s", userId.String())
	}
}
