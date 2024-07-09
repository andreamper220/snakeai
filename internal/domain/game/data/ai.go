package data

import (
	"github.com/andreamper220/snakeai/pkg/logger"
	"regexp"
	"strconv"
	"strings"
)

type obstacleType string

const (
	ObstacleEdge  = obstacleType("edge")
	ObstacleFood  = obstacleType("food")
	ObstacleSnake = obstacleType("snake")
)

type obstacleDirection string

const (
	Forward = obstacleDirection("forward")
	Left    = obstacleDirection("left")
	Right   = obstacleDirection("right")
)

type obstacleCondition string

const (
	Equal          = obstacleCondition("==")
	NotEqual       = obstacleCondition("!=")
	LessThan       = obstacleCondition("<<")
	GreaterThan    = obstacleCondition(">>")
	LessOrEqual    = obstacleCondition("<=")
	GreaterOrEqual = obstacleCondition(">=")
)

type AiCondition struct {
	obstacleType      obstacleType
	obstacleDirection obstacleDirection
	obstacleCondition obstacleCondition
	obstacleDistance  int
}

func (condition AiCondition) Check(snake *Snake, game *Game) bool {
	head := snake.Body[0]
	direction := snake.Direction
	var obstaclePoint = Point{}
	switch condition.obstacleType {
	case ObstacleEdge:
		switch condition.obstacleDirection {
		case Forward:
			if direction.X == 0 {
				if direction.Y == 1 {
					obstaclePoint = Point{X: head.X, Y: game.Height}
				} else {
					obstaclePoint = Point{X: head.X, Y: 1}
				}
			} else {
				if direction.X == 1 {
					obstaclePoint = Point{X: game.Width, Y: head.Y}
				} else {
					obstaclePoint = Point{X: 1, Y: head.Y}
				}
			}
		case Right:
			if direction.X == 0 {
				if direction.Y == 1 {
					obstaclePoint = Point{X: 1, Y: head.Y}
				} else {
					obstaclePoint = Point{X: game.Width, Y: head.Y}
				}
			} else {
				if direction.X == 1 {
					obstaclePoint = Point{X: head.X, Y: game.Height}
				} else {
					obstaclePoint = Point{X: head.X, Y: 1}
				}
			}
		case Left:
			if direction.X == 0 {
				if direction.Y == 1 {
					obstaclePoint = Point{X: game.Width, Y: head.Y}
				} else {
					obstaclePoint = Point{X: 1, Y: head.Y}
				}
			} else {
				if direction.X == 1 {
					obstaclePoint = Point{X: head.X, Y: 1}
				} else {
					obstaclePoint = Point{X: head.X, Y: game.Height}
				}
			}
		}
	case ObstacleFood:
		x := game.Food.Position.X
		y := game.Food.Position.Y
		if direction.X == 0 {
			if condition.obstacleDirection == Left {
				x++
			} else if condition.obstacleDirection == Right {
				x--
			}
		} else {
			if condition.obstacleDirection == Left {
				y--
			} else if condition.obstacleDirection == Right {
				y++
			}
		}
		obstaclePoint = Point{X: x, Y: y}
	}

	// TODO add snake condition

	switch condition.obstacleDirection {
	case Forward:
		logger.Log.Infof("%v %v", obstaclePoint, head)
		switch condition.obstacleCondition {
		case Equal:
			if (direction.Y == 0 && abs(obstaclePoint.X-head.X) == condition.obstacleDistance) ||
				(direction.X == 0 && abs(obstaclePoint.Y-head.Y) == condition.obstacleDistance) {
				return true
			}
		case NotEqual:
			if (direction.Y == 0 && abs(obstaclePoint.X-head.X) != condition.obstacleDistance) ||
				(direction.X == 0 && abs(obstaclePoint.Y-head.Y) != condition.obstacleDistance) {
				return true
			}
		case LessThan:
			if (direction.Y == 0 && abs(obstaclePoint.X-head.X) < condition.obstacleDistance) ||
				(direction.X == 0 && abs(obstaclePoint.Y-head.Y) < condition.obstacleDistance) {
				return true
			}
		case GreaterThan:
			if (direction.Y == 0 && abs(obstaclePoint.X-head.X) > condition.obstacleDistance) ||
				(direction.X == 0 && abs(obstaclePoint.Y-head.Y) > condition.obstacleDistance) {
				return true
			}
		case LessOrEqual:
			if (direction.Y == 0 && abs(obstaclePoint.X-head.X) <= condition.obstacleDistance) ||
				(direction.X == 0 && abs(obstaclePoint.Y-head.Y) <= condition.obstacleDistance) {
				return true
			}
		case GreaterOrEqual:
			if (direction.Y == 0 && abs(obstaclePoint.X-head.X) >= condition.obstacleDistance) ||
				(direction.X == 0 && abs(obstaclePoint.Y-head.Y) >= condition.obstacleDistance) {
				return true
			}
		}
	case Right, Left:
		logger.Log.Infof("%v %v", obstaclePoint, head)
		switch condition.obstacleCondition {
		case Equal:
			if (direction.Y == 0 && obstaclePoint.X == head.X && abs(obstaclePoint.Y-head.Y) == condition.obstacleDistance) ||
				(direction.X == 0 && obstaclePoint.Y == head.Y && abs(obstaclePoint.X-head.X) == condition.obstacleDistance) {
				return true
			}
		case NotEqual:
			if (direction.Y == 0 && obstaclePoint.X == head.X && abs(obstaclePoint.Y-head.Y) != condition.obstacleDistance) ||
				(direction.X == 0 && obstaclePoint.Y == head.Y && abs(obstaclePoint.X-head.X) != condition.obstacleDistance) {
				return true
			}
		case LessThan:
			if (direction.Y == 0 && obstaclePoint.X == head.X && abs(obstaclePoint.Y-head.Y) < condition.obstacleDistance) ||
				(direction.X == 0 && obstaclePoint.Y == head.Y && abs(obstaclePoint.X-head.X) < condition.obstacleDistance) {
				return true
			}
		case GreaterThan:
			if (direction.Y == 0 && obstaclePoint.X == head.X && abs(obstaclePoint.Y-head.Y) > condition.obstacleDistance) ||
				(direction.X == 0 && obstaclePoint.Y == head.Y && abs(obstaclePoint.X-head.X) > condition.obstacleDistance) {
				return true
			}
		case LessOrEqual:
			if (direction.Y == 0 && obstaclePoint.X == head.X && abs(obstaclePoint.Y-head.Y) <= condition.obstacleDistance) ||
				(direction.X == 0 && obstaclePoint.Y == head.Y && abs(obstaclePoint.X-head.X) <= condition.obstacleDistance) {
				return true
			}
		case GreaterOrEqual:
			if (direction.Y == 0 && obstaclePoint.X == head.X && abs(obstaclePoint.Y-head.Y) >= condition.obstacleDistance) ||
				(direction.X == 0 && obstaclePoint.Y == head.Y && abs(obstaclePoint.X-head.X) >= condition.obstacleDistance) {
				return true
			}
		}
	}
	return false
}

func GenerateAiFunctions(ai string) ([]func(snake *Snake), error) {
	aiFunctions := processAi(ai)

	return aiFunctions, nil
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
	actionsString := ""
	conditionString, index := getValueBetweenSymbols("(", ")", ai)
	if conditionString != "" {
		actionsString, index = getValueBetweenSymbols("{", "}", ai)
		if actionsString != "" {
			conditionStrings := strings.Split(conditionString, `_`)
			numberRegExp := regexp.MustCompile("[0-9]+")
			numbers := numberRegExp.FindAllString(conditionStrings[1], 1)
			if len(numbers) > 0 {
				number := numbers[0]
				numberIndex := strings.Index(conditionStrings[1], number)
				conditionSeparator := conditionStrings[1][numberIndex-2 : numberIndex]
				conditionStringsInner := strings.Split(conditionStrings[1], conditionSeparator)
				obstacleDistance, _ := strconv.Atoi(conditionStringsInner[1])

				condition := AiCondition{
					obstacleType:      obstacleType(conditionStrings[0]),
					obstacleDirection: obstacleDirection(conditionStringsInner[0]),
					obstacleCondition: obstacleCondition(conditionSeparator),
					obstacleDistance:  obstacleDistance,
				}
				actions := processAi(actionsString)

				aiFunctions = []func(snake *Snake){
					func(snake *Snake) { snake.DoIf(condition, len(actions)) },
				}
				aiFunctions = append(aiFunctions, actions...)
			}
		}
	}
	return append(aiFunctions, processAi(ai[index+1:])...)
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
