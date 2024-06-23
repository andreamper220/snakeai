package delete_handlers

import (
	"github.com/google/uuid"
	"net/http"

	gamedata "snake_ai/internal/domain/game/data"
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
