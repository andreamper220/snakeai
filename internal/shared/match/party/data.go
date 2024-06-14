package party

import (
	"sync"

	"snake_ai/internal/shared/match/player"
)

type Party struct {
	mux       sync.Mutex
	id        string
	players   []*player.Player
	avgSkill  int
	createdAt int64
}
