package game_test

import (
	"github.com/andreamper220/snakeai/internal/server/domain/game/data"
	matchdata "github.com/andreamper220/snakeai/internal/server/domain/match/data"
	"github.com/google/uuid"
	"time"
)

type testCondition struct {
	condition data.ObstacleCondition
	distance  int
}

var conditions = []testCondition{
	{
		condition: data.Equal,
		distance:  1,
	},
	{
		condition: data.NotEqual,
		distance:  0,
	},
	{
		condition: data.GreaterThan,
		distance:  0,
	},
	{
		condition: data.LessThan,
		distance:  2,
	},
	{
		condition: data.GreaterOrEqual,
		distance:  1,
	},
	{
		condition: data.LessOrEqual,
		distance:  1,
	},
}

func (s *GameTestSuite) TestSnakeIfEdge() {
	pa := matchdata.NewParty()
	g := data.NewGame(gameWidth, gameHeight, &pa)
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
					cond := data.AiConditions{
						Condition: data.AiCondition{
							ObstacleType:      data.ObstacleEdge,
							ObstacleDirection: data.ObstacleDirection(tt.name),
							ObstacleCondition: cond.condition,
							ObstacleDistance:  cond.distance,
						},
					}
					sn := data.NewSnake(ttt.initX, ttt.initY, ttt.xTo, ttt.yTo, []func(snake *data.Snake){
						func(snake *data.Snake) { snake.DoIf(cond, 1) },
						func(snake *data.Snake) { snake.Move() },
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
	g := data.NewGame(gameWidth, gameHeight, &pa)
	g.Food = &data.Food{
		Position: data.Point{
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
					cond := data.AiConditions{
						Condition: data.AiCondition{
							ObstacleType:      data.ObstacleFood,
							ObstacleDirection: data.ObstacleDirection(tt.name),
							ObstacleCondition: cond.condition,
							ObstacleDistance:  cond.distance,
						},
					}
					sn := data.NewSnake(ttt.initX, ttt.initY, ttt.xTo, ttt.yTo, []func(snake *data.Snake){
						func(snake *data.Snake) { snake.DoIf(cond, 1) },
						func(snake *data.Snake) { snake.Move() },
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
	g := data.NewGame(gameWidth, gameHeight, &pa)
	s.games.AddGame(g)

	sn2 := data.NewSnake(initX2, initY2, 1, 0, []func(snake *data.Snake){})
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
					cond := data.AiConditions{
						Condition: data.AiCondition{
							ObstacleType:      data.ObstacleSnake,
							ObstacleDirection: data.ObstacleDirection(tt.name),
							ObstacleCondition: cond.condition,
							ObstacleDistance:  cond.distance,
						},
					}
					sn1 := data.NewSnake(ttt.initX, ttt.initY, ttt.xTo, ttt.yTo, []func(snake *data.Snake){
						func(snake *data.Snake) { snake.DoIf(cond, 1) },
						func(snake *data.Snake) { snake.Move() },
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
	g := data.NewGame(gameWidth, gameHeight, &pa)
	s.games.AddGame(g)

	condIf := data.AiConditions{
		Condition: data.AiCondition{
			ObstacleType:      data.ObstacleEdge,
			ObstacleDirection: data.Right,
			ObstacleCondition: data.Equal,
			ObstacleDistance:  0,
		},
	}
	condElseIf := data.AiConditions{
		Condition: data.AiCondition{
			ObstacleType:      data.ObstacleEdge,
			ObstacleDirection: data.Right,
			ObstacleCondition: data.Equal,
			ObstacleDistance:  1,
		},
	}

	s.Run("if", func() {
		sn := data.NewSnake(5, 3, 0, -1, []func(snake *data.Snake){
			func(snake *data.Snake) { snake.DoIf(condIf, 1) },
			func(snake *data.Snake) { snake.Move() },
			func(snake *data.Snake) { snake.DoElseIf(condElseIf, 2) },
			func(snake *data.Snake) { snake.Left() },
			func(snake *data.Snake) { snake.Move() },
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
		sn := data.NewSnake(4, 3, 0, -1, []func(snake *data.Snake){
			func(snake *data.Snake) { snake.DoIf(condIf, 1) },
			func(snake *data.Snake) { snake.Move() },
			func(snake *data.Snake) { snake.DoElseIf(condElseIf, 2) },
			func(snake *data.Snake) { snake.Left() },
			func(snake *data.Snake) { snake.Move() },
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
	g := data.NewGame(gameWidth, gameHeight, &pa)
	s.games.AddGame(g)

	condIf := data.AiConditions{
		Condition: data.AiCondition{
			ObstacleType:      data.ObstacleEdge,
			ObstacleDirection: data.Right,
			ObstacleCondition: data.Equal,
			ObstacleDistance:  0,
		},
	}
	condElseIf := data.AiConditions{
		Condition: data.AiCondition{
			ObstacleType:      data.ObstacleEdge,
			ObstacleDirection: data.Right,
			ObstacleCondition: data.Equal,
			ObstacleDistance:  1,
		},
	}

	s.Run("if", func() {
		sn := data.NewSnake(5, 3, 0, -1, []func(snake *data.Snake){
			func(snake *data.Snake) { snake.DoIf(condIf, 1) },
			func(snake *data.Snake) { snake.Move() },
			func(snake *data.Snake) { snake.DoElseIf(condElseIf, 2) },
			func(snake *data.Snake) { snake.Left() },
			func(snake *data.Snake) { snake.Move() },
			func(snake *data.Snake) { snake.DoElse(2) },
			func(snake *data.Snake) { snake.Right() },
			func(snake *data.Snake) { snake.Move() },
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
		sn := data.NewSnake(4, 3, 0, -1, []func(snake *data.Snake){
			func(snake *data.Snake) { snake.DoIf(condIf, 1) },
			func(snake *data.Snake) { snake.Move() },
			func(snake *data.Snake) { snake.DoElseIf(condElseIf, 2) },
			func(snake *data.Snake) { snake.Left() },
			func(snake *data.Snake) { snake.Move() },
			func(snake *data.Snake) { snake.DoElse(2) },
			func(snake *data.Snake) { snake.Right() },
			func(snake *data.Snake) { snake.Move() },
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
		sn := data.NewSnake(3, 3, 0, -1, []func(snake *data.Snake){
			func(snake *data.Snake) { snake.DoIf(condIf, 1) },
			func(snake *data.Snake) { snake.Move() },
			func(snake *data.Snake) { snake.DoElseIf(condElseIf, 2) },
			func(snake *data.Snake) { snake.Left() },
			func(snake *data.Snake) { snake.Move() },
			func(snake *data.Snake) { snake.DoElse(2) },
			func(snake *data.Snake) { snake.Right() },
			func(snake *data.Snake) { snake.Move() },
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
