package delete_handlers

import (
	gamedata "github.com/andreamper220/snakeai/internal/server/domain/game/data"
	"github.com/google/uuid"
	"net/http"
)

func PlayerRemoveAi(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
out:
	for _, g := range gamedata.CurrentGames.Games {
		for _, p := range g.Party.Players {
			if p.Id == userId {
				g.RemoveSnake(userId)
				break out
			}
		}
	}
}
