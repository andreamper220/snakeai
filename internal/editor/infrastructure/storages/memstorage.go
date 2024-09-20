package storages

import (
	"sync"

	"github.com/andreamper220/snakeai/internal/editor/domain"
	"github.com/google/uuid"
)

type MemStorage struct {
	mu   sync.RWMutex
	maps []*domain.Map
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		maps: make([]*domain.Map, 0),
	}
}
func (ms *MemStorage) AddMap(gameMap *domain.Map) (uuid.UUID, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	gameMap.Id = uuid.New()
	ms.maps = append(ms.maps, gameMap)

	return gameMap.Id, nil
}
func (ms *MemStorage) GetMap(id uuid.UUID) (*domain.Map, error) {
	for _, m := range ms.maps {
		if m.Id == id {
			return m, nil
		}
	}
	return nil, ErrMapNotFound
}
