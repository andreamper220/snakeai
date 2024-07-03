package routines

import (
	gamedata "snakeai/internal/domain/game/data"
)

var GameJobsChannel = make(chan *gamedata.Game, 100)

func GameWorker() {
	for g := range GameJobsChannel {
		g.Update()
	}
}
