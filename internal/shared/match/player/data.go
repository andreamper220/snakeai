package player

import (
	"sync"

	"snake_ai/internal/shared/match/party"
)

type Player struct {
	mux        sync.Mutex
	name       string
	skill      int
	foundParty bool
	delta      int
	party      *party.Party
	inParty    bool
	inProcess  bool
}
