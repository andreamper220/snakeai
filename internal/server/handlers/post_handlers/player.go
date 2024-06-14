package post_handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"

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
	routines.PlayerJobsChannel <- p
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

	routines.PlayerJobsChannel <- p
}
