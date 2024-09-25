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
	conns sync.Map
}

func (c *connections) Exists(userId uuid.UUID) bool {
	_, exists := c.conns.Load(userId)
	return exists
}
func (c *connections) Add(userId uuid.UUID, messagesChannel chan []byte, closeChannel chan bool) {
	c.conns.LoadOrStore(userId, connection{
		messagesChannel: messagesChannel,
		closeChannel:    closeChannel,
	})
	logger.Log.Infof("ws connection added to player with ID = %s", userId.String())
}
func (c *connections) Remove(userId uuid.UUID) {
	value, existed := c.conns.LoadAndDelete(userId)
	if existed {
		conn := value.(connection)
		conn.closeChannel <- true
	}
}
func (c *connections) WriteJSON(userId uuid.UUID, data interface{}) error {
	if value, exists := c.conns.Load(userId); exists {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(data); err != nil {
			return err
		}
		mess := buf.Bytes()
		conn := value.(connection)
		conn.messagesChannel <- mess
		return nil
	}
	return ErrConnectionNotFound
}
