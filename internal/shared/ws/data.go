package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var Connections ConnectionsMap

type ConnectionsMap map[uuid.UUID]*websocket.Conn
