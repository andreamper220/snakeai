package game_test

import (
	"github.com/stretchr/testify/assert"
	gamedata "snake_ai/internal/domain/game/data"
	matchdata "snake_ai/internal/domain/match/data"
	"testing"
	"time"
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
		go gamedata.CurrentGames.AddGame(g)
		go gamedata.CurrentGames.AddGame(g)
		time.Sleep(1 * time.Second)
		count := 0
		for _, gg := range gamedata.CurrentGames.Games {
			if gg.Id == g.Id {
				count++
			}
		}
		assert.Equal(t, 1, count)
	})
}
