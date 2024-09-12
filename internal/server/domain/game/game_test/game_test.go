package game_test

import (
	gamedata "github.com/andreamper220/snakeai/internal/server/domain/game/data"
	matchdata "github.com/andreamper220/snakeai/internal/server/domain/match/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddGame(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		pa := matchdata.NewParty()
		g := gamedata.NewGame(gameWidth, gameHeight, &pa)
		gamedata.CurrentGames.AddGame(g)
		for _, gg := range gamedata.CurrentGames.Games {
			if gg.Id == g.Id {
				return
			}
		}
		t.Fail()
	})

	t.Run("existing", func(t *testing.T) {
		pa := matchdata.NewParty()
		g := gamedata.NewGame(gameWidth, gameHeight, &pa)
		gamedata.CurrentGames.AddGame(g)
		gamedata.CurrentGames.AddGame(g)
		games := gamedata.CurrentGames.GetGames()
		count := 0
		for _, gg := range games {
			if gg.Id == g.Id {
				count++
			}
		}
		assert.Equal(t, 1, count)
	})
}
