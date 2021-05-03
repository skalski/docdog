package javaLang

import (
	"docdog/src/helper"
	"docdog/src/notations"
	"errors"
	"strings"
)

const Private = "private"
const Public = "public"
const Abstract = "abstract"
const Implements = "implements"
const Interface = "interface"

const arrayIdentifier = "[]"
const listIdentifier = "List<"

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
			tempVar.Typ = CreateArrayType(temp[1])
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

func IsArrayType(line string) bool {
	return strings.Contains(line, arrayIdentifier) || strings.Contains(line, listIdentifier)
}

func CreateArrayType(line string) string {
	return strings.Replace(strings.Replace(strings.Replace(line, ">", "", 1), listIdentifier, "", 1), arrayIdentifier, "", 1)
}

func ChckImpl(fls []string) string {
	for _, s := range fls {
		if strings.Contains(s, Implements) {
			temp := helper.SeparateLineByTags(s)
			for i, s := range temp {
				if strings.Contains(s, Implements) {
					return temp[i+1]
				}
			}
		}
	}
	return ""
}

func IsItrf(fls []string) bool {
	for _, s := range fls {
		if strings.Contains(s, Interface) {
			return true
		}
	}
	return false
}

func IsAbstrc(fls []string) bool {
	for _, s := range fls {
		if strings.Contains(s, Abstract) {
			return true
		}
	}
	return false
}
