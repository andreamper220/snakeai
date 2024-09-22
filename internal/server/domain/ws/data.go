package ws

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"sync"

	"github.com/andreamper220/snakeai/pkg/logger"
)

var ErrConnectionNotFound = errors.New("connection not found")

var Connections connections

type connection struct {
	messagesChannel chan []byte
	closeChannel    chan bool
}
type connections struct {
	mu    sync.RWMutex
	conns map[uuid.UUID]connection
}

func (c *connections) Exists(userId uuid.UUID) bool {
	if c.conns == nil {
		c.conns = make(map[uuid.UUID]connection)
	} else {
		c.mu.RLock()
		if _, exists := c.conns[userId]; exists {
			c.mu.RUnlock()
			return true
		}
		c.mu.RUnlock()
	}
	return false
}
func (c *connections) Add(userId uuid.UUID, messagesChannel chan []byte, closeChannel chan bool) {
	if c.Exists(userId) {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.conns[userId] = connection{
		messagesChannel: messagesChannel,
		closeChannel:    closeChannel,
	}
	logger.Log.Infof("ws connection added to player with ID = %s", userId.String())
}
func (c *connections) Remove(userId uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if conn, exists := c.conns[userId]; exists {
		conn.closeChannel <- true
		delete(c.conns, userId)
	}
}
func (c *connections) WriteJSON(userId uuid.UUID, data interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if conn, exists := c.conns[userId]; exists {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(data); err != nil {
			return err
		}
		mess := buf.Bytes()
		conn.messagesChannel <- mess
		return nil
	}
	return ErrConnectionNotFound
}
