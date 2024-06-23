package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

var Connections connections

type connections struct {
	mu    sync.RWMutex
	conns map[uuid.UUID]*websocket.Conn
}

func (c *connections) Add(userId uuid.UUID, conn *websocket.Conn) {
	if c.conns == nil {
		c.conns = make(map[uuid.UUID]*websocket.Conn)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.conns[userId] = conn
}
func (c *connections) Remove(userId uuid.UUID) {
	c.mu.Lock()
	conn, exists := c.conns[userId]
	if exists {
		conn.Close()
		delete(c.conns, userId)
	}
	c.mu.Unlock()
}
func (c *connections) Get(userId uuid.UUID) (*websocket.Conn, bool) {
	c.mu.RLock()
	conn, exists := c.conns[userId]
	c.mu.RUnlock()
	return conn, exists
}
