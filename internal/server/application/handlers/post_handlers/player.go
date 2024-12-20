package post_handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/andreamper220/snakeai/internal/server/domain/ws"
	"github.com/google/uuid"
	"net/http"
	"time"

	"github.com/andreamper220/snakeai/internal/server/domain/game/data"
	gamejson "github.com/andreamper220/snakeai/internal/server/domain/game/json"
	matchdata "github.com/andreamper220/snakeai/internal/server/domain/match/data"
	matchjson "github.com/andreamper220/snakeai/internal/server/domain/match/json"
	matchroutines "github.com/andreamper220/snakeai/internal/server/domain/match/routines"
	grpcclients "github.com/andreamper220/snakeai/internal/server/infrastructure/grpc"
	"github.com/andreamper220/snakeai/internal/server/infrastructure/storages"
	"github.com/andreamper220/snakeai/pkg/logger"
	pb "github.com/andreamper220/snakeai/proto"
)

var ErrIncorrectFieldParams = errors.New("incorrect field parameters")

func CreateOrJoinParty(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var partyJson matchjson.PartyJson
	if err := json.NewDecoder(r.Body).Decode(&partyJson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if (partyJson.Width < 5 || partyJson.Width > 30) || (partyJson.Height < 5 || partyJson.Height > 30) || partyJson.Size > 10 {
		http.Error(w, ErrIncorrectFieldParams.Error(), http.StatusBadRequest)
		return
	}

	mapID := ""
	if len(partyJson.Obstacles) > 0 {
		obstacles := make([]*pb.Obstacle, len(partyJson.Obstacles))
		for i := 0; i < len(obstacles); i++ {
			obstacle := &pb.Obstacle{
				Cx: partyJson.Obstacles[i][0],
				Cy: partyJson.Obstacles[i][1],
			}
			obstacles[i] = obstacle
		}
		requestMap := &pb.SaveMapRequest{
			Struct: &pb.MapStruct{
				Width:     int32(partyJson.Width),
				Height:    int32(partyJson.Height),
				Obstacles: obstacles,
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		responseMap, err := grpcclients.EditorClient.SaveMap(ctx, requestMap)
		if err != nil {
			logger.Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mapID = responseMap.GetMap().GetId()
	}

	pa := matchdata.NewParty()
	pa.Size = partyJson.Size
	pa.Width = partyJson.Width
	pa.Height = partyJson.Height
	if partyJson.ToUseById != "" {
		pa.ToConnectById = true
	}
	pa.MapId = mapID

	p, err := storages.Storage.GetPlayerById(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pa.AddPlayer(p)
	matchroutines.PlayerJobsChannel <- p

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&pa); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RestoreParty(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	if game := data.GetGameByPlayer(userId); game != nil {
		err := ws.Connections.WriteJSON(userId, game.Party)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func JoinParty(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	p, err := storages.Storage.GetPlayerById(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var playerPartyJson matchjson.PlayerPartyJson
	if err := json.NewDecoder(r.Body).Decode(&playerPartyJson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.PartyId = playerPartyJson.PartyId

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

func PlayerMapCheck(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var partyJson matchjson.PartyJson
	if err := json.NewDecoder(r.Body).Decode(&partyJson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	obstacles := make([]*pb.Obstacle, len(partyJson.Obstacles))
	if len(partyJson.Obstacles) > 0 {
		for i := 0; i < len(obstacles); i++ {
			obstacle := &pb.Obstacle{
				Cx: partyJson.Obstacles[i][0],
				Cy: partyJson.Obstacles[i][1],
			}
			obstacles[i] = obstacle
		}
	}

	requestMap := &pb.CheckMapRequest{
		Struct: &pb.MapStruct{
			Width:     int32(partyJson.Width),
			Height:    int32(partyJson.Height),
			Obstacles: obstacles,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := grpcclients.EditorClient.CheckMap(ctx, requestMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
