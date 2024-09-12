package post_handlers

import (
	"encoding/json"
	"errors"
	"github.com/andreamper220/snakeai/internal/server/domain/game/data"
	gamejson "github.com/andreamper220/snakeai/internal/server/domain/game/json"
	matchdata "github.com/andreamper220/snakeai/internal/server/domain/match/data"
	matchjson "github.com/andreamper220/snakeai/internal/server/domain/match/json"
	matchroutines "github.com/andreamper220/snakeai/internal/server/domain/match/routines"
	"github.com/andreamper220/snakeai/internal/server/infrastructure/storages"
	"github.com/andreamper220/snakeai/pkg/logger"
	"github.com/google/uuid"
	"net/http"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pa.AddPlayer(p)
	matchroutines.PlayerJobsChannel <- p
}

func PlayerEnqueue(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	p, err := storages.Storage.GetPlayerById(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	ai, err := data.GenerateAiFunctions(aiJson.Ai)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	snake := data.NewSnake(aiJson.X, aiJson.Y, aiJson.XTo, aiJson.YTo, ai)
out:
	for _, g := range data.CurrentGames.Games {
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
					g.Lock()
					g.Scores[userId] = pl.Skill
					g.Unlock()
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
