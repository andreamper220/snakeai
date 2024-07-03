package ws

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"sync"

	"snakeai/pkg/logger"
)

var ErrConnectionNotFound = errors.New("connection not found")

var Connections connections

type connection struct {
	mu              sync.Mutex
	messagesChannel chan []byte
	closeChannel    chan bool
}
type connections struct {
	mu    sync.RWMutex
	conns map[uuid.UUID]connection
}

func (c *connections) Add(userId uuid.UUID, messagesChannel chan []byte, closeChannel chan bool) {
	if c.conns == nil {
		c.conns = make(map[uuid.UUID]connection)
	} else {
		c.mu.RLock()
		_, exists := c.conns[userId]
		c.mu.RUnlock()
		if exists {
			return
		}
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
	conn, exists := c.conns[userId]
	if exists {
		conn.closeChannel <- true
		delete(c.conns, userId)
	}
}
func (c *connections) WriteJSON(userId uuid.UUID, data interface{}) error {
	c.mu.RLock()
	conn, exists := c.conns[userId]
	c.mu.RUnlock()

	if exists {
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
func (c *connections) Write(userId uuid.UUID, data []byte) error {
	c.mu.RLock()
	conn, exists := c.conns[userId]
	c.mu.RUnlock()

	if exists {
		conn.messagesChannel <- data
		return nil
	}
	return ErrConnectionNotFound
}
