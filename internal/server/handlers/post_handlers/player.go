package post_handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"snake_ai/internal/logger"

	game "snake_ai/internal/server/ai/data"
	aijs "snake_ai/internal/server/ai/json"
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
	pa.Width = partyJson.Width
	pa.Height = partyJson.Height

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
	var aiJson aijs.AiJson
	if err := json.NewDecoder(r.Body).Decode(&aiJson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	snake := game.NewSnake(aiJson.X, aiJson.Y, aiJson.XTo, aiJson.YTo, game.GenerateAiFunctions(aiJson.Ai), userId)
	// TODO refactor inner cycles for optimization
out:
	for _, g := range game.CurrentGames.Games {
		if g != nil && g.Party != nil {
			for _, p := range g.Party.Players {
				if p.Id == userId {
					g.AddSnake(snake, userId)
					player, err := storages.Storage.GetPlayerById(userId)
					if err != nil {
						g.Scores[userId] += 0
						logger.Log.Error(err.Error())
					} else {
						g.Scores[userId] = player.Skill
					}
					break out
				}
			}
		}
	}
}
