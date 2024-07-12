package game_test

import (
	"github.com/google/uuid"
	"time"

	gamedata "github.com/andreamper220/snakeai/internal/domain/game/data"
	matchdata "github.com/andreamper220/snakeai/internal/domain/match/data"
)

type testCondition struct {
	condition gamedata.ObstacleCondition
	distance  int
}

var conditions = []testCondition{
	{
		condition: gamedata.Equal,
		distance:  1,
	},
	{
		condition: gamedata.NotEqual,
		distance:  0,
	},
	{
		condition: gamedata.GreaterThan,
		distance:  0,
	},
	{
		condition: gamedata.LessThan,
		distance:  2,
	},
	{
		condition: gamedata.GreaterOrEqual,
		distance:  1,
	},
	{
		condition: gamedata.LessOrEqual,
		distance:  1,
	},
}

func (s *GameTestSuite) TestSnakeIfEdge() {
	pa := matchdata.NewParty()
	g := gamedata.NewGame(gameWidth, gameHeight, &pa)
	s.games.AddGame(g)

	type subtest struct {
		name  string
		initX int
		initY int
		xTo   int
		yTo   int
	}

	// test on 'move;' command
	tests := []struct {
		name     string
		subtests []subtest
	}{
		{
			name: "right",
			subtests: []subtest{
				{
					name:  "bottom",
					initX: 3,
					initY: 4,
					xTo:   1,
					yTo:   0,
				},
				{
					name:  "right",
					initX: 4,
					initY: 3,
					xTo:   0,
					yTo:   -1,
				},
				{
					name:  "top",
					initX: 3,
					initY: 2,
					xTo:   -1,
					yTo:   0,
				},
				{
					name:  "left",
					initX: 2,
					initY: 3,
					xTo:   0,
					yTo:   1,
				},
			},
		},
		{
			name: "left",
			subtests: []subtest{
				{
					name:  "bottom",
					initX: 3,
					initY: 4,
					xTo:   -1,
					yTo:   0,
				},
				{
					name:  "right",
					initX: 4,
					initY: 3,
					xTo:   0,
					yTo:   1,
				},
				{
					name:  "top",
					initX: 3,
					initY: 2,
					xTo:   1,
					yTo:   0,
				},
				{
					name:  "left",
					initX: 2,
					initY: 3,
					xTo:   0,
					yTo:   -1,
				},
			},
		},
		{
			name: "forward",
			subtests: []subtest{
				{
					name:  "bottom",
					initX: 3,
					initY: 4,
					xTo:   0,
					yTo:   1,
				},
				{
					name:  "right",
					initX: 4,
					initY: 3,
					xTo:   1,
					yTo:   0,
				},
				{
					name:  "top",
					initX: 3,
					initY: 2,
					xTo:   0,
					yTo:   -1,
				},
				{
					name:  "left",
					initX: 2,
					initY: 3,
					xTo:   -1,
					yTo:   0,
				},
			},
		},
	}

	for _, tt := range tests {
		for _, ttt := range tt.subtests {
			for _, cond := range conditions {
				s.Run(tt.name+"_"+ttt.name+"_"+string(cond.condition), func() {
					cond := gamedata.AiCondition{
						ObstacleType:      gamedata.ObstacleEdge,
						ObstacleDirection: gamedata.ObstacleDirection(tt.name),
						ObstacleCondition: cond.condition,
						ObstacleDistance:  cond.distance,
					}
					sn := gamedata.NewSnake(ttt.initX, ttt.initY, ttt.xTo, ttt.yTo, []func(snake *gamedata.Snake){
						func(snake *gamedata.Snake) { snake.DoIf(cond, 1) },
						func(snake *gamedata.Snake) { snake.Move() },
					})
					userId := uuid.New()
					g.AddSnake(sn, userId)
					time.Sleep(100 * time.Millisecond)
					g.Update()
					time.Sleep(100 * time.Millisecond)
					g.Update()
					time.Sleep(100 * time.Millisecond)
					s.Assert().Equal(ttt.initX+ttt.xTo, sn.Body[0].X)
					s.Assert().Equal(ttt.initY+ttt.yTo, sn.Body[0].Y)
					g.RemoveSnake(userId)
				})
			}
		}
	}

	s.games.RemoveGame(g)
}

func (s *GameTestSuite) TestSnakeIfFood() {
	pa := matchdata.NewParty()
	g := gamedata.NewGame(gameWidth, gameHeight, &pa)
	g.Food = &gamedata.Food{
		Position: gamedata.Point{
			X: 3,
			Y: 3,
		},
	}
	s.games.AddGame(g)

	type subtest struct {
		name  string
		initX int
		initY int
		xTo   int
		yTo   int
	}

	// test on 'move;' command
	tests := []struct {
		name     string
		subtests []subtest
	}{
		{
			name: "right",
			subtests: []subtest{
				{
					name:  "bottom",
					initX: 3,
					initY: 1,
					xTo:   1,
					yTo:   0,
				},
				{
					name:  "right",
					initX: 1,
					initY: 3,
					xTo:   0,
					yTo:   -1,
				},
				{
					name:  "top",
					initX: 3,
					initY: 5,
					xTo:   -1,
					yTo:   0,
				},
				{
					name:  "left",
					initX: 5,
					initY: 3,
					xTo:   0,
					yTo:   1,
				},
			},
		},
		{
			name: "left",
			subtests: []subtest{
				{
					name:  "bottom",
					initX: 3,
					initY: 1,
					xTo:   -1,
					yTo:   0,
				},
				{
					name:  "right",
					initX: 1,
					initY: 3,
					xTo:   0,
					yTo:   1,
				},
				{
					name:  "top",
					initX: 3,
					initY: 5,
					xTo:   1,
					yTo:   0,
				},
				{
					name:  "left",
					initX: 5,
					initY: 3,
					xTo:   0,
					yTo:   -1,
				},
			},
		},
		{
			name: "forward",
			subtests: []subtest{
				{
					name:  "bottom",
					initX: 3,
					initY: 1,
					xTo:   0,
					yTo:   1,
				},
				{
					name:  "right",
					initX: 1,
					initY: 3,
					xTo:   1,
					yTo:   0,
				},
				{
					name:  "top",
					initX: 3,
					initY: 5,
					xTo:   0,
					yTo:   -1,
				},
				{
					name:  "left",
					initX: 5,
					initY: 3,
					xTo:   -1,
					yTo:   0,
				},
			},
		},
	}

	for _, tt := range tests {
		for _, ttt := range tt.subtests {
			for _, cond := range conditions {
				s.Run(tt.name+"_"+ttt.name+"_"+string(cond.condition), func() {
					cond := gamedata.AiCondition{
						ObstacleType:      gamedata.ObstacleFood,
						ObstacleDirection: gamedata.ObstacleDirection(tt.name),
						ObstacleCondition: cond.condition,
						ObstacleDistance:  cond.distance,
					}
					sn := gamedata.NewSnake(ttt.initX, ttt.initY, ttt.xTo, ttt.yTo, []func(snake *gamedata.Snake){
						func(snake *gamedata.Snake) { snake.DoIf(cond, 1) },
						func(snake *gamedata.Snake) { snake.Move() },
					})
					userId := uuid.New()
					g.AddSnake(sn, userId)
					time.Sleep(100 * time.Millisecond)
					g.Update()
					time.Sleep(100 * time.Millisecond)
					g.Update()
					time.Sleep(100 * time.Millisecond)
					s.Assert().Equal(ttt.initX+ttt.xTo, sn.Body[0].X)
					s.Assert().Equal(ttt.initY+ttt.yTo, sn.Body[0].Y)
					g.RemoveSnake(userId)
				})
			}
		}
	}

	s.games.RemoveGame(g)
}

func (s *GameTestSuite) TestSnakeIfSnake() {
	initX2 := 3
	initY2 := 3

	pa := matchdata.NewParty()
	g := gamedata.NewGame(gameWidth, gameHeight, &pa)
	s.games.AddGame(g)

	sn2 := gamedata.NewSnake(initX2, initY2, 1, 0, []func(snake *gamedata.Snake){})
	userId2 := uuid.New()
	g.AddSnake(sn2, userId2)

	type subtest struct {
		name  string
		initX int
		initY int
		xTo   int
		yTo   int
	}

	// test on 'move;' command
	tests := []struct {
		name     string
		subtests []subtest
	}{
		{
			name: "right",
			subtests: []subtest{
				{
					name:  "bottom",
					initX: 3,
					initY: 1,
					xTo:   1,
					yTo:   0,
				},
				{
					name:  "right",
					initX: 1,
					initY: 3,
					xTo:   0,
					yTo:   -1,
				},
				{
					name:  "top",
					initX: 3,
					initY: 5,
					xTo:   -1,
					yTo:   0,
				},
				{
					name:  "left",
					initX: 5,
					initY: 3,
					xTo:   0,
					yTo:   1,
				},
			},
		},
		{
			name: "left",
			subtests: []subtest{
				{
					name:  "bottom",
					initX: 3,
					initY: 1,
					xTo:   -1,
					yTo:   0,
				},
				{
					name:  "right",
					initX: 1,
					initY: 3,
					xTo:   0,
					yTo:   1,
				},
				{
					name:  "top",
					initX: 3,
					initY: 5,
					xTo:   1,
					yTo:   0,
				},
				{
					name:  "left",
					initX: 5,
					initY: 3,
					xTo:   0,
					yTo:   -1,
				},
			},
		},
		{
			name: "forward",
			subtests: []subtest{
				{
					name:  "bottom",
					initX: 3,
					initY: 1,
					xTo:   0,
					yTo:   1,
				},
				{
					name:  "right",
					initX: 1,
					initY: 3,
					xTo:   1,
					yTo:   0,
				},
				{
					name:  "top",
					initX: 3,
					initY: 5,
					xTo:   0,
					yTo:   -1,
				},
				{
					name:  "left",
					initX: 5,
					initY: 3,
					xTo:   -1,
					yTo:   0,
				},
			},
		},
	}

	for _, tt := range tests {
		for _, ttt := range tt.subtests {
			for _, cond := range conditions {
				s.Run(tt.name+"_"+ttt.name+"_"+string(cond.condition), func() {
					cond := gamedata.AiCondition{
						ObstacleType:      gamedata.ObstacleSnake,
						ObstacleDirection: gamedata.ObstacleDirection(tt.name),
						ObstacleCondition: cond.condition,
						ObstacleDistance:  cond.distance,
					}
					sn1 := gamedata.NewSnake(ttt.initX, ttt.initY, ttt.xTo, ttt.yTo, []func(snake *gamedata.Snake){
						func(snake *gamedata.Snake) { snake.DoIf(cond, 1) },
						func(snake *gamedata.Snake) { snake.Move() },
					})
					userId1 := uuid.New()
					g.AddSnake(sn1, userId1)

					time.Sleep(100 * time.Millisecond)
					g.Update()
					time.Sleep(100 * time.Millisecond)
					g.Update()
					time.Sleep(100 * time.Millisecond)
					s.Assert().Equal(ttt.initX+ttt.xTo, sn1.Body[0].X)
					s.Assert().Equal(ttt.initY+ttt.yTo, sn1.Body[0].Y)
					g.RemoveSnake(userId1)
				})
			}
		}
	}

	s.games.RemoveGame(g)
}

func (s *GameTestSuite) TestSnakeElseIf() {
	pa := matchdata.NewParty()
	g := gamedata.NewGame(gameWidth, gameHeight, &pa)
	s.games.AddGame(g)

	condIf := gamedata.AiCondition{
		ObstacleType:      gamedata.ObstacleEdge,
		ObstacleDirection: gamedata.Right,
		ObstacleCondition: gamedata.Equal,
		ObstacleDistance:  0,
	}
	condElseIf := gamedata.AiCondition{
		ObstacleType:      gamedata.ObstacleEdge,
		ObstacleDirection: gamedata.Right,
		ObstacleCondition: gamedata.Equal,
		ObstacleDistance:  1,
	}

	s.Run("if", func() {
		sn := gamedata.NewSnake(5, 3, 0, -1, []func(snake *gamedata.Snake){
			func(snake *gamedata.Snake) { snake.DoIf(condIf, 1) },
			func(snake *gamedata.Snake) { snake.Move() },
			func(snake *gamedata.Snake) { snake.DoElseIf(condElseIf, 2) },
			func(snake *gamedata.Snake) { snake.Left() },
			func(snake *gamedata.Snake) { snake.Move() },
		})
		userId := uuid.New()
		g.AddSnake(sn, userId)

		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		s.Assert().Equal(5, sn.Body[0].X)
		s.Assert().Equal(2, sn.Body[0].Y)
		g.RemoveSnake(userId)
	})

	s.Run("else_if", func() {
		sn := gamedata.NewSnake(4, 3, 0, -1, []func(snake *gamedata.Snake){
			func(snake *gamedata.Snake) { snake.DoIf(condIf, 1) },
			func(snake *gamedata.Snake) { snake.Move() },
			func(snake *gamedata.Snake) { snake.DoElseIf(condElseIf, 2) },
			func(snake *gamedata.Snake) { snake.Left() },
			func(snake *gamedata.Snake) { snake.Move() },
		})
		userId := uuid.New()
		g.AddSnake(sn, userId)

		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		s.Assert().Equal(3, sn.Body[0].X)
		s.Assert().Equal(3, sn.Body[0].Y)
		g.RemoveSnake(userId)
	})

	s.games.RemoveGame(g)
}

func (s *GameTestSuite) TestSnakeElse() {
	pa := matchdata.NewParty()
	g := gamedata.NewGame(gameWidth, gameHeight, &pa)
	s.games.AddGame(g)

	condIf := gamedata.AiCondition{
		ObstacleType:      gamedata.ObstacleEdge,
		ObstacleDirection: gamedata.Right,
		ObstacleCondition: gamedata.Equal,
		ObstacleDistance:  0,
	}
	condElseIf := gamedata.AiCondition{
		ObstacleType:      gamedata.ObstacleEdge,
		ObstacleDirection: gamedata.Right,
		ObstacleCondition: gamedata.Equal,
		ObstacleDistance:  1,
	}

	s.Run("if", func() {
		sn := gamedata.NewSnake(5, 3, 0, -1, []func(snake *gamedata.Snake){
			func(snake *gamedata.Snake) { snake.DoIf(condIf, 1) },
			func(snake *gamedata.Snake) { snake.Move() },
			func(snake *gamedata.Snake) { snake.DoElseIf(condElseIf, 2) },
			func(snake *gamedata.Snake) { snake.Left() },
			func(snake *gamedata.Snake) { snake.Move() },
			func(snake *gamedata.Snake) { snake.DoElse(2) },
			func(snake *gamedata.Snake) { snake.Right() },
			func(snake *gamedata.Snake) { snake.Move() },
		})
		userId := uuid.New()
		g.AddSnake(sn, userId)

		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		s.Assert().Equal(5, sn.Body[0].X)
		s.Assert().Equal(2, sn.Body[0].Y)
		g.RemoveSnake(userId)
	})

	s.Run("else_if", func() {
		sn := gamedata.NewSnake(4, 3, 0, -1, []func(snake *gamedata.Snake){
			func(snake *gamedata.Snake) { snake.DoIf(condIf, 1) },
			func(snake *gamedata.Snake) { snake.Move() },
			func(snake *gamedata.Snake) { snake.DoElseIf(condElseIf, 2) },
			func(snake *gamedata.Snake) { snake.Left() },
			func(snake *gamedata.Snake) { snake.Move() },
			func(snake *gamedata.Snake) { snake.DoElse(2) },
			func(snake *gamedata.Snake) { snake.Right() },
			func(snake *gamedata.Snake) { snake.Move() },
		})
		userId := uuid.New()
		g.AddSnake(sn, userId)

		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		s.Assert().Equal(3, sn.Body[0].X)
		s.Assert().Equal(3, sn.Body[0].Y)
		g.RemoveSnake(userId)
	})

	s.Run("else", func() {
		sn := gamedata.NewSnake(3, 3, 0, -1, []func(snake *gamedata.Snake){
			func(snake *gamedata.Snake) { snake.DoIf(condIf, 1) },
			func(snake *gamedata.Snake) { snake.Move() },
			func(snake *gamedata.Snake) { snake.DoElseIf(condElseIf, 2) },
			func(snake *gamedata.Snake) { snake.Left() },
			func(snake *gamedata.Snake) { snake.Move() },
			func(snake *gamedata.Snake) { snake.DoElse(2) },
			func(snake *gamedata.Snake) { snake.Right() },
			func(snake *gamedata.Snake) { snake.Move() },
		})
		userId := uuid.New()
		g.AddSnake(sn, userId)

		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		g.Update()
		time.Sleep(100 * time.Millisecond)
		s.Assert().Equal(4, sn.Body[0].X)
		s.Assert().Equal(3, sn.Body[0].Y)
		g.RemoveSnake(userId)
	})

	s.games.RemoveGame(g)
}
