package game_test

import (
	"github.com/andreamper220/snakeai/internal/server/domain/game/data"
	matchdata "github.com/andreamper220/snakeai/internal/server/domain/match/data"
	"github.com/google/uuid"
	"time"
)

const (
	gameWidth  = 5
	gameHeight = 5
	initX      = 3
	initY      = 3
	initXTo    = 1
	initYTo    = 0
)

func (s *GameTestSuite) TestSnakeMove() {
	pa := matchdata.NewParty()
	g := data.NewGame(gameWidth, gameHeight, &pa)
	g.Food = &data.Food{
		Position: data.Point{
			X: 1,
			Y: 1,
		},
	}
	s.games.AddGame(g)

	tests := []struct {
		name     string
		commands []func(snake *data.Snake)
		x        int
		y        int
	}{
		{
			name: "move",
			commands: []func(snake *data.Snake){
				func(snake *data.Snake) { snake.Move() },
			},
			x: 4,
			y: 3,
		},
		{
			name: "right",
			commands: []func(snake *data.Snake){
				func(snake *data.Snake) { snake.Right() },
				func(snake *data.Snake) { snake.Move() },
			},
			x: 3,
			y: 4,
		},
		{
			name: "left",
			commands: []func(snake *data.Snake){
				func(snake *data.Snake) { snake.Left() },
				func(snake *data.Snake) { snake.Move() },
			},
			x: 3,
			y: 2,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			sn := data.NewSnake(initX, initY, initXTo, initYTo, tt.commands)
			userId := uuid.New()
			g.AddSnake(sn, userId)
			for range tt.commands {
				g.Update()
				time.Sleep(100 * time.Millisecond)
			}
			sn.Lock()
			s.Assert().Equal(tt.x, sn.Body[0].X)
			s.Assert().Equal(tt.y, sn.Body[0].Y)
			sn.Unlock()
			g.RemoveSnake(userId)
		})
	}

	s.games.RemoveGame(g)
}

func (s *GameTestSuite) TestSnakeRotation() {
	pa := matchdata.NewParty()
	g := data.NewGame(gameWidth, gameHeight, &pa)
	g.Food = &data.Food{
		Position: data.Point{
			X: 1,
			Y: 1,
		},
	}
	s.games.AddGame(g)

	tests := []struct {
		name       string
		command    func(snake *data.Snake)
		directions []data.Point
	}{
		{
			name:    "right",
			command: func(snake *data.Snake) { snake.Right() },
			directions: []data.Point{
				{X: 0, Y: 1},
				{X: -1, Y: 0},
				{X: 0, Y: -1},
				{X: 1, Y: 0},
			},
		},
		{
			name:    "left",
			command: func(snake *data.Snake) { snake.Left() },
			directions: []data.Point{
				{X: 0, Y: -1},
				{X: -1, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 0},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			sn := data.NewSnake(initX, initY, initXTo, initYTo, []func(snake *data.Snake){tt.command})
			userId := uuid.New()
			g.AddSnake(sn, userId)
			for _, direction := range tt.directions {
				g.Update()
				time.Sleep(100 * time.Millisecond)
				sn.Lock()
				s.Assert().Equal(direction.X, sn.Direction.X)
				s.Assert().Equal(direction.Y, sn.Direction.Y)
				sn.Unlock()
			}
			g.RemoveSnake(userId)
		})
	}

	s.games.RemoveGame(g)
}

func (s *GameTestSuite) TestSnakeEdgeCollision() {
	pa := matchdata.NewParty()
	g := data.NewGame(gameWidth, gameHeight, &pa)
	g.Food = &data.Food{
		Position: data.Point{
			X: 3,
			Y: 3,
		},
	}
	s.games.AddGame(g)

	tests := []struct {
		name  string
		initX int
		initY int
		xTo   int
		yTo   int
	}{
		{
			name:  "up",
			initX: 1,
			initY: 5,
			xTo:   0,
			yTo:   1,
		},
		{
			name:  "right",
			initX: 5,
			initY: 1,
			xTo:   1,
			yTo:   0,
		},
		{
			name:  "down",
			initX: 1,
			initY: 1,
			xTo:   0,
			yTo:   -1,
		},
		{
			name:  "left",
			initX: 1,
			initY: 1,
			xTo:   -1,
			yTo:   0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			sn := data.NewSnake(tt.initX, tt.initY, tt.xTo, tt.yTo, []func(snake *data.Snake){
				func(snake *data.Snake) { snake.Move() },
			})
			userId := uuid.New()
			g.AddSnake(sn, userId)
			time.Sleep(100 * time.Millisecond)
			g.Update()
			time.Sleep(100 * time.Millisecond)
			s.Assert().Equal(0, len(g.GetSnakes()))
		})
	}

	s.games.RemoveGame(g)
}

func (s *GameTestSuite) TestSnakeSnakeCollision() {
	user1 := s.AddNewUser("test@test.com")
	user2 := s.AddNewUser("test@test1.com")

	pa := matchdata.NewParty()
	g := data.NewGame(gameWidth, gameHeight, &pa)
	g.Food = &data.Food{
		Position: data.Point{
			X: 1,
			Y: 1,
		},
	}
	s.games.AddGame(g)

	sn1 := data.NewSnake(3, 3, 1, 0, []func(snake *data.Snake){
		func(snake *data.Snake) { snake.Move() },
	})
	sn2 := data.NewSnake(4, 3, -1, 0, []func(snake *data.Snake){
		func(snake *data.Snake) { snake.Move() },
	})
	s.T().Log(user1.Id, user2.Id)
	g.AddSnake(sn1, user1.Id)
	g.AddSnake(sn2, user2.Id)
	g.Update()
	time.Sleep(200 * time.Millisecond)
	s.Assert().Equal(0, len(g.Snakes.Data))

	s.games.RemoveGame(g)
}

func (s *GameTestSuite) TestSnakeFoodEating() {
	user := s.AddNewUser("test@test.com")

	pa := matchdata.NewParty()
	g := data.NewGame(gameWidth, gameHeight, &pa)
	g.Food = &data.Food{
		Position: data.Point{
			X: 4,
			Y: 3,
		},
	}
	s.games.AddGame(g)

	sn := data.NewSnake(initX, initY, initXTo, initYTo, []func(snake *data.Snake){
		func(snake *data.Snake) { snake.Move() },
	})
	g.AddSnake(sn, user.Id)
	g.Update()
	time.Sleep(100 * time.Millisecond)
	g.Update()
	time.Sleep(100 * time.Millisecond)
	sn.Lock()
	s.Assert().Equal(2, len(sn.Body))
	sn.Unlock()

	s.games.RemoveGame(g)
}
