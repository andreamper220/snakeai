package data

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
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

// AiCondition contains obstacle type, direction, condition and distance.
type AiCondition struct {
	ObstacleType        ObstacleType
	ObstacleDirection   ObstacleDirection
	ObstacleCondition   ObstacleCondition
	ObstacleDistance    int
	IsNegativeCondition bool
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
			if condition.IsNegativeCondition {
				return false
			}
			return true
		}
	}
	if condition.IsNegativeCondition {
		return true
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

func GenerateAiFunctions(ai string) ([]func(snake *Snake), error) {
	if strings.Count(ai, "(") != strings.Count(ai, ")") {
		return nil, errors.New("parenthesis count does not match")
	}
	if strings.Count(ai, "{") != strings.Count(ai, "}") {
		return nil, errors.New("curly brackets count does not match")
	}

	return processAi(ai), nil
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

func processAi(ai string) []func(snake *Snake) {
	aiFunctions := make([]func(snake *Snake), 0)
	if strings.Index(ai, "if") == 0 {
		aiFunctions = append(aiFunctions, processConditions(ai)...)
	} else {
		aiFunctions = append(aiFunctions, processActions(ai)...)
	}

	return aiFunctions
}

func processActions(ai string) []func(snake *Snake) {
	aiStrings := strings.Split(ai, `;`)
	aiFunctions := make([]func(snake *Snake), 0)
	for i, aiString := range aiStrings {
		if aiString != "" {
			switch aiString {
			case "right":
				aiFunctions = append(aiFunctions, func(snake *Snake) { snake.Right() })
			case "left":
				aiFunctions = append(aiFunctions, func(snake *Snake) { snake.Left() })
			case "move":
				aiFunctions = append(aiFunctions, func(snake *Snake) { snake.Move() })
			default:
				aiStr := aiString + ";" + strings.Join(aiStrings[i+1:], `;`)
				aiFunctions = append(aiFunctions, processAi(aiStr)...)
				return aiFunctions
			}
		}
	}

	return aiFunctions
}

func processConditions(ai string) []func(snake *Snake) {
	aiFunctions := make([]func(snake *Snake), 0)
	// process 'if'
	ifCondition, ifActions, aiNotProcessedString := processConditionString(ai)
	if len(ifActions) > 0 {
		aiFunctionsIf := []func(snake *Snake){
			func(snake *Snake) { snake.DoIf(ifCondition, len(ifActions)) },
		}
		aiFunctionsIf = append(aiFunctionsIf, ifActions...)
		// process 'elseif'
		var aiFunctionsElseIf []func(snake *Snake)
		for i := 0; i <= strings.Count(aiNotProcessedString, "elseif"); i++ {
			if strings.Index(aiNotProcessedString, "elseif") == 0 {
				elseIfCondition, elseIfActions, notProcessedString := processConditionString(aiNotProcessedString)
				aiNotProcessedString = notProcessedString
				if len(elseIfActions) > 0 {
					aiFunctionsElseIf = append(aiFunctionsElseIf, func(snake *Snake) {
						snake.DoElseIf(elseIfCondition, len(elseIfActions))
					})
					aiFunctionsElseIf = append(aiFunctionsElseIf, elseIfActions...)
					aiFunctionsIf = append(aiFunctionsIf, aiFunctionsElseIf...)
				}
			}
		}
		// process 'else'
		if strings.Index(aiNotProcessedString, "else") == 0 {
			elseActions, _ := processConditionActionsString(aiNotProcessedString)
			if len(elseActions) > 0 {
				aiFunctionsElse := []func(snake *Snake){
					func(snake *Snake) { snake.DoElse(len(elseActions)) },
				}
				aiFunctionsElse = append(aiFunctionsElse, elseActions...)
				aiFunctionsIf = append(aiFunctionsIf, aiFunctionsElse...)
			}
		}
		aiFunctions = append(aiFunctions, aiFunctionsIf...)
	}
	return aiFunctions
}

func processConditionString(ai string) (AiCondition, []func(snake *Snake), string) {
	aiNotProcessedString := ""
	actions := make([]func(snake *Snake), 0)
	condition := AiCondition{}
	conditionString, index := getValueBetweenSymbols("(", ")", ai)
	if conditionString != "" {
		isNegativeCondition := false
		if strings.Index(conditionString, "!") == 0 {
			isNegativeCondition = true
			conditionString = conditionString[2:]
		}
		conditionStrings := strings.Split(conditionString, `_`)
		numberRegExp := regexp.MustCompile("[0-9]+")
		numbers := numberRegExp.FindAllString(conditionStrings[1], 1)
		if len(numbers) > 0 {
			number := numbers[0]
			numberIndex := strings.Index(conditionStrings[1], number)
			conditionSeparator := conditionStrings[1][numberIndex-2 : numberIndex]
			conditionStringsInner := strings.Split(conditionStrings[1], conditionSeparator)
			obstacleDistance, _ := strconv.Atoi(conditionStringsInner[1])

			condition = AiCondition{
				ObstacleType:        ObstacleType(conditionStrings[0]),
				ObstacleDirection:   ObstacleDirection(conditionStringsInner[0]),
				ObstacleCondition:   ObstacleCondition(conditionSeparator),
				ObstacleDistance:    obstacleDistance,
				IsNegativeCondition: isNegativeCondition,
			}
		}
		actions, aiNotProcessedString = processConditionActionsString(ai)

		return condition, actions, aiNotProcessedString
	}

	return condition, actions, ai[index+1:]
}

func processConditionActionsString(ai string) ([]func(snake *Snake), string) {
	actions := make([]func(snake *Snake), 0)
	actionsString, index := getValueBetweenSymbols("{", "}", ai)
	if actionsString != "" {
		actions = processAi(actionsString)
	}
	return actions, ai[index+1:]
}

func getValueBetweenSymbols(first, second, haystack string) (string, int) {
	indexBegin := strings.Index(haystack, first)
	if indexBegin >= 0 {
		indexEnd := strings.Index(haystack, second)
		if indexEnd >= 0 {
			return haystack[indexBegin+1 : indexEnd], indexEnd
		}
	}
	return "", indexBegin
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
