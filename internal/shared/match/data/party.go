package data

import (
	"math/rand"
	"sync"
)

type Party struct {
	mux      sync.Mutex
	Id       string
	Players  []*Player
	AvgSkill int
	Size     int
}

func NewParty() Party {
	return Party{Id: randSeq(10), Players: []*Player{}, AvgSkill: 0}
}

func (party *Party) lock() {
	party.mux.Lock()
}

func (party *Party) unlock() {
	party.mux.Unlock()
}
func (party *Party) AddPlayer(player *Player) {
	player.Lock()
	if player.InParty {
		player.Unlock()
		return
	}
	player.Party = party
	player.InParty = true
	player.Unlock()

	party.lock()
	defer party.unlock()
	party.Players = append(party.Players, player)
	party.computeAvgSkill()
}
func (party *Party) RemovePlayer(p *Player) {
	result := make([]*Player, 0)
	party.lock()
	defer party.unlock()
	for _, pl := range party.Players {
		if pl.Name != p.Name {
			result = append(result, pl)
		}
	}
	party.Players = result
	party.computeAvgSkill()
}
func (party *Party) computeAvgSkill() {
	if len(party.Players) == 0 {
		party.AvgSkill = 0
		return
	}
	sum := 0
	for _, p := range party.Players {
		sum += p.Skill
	}
	party.AvgSkill = sum / len(party.Players)
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
