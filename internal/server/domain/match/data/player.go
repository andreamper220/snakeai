package data

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

// Player represent a thread-safe object with ID, name, skill and party match-making fields.
type Player struct {
	mux        sync.Mutex
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Skill      int       `json:"skill"`
	SkillDelta int       `json:"-"`
	Party      *Party    `json:"-"`
	InParty    bool      `json:"-"`
	InProcess  bool      `json:"-"`
	PartyId    string    `json:"-"`
}

// NewPlayer creates a player with skill delta = 2.
func NewPlayer() Player {
	return Player{SkillDelta: 2}
}

func (player *Player) Lock() {
	player.mux.Lock()
}
func (player *Player) Unlock() {
	player.mux.Unlock()
}
func (player *Player) FindParty() (*Party, error) {
	if player.InProcess {
		return nil, errors.New("player already searching for party")
	}

	player.InProcess = true
	parties := CurrentParties.GetParties()
	if player.PartyId != "" {
		for _, party := range parties {
			if player.PartyId == party.Id {
				return party, nil
			}
		}
	} else {
		goodParties := player.getGoodParties(parties)
		return findBestParty(goodParties), nil
	}
	return nil, errors.New("could not find party")
}
func (player *Player) getGoodParties(parties []*Party) []*Party {
	var goodParties []*Party
	for _, pa := range parties {
		if !pa.ToConnectById && isPartyGoodForPlayer(player, pa) {
			goodParties = append(goodParties, pa)
		}
	}
	return goodParties
}
func isPartyGoodForPlayer(player *Player, party *Party) bool {
	// find party with [as-d; as+d] interval player skill satisfying
	ps := player.Skill
	as := party.AvgSkill
	d := player.SkillDelta

	return ps >= (as-d) && ps <= (as+d)
}
func findBestParty(parties []*Party) *Party {
	if len(parties) == 0 {
		return nil
	}

	// find the most full party
	bestParty := parties[0]
	maxLen := len(bestParty.Players)
	for _, p := range parties {
		if len(p.Players) > maxLen {
			maxLen = len(p.Players)
			bestParty = p
		}
	}
	return bestParty
}
