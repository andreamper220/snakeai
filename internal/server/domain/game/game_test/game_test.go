package game_test

import (
	"context"
	editorserver "github.com/andreamper220/snakeai/internal/editor/application"
	"github.com/andreamper220/snakeai/internal/server/application"
	gamedata "github.com/andreamper220/snakeai/internal/server/domain/game/data"
	matchdata "github.com/andreamper220/snakeai/internal/server/domain/match/data"
	grpcclients "github.com/andreamper220/snakeai/internal/server/infrastructure/grpc"
	pb "github.com/andreamper220/snakeai/proto"
	"github.com/stretchr/testify/assert"
	"os"
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

	t.Run("new with walls", func(t *testing.T) {
		editorServer := editorserver.InitGRPCServer()
		go func() {
			assert.NoError(t, editorserver.Run(50051, true))
		}()
		time.Sleep(500 * time.Millisecond)

		assert.NoError(t, os.Setenv("EDITOR_ADDRESS", "0.0.0.0:50051"))
		application.ParseFlags()
		editorConn, err := application.ConnectEditor()
		assert.NoError(t, err)
		grpcclients.EditorClient = pb.NewEditorClient(editorConn)
		defer editorServer.Stop()

		mapResponse, err := grpcclients.EditorClient.SaveMap(context.Background(), &pb.SaveMapRequest{
			Struct: &pb.MapStruct{
				Width:  gameWidth,
				Height: gameHeight,
				Obstacles: []*pb.Obstacle{
					{
						Cx: 0,
						Cy: 0,
					},
					{
						Cx: 0,
						Cy: 1,
					},
					{
						Cx: 0,
						Cy: 2,
					},
					{
						Cx: 0,
						Cy: 3,
					},
					{
						Cx: 0,
						Cy: 4,
					},
					{
						Cx: 1,
						Cy: 0,
					},
					{
						Cx: 1,
						Cy: 1,
					},
					{
						Cx: 1,
						Cy: 2,
					},
					{
						Cx: 1,
						Cy: 3,
					},
					{
						Cx: 1,
						Cy: 4,
					},
					{
						Cx: 2,
						Cy: 0,
					},
					{
						Cx: 2,
						Cy: 1,
					},
					{
						Cx: 2,
						Cy: 2,
					},
					{
						Cx: 2,
						Cy: 3,
					},
					{
						Cx: 2,
						Cy: 4,
					},
					{
						Cx: 3,
						Cy: 0,
					},
					{
						Cx: 3,
						Cy: 1,
					},
					{
						Cx: 3,
						Cy: 2,
					},
					{
						Cx: 3,
						Cy: 3,
					},
					{
						Cx: 3,
						Cy: 4,
					},
					{
						Cx: 4,
						Cy: 0,
					},
					{
						Cx: 4,
						Cy: 1,
					},
					{
						Cx: 4,
						Cy: 2,
					},
					{
						Cx: 4,
						Cy: 3,
					},
				},
			},
		})
		assert.NoError(t, err)

		pa := matchdata.NewParty()
		pa.MapId = mapResponse.Map.Id
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
