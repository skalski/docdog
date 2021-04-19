package spring

import (
	"docdog/src/notations"
	"fmt"
	"strings"
)

const arrayIdentifier = "[]"
const listIdentifier = "List<"
const mapping = "Mapping("
const requestBody = "@RequestBody"
const pathVariable = "@PathVariable"

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
		fmt.Println(wholeFile[i])
		if strings.Contains(wholeFile[i], mapping) {
			tempVar.Url = GetStringFromQouteLine(wholeFile[i])
			tempVar.HttpType = GetProtocolFormMappingTag(wholeFile[i])
		}
		ls := strings.Split(wholeFile[i], " ")
		for i, command := range ls {
			if strings.Contains(command, requestBody) {
				tempVar.Objects = append(tempVar.Objects, ls[i+1])
			}
			if strings.Contains(command, pathVariable) {
				params := notations.Params{
					Name:    strings.Replace(ls[i+2], ")", "", 1),
					VarType: ls[i+1],
				}
				if IsArrayType(ls[i+1]) {
					params.IsArray = true
				}
				tempVar.Params = append(tempVar.Params, params)
			}
		}
		if strings.Contains(wholeFile[i], "{") && !strings.Contains(wholeFile[i], "}") {
			break
		}
		i++
	}
	return tempVar
}

func IsArrayType(line string) bool {
	return strings.Contains(line, arrayIdentifier) || strings.Contains(line, listIdentifier)
}

func GetProtocolFormMappingTag(str string) (result string) {
	s := strings.Index(str, "@")
	if s == -1 {
		return
	}
	s += len("M")
	e := strings.Index(str[s:], "M")
	if e == -1 {
		return
	}
	return str[s : s+e]
}

func GetStringFromQouteLine(str string) (result string) {
	s := strings.Index(str, "\"")
	if s == -1 {
		return
	}
	s += len("\"")
	e := strings.Index(str[s:], "\"")
	if e == -1 {
		return
	}
	return str[s : s+e]
}
