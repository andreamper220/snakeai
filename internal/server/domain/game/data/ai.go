package data

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func GenerateAiFunctions(ai string) ([]func(snake *Snake), error) {
	if strings.Count(ai, "(") != strings.Count(ai, ")") {
		return nil, errors.New("parenthesis count does not match")
	}
	if strings.Count(ai, "{") != strings.Count(ai, "}") {
		return nil, errors.New("curly brackets count does not match")
	}

	return processAi(ai), nil
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

func processConditionString(ai string) (AiConditions, []func(snake *Snake), string) {
	aiNotProcessedString := ""
	actions := make([]func(snake *Snake), 0)
	condition := AiConditions{}

	conditionString, index := getValueBetweenSymbols("if", "then", ai)
	if conditionString != "" {
		conditionString, _ = getConditionStringFromString(conditionString)
		condition, _, index = processConditionsString(conditionString)
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

func processConditionsString(conditionsString string) (AiConditions, ConditionOperator, int) {
	index := 0
	operator := Default
	conditionsStringInner := ""
	aiConditions := AiConditions{}
	condition := AiConditions{}

	if strings.Index(conditionsString, "!") == 0 {
		aiConditions.IsNegativeCondition = true
		conditionsStringInner, _ = getConditionStringFromString(conditionsString)
		if conditionsString[4] == '!' || conditionsString[4] == '(' || conditionsString[1] == '(' {
			condition, operator, index = processConditionsString(conditionsStringInner)
			aiConditions.Conditions = append(aiConditions.Conditions, condition)
		} else {
			conditionsStringInner, index = getConditionStringFromString(conditionsString)
		}
	} else {
		conditionsStringInner, index = getConditionStringFromString(conditionsString)
	}

	if strings.Index(conditionsStringInner, "&&") == index+1 {
		aiConditions.Operator = And
		index += 2
	} else if strings.Index(conditionsStringInner, "||") == index+1 {
		aiConditions.Operator = Or
		index += 2
	}

	// if with down-up flow
	if len(aiConditions.Conditions) > 0 {
		conditionsStringInner = conditionsStringInner[index+1:]
	}

	for {
		aiCondition := AiConditions{}

		if strings.Index(conditionsStringInner, "!") == 0 {
			aiCondition.IsNegativeCondition = true
		}

		conditionString, indexInner := getValueBetweenSymbols("(", ")", conditionsStringInner)
		if conditionString == "" {
			break
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

			aiCondition.Condition = AiCondition{
				ObstacleType:      ObstacleType(conditionStrings[0]),
				ObstacleDirection: ObstacleDirection(conditionStringsInner[0]),
				ObstacleCondition: ObstacleCondition(conditionSeparator),
				ObstacleDistance:  obstacleDistance,
			}
		}

		conditionsStringInner = conditionsStringInner[indexInner+1:]
		if strings.Index(conditionsStringInner, "&&") == 0 {
			aiConditions.Operator = And
			conditionsStringInner = conditionsStringInner[2:]
		} else if strings.Index(conditionsStringInner, "||") == 0 {
			aiConditions.Operator = Or
			conditionsStringInner = conditionsStringInner[2:]
		}

		aiConditions.Conditions = append(aiConditions.Conditions, aiCondition)
	}

	return aiConditions, operator, index
}

func getValueBetweenSymbols(first, second, haystack string) (string, int) {
	indexBegin := strings.Index(haystack, first)
	if indexBegin >= 0 {
		indexEnd := strings.Index(haystack, second)
		if indexEnd >= 0 {
			return haystack[indexBegin+len(first) : indexEnd], indexEnd
		}
	}
	return "", indexBegin
}

func getConditionStringFromString(conditionString string) (string, int) {
	indexBegin := strings.Index(conditionString, "(")
	indexEnd := len(conditionString) - 1
	var rightCount, leftCount int

	if len(conditionString) > 0 {
		for index := len(conditionString) - 1; index >= indexBegin; index-- {
			if conditionString[index] == ')' {
				indexEnd = index + 1
				if rightCount != 0 && leftCount == 0 {
					return conditionString[indexBegin+1 : indexEnd], indexEnd
				}
				rightCount++
			} else if conditionString[index] == '(' {
				leftCount++
			}
			if rightCount == leftCount {
				rightCount = 0
				leftCount = 0
			}
		}
	}
	return conditionString[indexBegin:indexEnd], indexEnd
}
