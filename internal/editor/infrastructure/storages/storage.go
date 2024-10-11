package storages

import (
	"errors"
	"github.com/google/uuid"

	"github.com/andreamper220/snakeai/internal/editor/domain"
)

var ErrMapNotFound = errors.New("map not found")

var EditorStorage EditorStorageInterface

type EditorStorageInterface interface {
	AddMap(gameMap *domain.Map) (uuid.UUID, error)
	GetMap(id uuid.UUID) (*domain.Map, error)
}
