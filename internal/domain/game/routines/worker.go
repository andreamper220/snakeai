package routines

import (
	gamedata "snake_ai/internal/domain/game/data"
)

var GameJobsChannel = make(chan *gamedata.Game, 100)

func GameWorker() {
	for g := range GameJobsChannel {
		g.Update()
	}
}
