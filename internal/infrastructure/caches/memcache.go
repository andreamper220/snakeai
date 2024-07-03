package caches

import (
	"sync"
	"time"
)

type MemCache struct {
	mu       sync.RWMutex
	sessions map[string]string
}

func NewMemCache() *MemCache {
	return &MemCache{
		sessions: make(map[string]string),
	}
}
func (mc *MemCache) AddSession(sessionId, userId string, expired time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.sessions["sessionId_"+sessionId] = userId
	return nil
}
func (mc *MemCache) DelSession(sessionId string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	_, ok := mc.sessions["sessionId_"+sessionId]
	if ok {
		delete(mc.sessions, "sessionId_"+sessionId)
	}
	return nil
}
