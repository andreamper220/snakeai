package post_handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"snake_ai/pkg/logger"

	gamedata "snake_ai/internal/domain/game/data"
	gamejson "snake_ai/internal/domain/game/json"
	matchdata "snake_ai/internal/domain/match/data"
	matchjson "snake_ai/internal/domain/match/json"
	matchroutines "snake_ai/internal/domain/match/routines"
	"snake_ai/internal/infrastructure/storages"
)

func PlayerPartyEnqueue(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var partyJson matchjson.PartyJson
	if err := json.NewDecoder(r.Body).Decode(&partyJson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pa := matchdata.NewParty()
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
	matchroutines.PlayerJobsChannel <- p
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

	matchroutines.PlayerJobsChannel <- p
}

func PlayerRunAi(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var aiJson gamejson.AiJson
	if err := json.NewDecoder(r.Body).Decode(&aiJson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	snake := gamedata.NewSnake(aiJson.X, aiJson.Y, aiJson.XTo, aiJson.YTo, gamedata.GenerateAiFunctions(aiJson.Ai))
out:
	for _, g := range gamedata.CurrentGames.Games {
		for _, p := range g.Party.Players {
			if p.Id == userId {
				g.AddSnake(snake, userId)
				pl, err := storages.Storage.GetPlayerById(userId)
				if err != nil {
					g.Scores[userId] += 0
					logger.Log.Error(err.Error())
					switch {
					case errors.Is(err, storages.ErrRecordNotFound):
						http.Error(w, err.Error(), http.StatusNotFound)
					default:
						http.Error(w, err.Error(), http.StatusBadRequest)
					}
					return
				} else {
					g.Scores[userId] = pl.Skill
					w.Header().Set("Content-Type", "application/json")
					if err = json.NewEncoder(w).Encode(pl); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
				break out
			}
		}
	}
}
