package data

import (
	"context"
	"github.com/andreamper220/snakeai/pkg/logger"
	"time"

	grpcclients "github.com/andreamper220/snakeai/internal/server/infrastructure/grpc"
	pb "github.com/andreamper220/snakeai/proto"
)

type ObstacleType string

const (
	ObstacleEdge  = ObstacleType("edge")
	ObstacleFood  = ObstacleType("food")
	ObstacleSnake = ObstacleType("snake")
)

type ObstacleDirection string

const (
	Forward = ObstacleDirection("forward")
	Left    = ObstacleDirection("left")
	Right   = ObstacleDirection("right")
)

type ObstacleCondition string

const (
	Equal          = ObstacleCondition("==")
	NotEqual       = ObstacleCondition("!=")
	LessThan       = ObstacleCondition("<<")
	GreaterThan    = ObstacleCondition(">>")
	LessOrEqual    = ObstacleCondition("<=")
	GreaterOrEqual = ObstacleCondition(">=")
)

type ConditionOperator string

const (
	And     = ConditionOperator("&&")
	Or      = ConditionOperator("||")
	Default = ConditionOperator("")
)

// AiConditions contains condition groups, condition, operator and is_negative property.
type AiConditions struct {
	Conditions          []AiConditions    // is nil - if no inner condition groups
	Condition           AiCondition       // is nil - if there are inner condition groups
	Operator            ConditionOperator // is empty string - if no inner condition groups
	IsNegativeCondition bool
}

func (conditions AiConditions) Check(snake *Snake, game *Game) bool {
	if len(conditions.Conditions) == 0 {
		// if no inner condition groups
		return conditions.Condition.Check(snake, game) && !conditions.IsNegativeCondition
	} else {
		// if there are inner condition groups
		for _, condition := range conditions.Conditions {
			conditionsCheck := condition.Check(snake, game) && !condition.IsNegativeCondition
			logger.Log.Info(conditionsCheck, conditions.Operator)
			if len(conditions.Conditions) == 1 {
				return conditionsCheck
			}

			if conditions.Operator == And && !conditionsCheck {
				return conditions.IsNegativeCondition
			}
			if conditions.Operator == Or && conditionsCheck {
				return !conditions.IsNegativeCondition
			}
		}

		switch conditions.Operator {
		case And:
			return !conditions.IsNegativeCondition
		case Or:
			return conditions.IsNegativeCondition
		default:
			return false
		}
	}
}

// AiCondition contains obstacle type, direction, condition and distance.
type AiCondition struct {
	ObstacleType      ObstacleType
	ObstacleDirection ObstacleDirection
	ObstacleCondition ObstacleCondition
	ObstacleDistance  int
}

func (condition AiCondition) Check(snake *Snake, game *Game) bool {
	direction := snake.Direction
	var obstaclePoints = make([]Point, 0)
	switch condition.ObstacleType {
	case ObstacleEdge:
		if game.Party.MapId != "" {
			obstaclePoints = append(obstaclePoints, checkObstacleWalls(game, direction, condition)...)
		}
		obstaclePoints = append(obstaclePoints, checkObstacleEdges(game, snake.Body[0], direction, condition)...)
	case ObstacleFood:
		obstaclePoints = append(obstaclePoints, checkObstacleFood(game.Food.Position.X, game.Food.Position.Y, direction, condition)...)
	case ObstacleSnake:
		obstaclePoints = append(obstaclePoints, checkObstacleSnakes(game, snake, direction, condition)...)
	}

	for _, obstaclePoint := range obstaclePoints {
		if check := condition.checkConditionDirection(direction, obstaclePoint, snake.Body[0]); check {
			return true
		}
	}
	return false
}

func (condition AiCondition) checkConditionDirection(direction, obstaclePoint, head Point) bool {
	switch condition.ObstacleDirection {
	case Forward:
		switch condition.ObstacleCondition {
		case Equal:
			if (direction.Y == 0 && abs(obstaclePoint.X-head.X) == condition.ObstacleDistance && obstaclePoint.Y == head.Y) ||
				(direction.X == 0 && abs(obstaclePoint.Y-head.Y) == condition.ObstacleDistance && obstaclePoint.X == head.X) {
				return true
			}
		case NotEqual:
			if (direction.Y == 0 && abs(obstaclePoint.X-head.X) != condition.ObstacleDistance && obstaclePoint.Y == head.Y) ||
				(direction.X == 0 && abs(obstaclePoint.Y-head.Y) != condition.ObstacleDistance && obstaclePoint.X == head.X) {
				return true
			}
		case LessThan:
			if (direction.Y == 0 && abs(obstaclePoint.X-head.X) < condition.ObstacleDistance && obstaclePoint.Y == head.Y) ||
				(direction.X == 0 && abs(obstaclePoint.Y-head.Y) < condition.ObstacleDistance && obstaclePoint.X == head.X) {
				return true
			}
		case GreaterThan:
			if (direction.Y == 0 && abs(obstaclePoint.X-head.X) > condition.ObstacleDistance && obstaclePoint.Y == head.Y) ||
				(direction.X == 0 && abs(obstaclePoint.Y-head.Y) > condition.ObstacleDistance && obstaclePoint.X == head.X) {
				return true
			}
		case LessOrEqual:
			if (direction.Y == 0 && abs(obstaclePoint.X-head.X) <= condition.ObstacleDistance && obstaclePoint.Y == head.Y) ||
				(direction.X == 0 && abs(obstaclePoint.Y-head.Y) <= condition.ObstacleDistance && obstaclePoint.X == head.X) {
				return true
			}
		case GreaterOrEqual:
			if (direction.Y == 0 && abs(obstaclePoint.X-head.X) >= condition.ObstacleDistance && obstaclePoint.Y == head.Y) ||
				(direction.X == 0 && abs(obstaclePoint.Y-head.Y) >= condition.ObstacleDistance && obstaclePoint.X == head.X) {
				return true
			}
		}
	case Right, Left:
		switch condition.ObstacleCondition {
		case Equal:
			if (direction.Y == 0 && obstaclePoint.X == head.X && abs(obstaclePoint.Y-head.Y) == condition.ObstacleDistance) ||
				(direction.X == 0 && obstaclePoint.Y == head.Y && abs(obstaclePoint.X-head.X) == condition.ObstacleDistance) {
				return true
			}
		case NotEqual:
			if (direction.Y == 0 && obstaclePoint.X == head.X && abs(obstaclePoint.Y-head.Y) != condition.ObstacleDistance) ||
				(direction.X == 0 && obstaclePoint.Y == head.Y && abs(obstaclePoint.X-head.X) != condition.ObstacleDistance) {
				return true
			}
		case LessThan:
			if (direction.Y == 0 && obstaclePoint.X == head.X && abs(obstaclePoint.Y-head.Y) < condition.ObstacleDistance) ||
				(direction.X == 0 && obstaclePoint.Y == head.Y && abs(obstaclePoint.X-head.X) < condition.ObstacleDistance) {
				return true
			}
		case GreaterThan:
			if (direction.Y == 0 && obstaclePoint.X == head.X && abs(obstaclePoint.Y-head.Y) > condition.ObstacleDistance) ||
				(direction.X == 0 && obstaclePoint.Y == head.Y && abs(obstaclePoint.X-head.X) > condition.ObstacleDistance) {
				return true
			}
		case LessOrEqual:
			if (direction.Y == 0 && obstaclePoint.X == head.X && abs(obstaclePoint.Y-head.Y) <= condition.ObstacleDistance) ||
				(direction.X == 0 && obstaclePoint.Y == head.Y && abs(obstaclePoint.X-head.X) <= condition.ObstacleDistance) {
				return true
			}
		case GreaterOrEqual:
			if (direction.Y == 0 && obstaclePoint.X == head.X && abs(obstaclePoint.Y-head.Y) >= condition.ObstacleDistance) ||
				(direction.X == 0 && obstaclePoint.Y == head.Y && abs(obstaclePoint.X-head.X) >= condition.ObstacleDistance) {
				return true
			}
		}
	}
	return false
}

func checkObstacleWalls(game *Game, direction Point, condition AiCondition) []Point {
	var obstaclePoints = make([]Point, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	gameMap, err := grpcclients.EditorClient.GetMap(ctx, &pb.GetMapRequest{
		Id: game.Party.MapId,
	})
	if err == nil {
		requestObstacles := gameMap.GetMap().GetStruct().GetObstacles()
		for i := 0; i < len(requestObstacles); i++ {
			x := requestObstacles[i].GetCx() + 1
			y := requestObstacles[i].GetCy() + 1
			if direction.X == 0 {
				if condition.ObstacleDirection == Left {
					if direction.Y == 1 {
						x--
					} else {
						x++
					}
				} else if condition.ObstacleDirection == Right {
					if direction.Y == 1 {
						x++
					} else {
						x--
					}
				} else {
					y -= int32(direction.Y)
				}
			} else {
				if condition.ObstacleDirection == Left {
					if direction.X == 1 {
						y++
					} else {
						y--
					}
				} else if condition.ObstacleDirection == Right {
					if direction.X == 1 {
						y--
					} else {
						y++
					}
				} else {
					x -= int32(direction.X)
				}
			}
			obstaclePoints = append(obstaclePoints, Point{X: int(x), Y: int(y)})
		}
	}

	return obstaclePoints
}

func checkObstacleEdges(game *Game, head Point, direction Point, condition AiCondition) []Point {
	var obstaclePoints = make([]Point, 0)

	switch condition.ObstacleDirection {
	case Forward:
		if direction.X == 0 {
			if direction.Y == 1 {
				obstaclePoints = append(obstaclePoints, Point{X: head.X, Y: game.Height})
			} else {
				obstaclePoints = append(obstaclePoints, Point{X: head.X, Y: 1})
			}
		} else {
			if direction.X == 1 {
				obstaclePoints = append(obstaclePoints, Point{X: game.Width, Y: head.Y})
			} else {
				obstaclePoints = append(obstaclePoints, Point{X: 1, Y: head.Y})
			}
		}
	case Right:
		if direction.X == 0 {
			if direction.Y == 1 {
				obstaclePoints = append(obstaclePoints, Point{X: 1, Y: head.Y})
			} else {
				obstaclePoints = append(obstaclePoints, Point{X: game.Width, Y: head.Y})
			}
		} else {
			if direction.X == 1 {
				obstaclePoints = append(obstaclePoints, Point{X: head.X, Y: game.Height})
			} else {
				obstaclePoints = append(obstaclePoints, Point{X: head.X, Y: 1})
			}
		}
	case Left:
		if direction.X == 0 {
			if direction.Y == 1 {
				obstaclePoints = append(obstaclePoints, Point{X: game.Width, Y: head.Y})
			} else {
				obstaclePoints = append(obstaclePoints, Point{X: 1, Y: head.Y})
			}
		} else {
			if direction.X == 1 {
				obstaclePoints = append(obstaclePoints, Point{X: head.X, Y: 1})
			} else {
				obstaclePoints = append(obstaclePoints, Point{X: head.X, Y: game.Height})
			}
		}
	}

	return obstaclePoints
}

func checkObstacleFood(x, y int, direction Point, condition AiCondition) []Point {
	var obstaclePoints = make([]Point, 0)

	if direction.X == 0 {
		if condition.ObstacleDirection == Left {
			if direction.Y == 1 {
				x--
			} else {
				x++
			}
		} else if condition.ObstacleDirection == Right {
			if direction.Y == 1 {
				x++
			} else {
				x--
			}
		} else {
			y -= direction.Y
		}
	} else {
		if condition.ObstacleDirection == Left {
			if direction.X == 1 {
				y++
			} else {
				y--
			}
		} else if condition.ObstacleDirection == Right {
			if direction.X == 1 {
				y--
			} else {
				y++
			}
		} else {
			x -= direction.X
		}
	}
	obstaclePoints = append(obstaclePoints, Point{X: x, Y: y})

	return obstaclePoints
}

func checkObstacleSnakes(game *Game, snake *Snake, direction Point, condition AiCondition) []Point {
	var obstaclePoints = make([]Point, 0)

	for _, sn := range game.GetSnakes() {
		sn.RLock()
		if sn != snake {
			for _, bodyPoint := range sn.Body {
				x := bodyPoint.X
				y := bodyPoint.Y
				if direction.X == 0 {
					if condition.ObstacleDirection == Left {
						if direction.Y == 1 {
							x--
						} else {
							x++
						}
					} else if condition.ObstacleDirection == Right {
						if direction.Y == 1 {
							x++
						} else {
							x--
						}
					} else {
						y -= direction.Y
					}
				} else {
					if condition.ObstacleDirection == Left {
						if direction.X == 1 {
							y++
						} else {
							y--
						}
					} else if condition.ObstacleDirection == Right {
						if direction.X == 1 {
							y--
						} else {
							y++
						}
					} else {
						x -= direction.X
					}
				}
				obstaclePoints = append(obstaclePoints, Point{X: x, Y: y})
			}
		}
		sn.RUnlock()
	}

	return obstaclePoints
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
