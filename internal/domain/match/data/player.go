package data

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

type Player struct {
	mux        sync.Mutex
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Skill      int       `json:"skill"`
	SkillDelta int       `json:"-"`
	Party      *Party    `json:"-"`
	InParty    bool      `json:"-"`
	InProcess  bool      `json:"-"`
}

func NewPlayer() Player {
	return Player{SkillDelta: 2}
}

func (player *Player) Lock() {
	player.mux.Lock()
}
func (player *Player) Unlock() {
	player.mux.Unlock()
}
func (player *Player) FindParty(parties []*Party) (*Party, error) {
	if player.InProcess {
		return nil, errors.New("player already searching for party")
	}

	player.InProcess = true
	goodParties := player.getGoodParties(parties)
	return findBestParty(goodParties), nil
}
func (player *Player) getGoodParties(parties []*Party) []*Party {
	var goodParties []*Party
	for _, pa := range parties {
		if isPartyGoodForPlayer(player, pa) {
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
