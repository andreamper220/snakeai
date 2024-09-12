package routines

import (
	gamedata "github.com/andreamper220/snakeai/internal/server/domain/game/data"
)

var GameJobsChannel = make(chan *gamedata.Game, 100)

func GameWorker() {
	for g := range GameJobsChannel {
		g.Update()
	}
}
