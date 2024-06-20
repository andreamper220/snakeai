package routines

import (
	"snake_ai/internal/server/ai/data"
)

var GameJobsChannel = make(chan *data.Game, 100)

func GameWorker() {
	for g := range GameJobsChannel {
		g.Update()
	}
}
