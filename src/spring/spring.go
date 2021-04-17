package spring

import (
	"docdog/src/helper"
	"docdog/src/notations"
	"errors"
	"strings"
)

const Private = "private"
const Public = "public"

const arrayIdentifier = "[]"
const listIdentifier = "List<"
const mapping = "Mapping("
const requestBody = "@RequestBody"
const pathVariable = "@PathVariable"

func JavaVariableHandler(line string, index int, wholeFile []string) (notations.Variable, error) {
	temp := helper.SeparateLineByTags(line)

	tempVar := notations.Variable{
		Name:        "",
		Description: "",
		Typ:         "",
		Notnull:     false,
		IsArray:     false,
	}

	if !strings.Contains(line, "(") && !strings.Contains(line, "class") && !strings.Contains(line, "{") {
		if len(temp) <= 2 {
			return tempVar, errors.New("is malformed function or variable")
		}
		tempVar.Name = strings.ReplaceAll(temp[2], ";", "")
		if IsArrayType(line) {
			tempVar.IsArray = true
			tempVar.Typ = CreateArrayType(line)
		} else {
			tempVar.Typ = temp[1]
		}
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

func CreateApiEndpoint(index int, wholeFile []string) notations.TempEndpoint {
	tempVar := notations.TempEndpoint{
		Url:         "",
		Description: "",
		HttpType:    "",
		Params:      nil,
		Objects:     nil,
	}
	i := index + 1

	for {
		if notations.CommentEndTag(wholeFile[i]) {
			break
		}
		if notations.IsEndpointNotation(wholeFile[i]) {
			tempVar.Url = helper.GetStringFromQouteLine(wholeFile[i])
		}
		if notations.IsDescriptionNotation(wholeFile[i]) {
			tempVar.Description = helper.GetStringFromQouteLine(wholeFile[i])
		}
		if notations.IsConnectionMethodNotation(wholeFile[i]) {
			tempPayload := helper.SeparateLineByTags(wholeFile[i])
			tempVar.HttpType = tempPayload[1]
		}
		if notations.IsPayloadNotation(wholeFile[i]) {
			tempPayload := helper.SeparateLineByTags(wholeFile[i])
			tempVar.Objects = append(tempVar.Objects, tempPayload[1])
		}
		i++
	}
	return tempVar
}

func IsArrayType(line string) bool {
	return strings.Contains(line, arrayIdentifier) || strings.Contains(line, listIdentifier)
}

func CreateArrayType(line string) string {
	return strings.Replace(strings.Replace(strings.Replace(line, ">", "", 1), listIdentifier, "", 1), arrayIdentifier, "", 1)
}
