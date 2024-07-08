package caches

import (
	"errors"
	"time"
)

var ErrNoSession = errors.New("no session found")

var Cache CacheInterface

type CacheInterface interface {
	AddSession(sessionId, userId string, expired time.Duration) error
	DelSession(sessionId string) error
}
