package javalang

import (
	"docdog/src/helper"
	"docdog/src/notations"
	"errors"
	"strings"
)

const Private = "private"
const Public = "public"

const arrayIdentifier = "array"
const listIdentifier = "List<"

func JavaVariableHandler(line string, index int, wholeFile []string) (notations.Variable, error) {
	temp := helper.SeparateLineByTags(line)

	tempVar := notations.Variable{
		Name:        "",
		Description: "",
		Typ:         "",
		Notnull:     false,
	}

	if !strings.Contains(line, "(") && !strings.Contains(line, "class") && !strings.Contains(line, "{") {
		tempVar.Name = strings.ReplaceAll(temp[2], ";", "")
		tempVar.Typ = temp[1]
		i := 1
		for i <= 3 {
			if notations.IsNotNullNotation(wholeFile[index-i]) {
				tempVar.Notnull = true
			}
			if notations.IsDescriptionNotation(wholeFile[index-i]) {
				tempVar.Description = helper.GetStringFromQouteLine(wholeFile[i])
			}

			i++
		}
		return tempVar, nil
	}

	return tempVar, errors.New("is function")
}
