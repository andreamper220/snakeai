package post_handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	game "snake_ai/internal/server/ai/data"

	"snake_ai/internal/server/routines"
	"snake_ai/internal/server/storages"
	"snake_ai/internal/shared/match/data"
	js "snake_ai/internal/shared/match/json"
)

func PlayerPartyEnqueue(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var partyJson js.PartyJson
	if err := json.NewDecoder(r.Body).Decode(&partyJson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pa := data.NewParty()
	pa.Size = partyJson.Size

	p, err := storages.Storage.GetPlayerById(userId)
	if err != nil {
		switch {
		case errors.Is(err, storages.ErrRecordNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	pa.AddPlayer(p)
	match_routines.PlayerJobsChannel <- p
}

func PlayerEnqueue(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	p, err := storages.Storage.GetPlayerById(userId)
	if err != nil {
		switch {
		case errors.Is(err, storages.ErrRecordNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	match_routines.PlayerJobsChannel <- p
}

func PlayerRunAi(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	// TODO handle user snakes
	snake := game.NewSnake(5, 5, 1, 0, []func(snake *game.Snake){
		func(snake *game.Snake) { snake.Move() },
		func(snake *game.Snake) { snake.Up() },
		func(snake *game.Snake) { snake.Move() },
		func(snake *game.Snake) { snake.Left() },
		func(snake *game.Snake) { snake.Move() },
		func(snake *game.Snake) { snake.Down() },
		func(snake *game.Snake) { snake.Move() },
		func(snake *game.Snake) { snake.Right() },
	})
	// TODO refactor inner cycles for optimization
out:
	for _, g := range game.Games {
		for _, p := range g.Party.Players {
			if p.Id == userId {
				g.Snakes[userId] = snake
				break out
			}
		}
	}
}
