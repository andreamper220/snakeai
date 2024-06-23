package caches

import (
	"time"
)

var Cache CacheInterface

type CacheInterface interface {
	AddSession(sessionId, userId string, expired time.Duration) error
	DelSession(sessionId string) error
}
