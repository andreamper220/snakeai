package delete_handlers

import (
	"github.com/google/uuid"
	"net/http"
	game "snake_ai/internal/server/ai/data"
)

func PlayerRemoveAi(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
out:
	for _, g := range game.CurrentGames.Games {
		for _, p := range g.Party.Players {
			if p.Id == userId {
				for _, s := range g.Snakes {
					if s.UserId == userId {
						g.RemoveSnake(s)
						break out
					}
				}
			}
		}
	}
}
