package spring

import (
	"docdog/src/helper"
	"docdog/src/notations"
	"errors"
	"regexp"
	"strings"
)

const arrayIdentifier = "[]"
const listIdentifier = "List<"
const mapping = "Mapping("
const requestMapping = "@RequestMapping("
const requestBody = "@RequestBody"
const pathVariable = "@PathVariable"
const finalTag = "final"

var notnull = []string{"@NotNull", "@NotBlank", "@NotEmpty"}
var ignore = []string{"@JsonIgnore", "@Ignore"}
var methods = [4]string{"post", "get", "delete", "put"}

func SBVariableHandler(line string, index int, wholeFile []string) (notations.Variable, error) {
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
		for i <= 2 {
			_, foundIgnore := helper.Find(ignore, wholeFile[index-i])
			_, foundNotNull := helper.Find(notnull, wholeFile[index-i])
			if foundNotNull && wholeFile[index-i] != "" {
				tempVar.Notnull = true
			}
			if foundIgnore && wholeFile[index-i] != "" {
				return tempVar, errors.New("is ignored")
			}
			i++
		}
		return tempVar, nil
	}
	return tempVar, errors.New("is function")
}

func CreateArrayType(line string) string {
	return strings.Replace(strings.Replace(strings.Replace(line, ">", "", 1), listIdentifier, "", 1), arrayIdentifier, "", 1)
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
		if strings.Contains(wholeFile[i], mapping) && !strings.Contains(wholeFile[i], requestMapping) {
			tempVar.Url = GetStringFromQouteLine(wholeFile[i])
			tempVar.HttpType = GetProtocolFormMappingTag(wholeFile[i])
		}
		if strings.Contains(wholeFile[i], requestMapping) {
			tempVar.HttpType = FetchMethodFromMapping(wholeFile[i])
			r, _ := regexp.Compile("value\\s*=\\s*\"(.+)\"")
			if r.MatchString(wholeFile[i]) {
				res := r.FindAllStringSubmatch(wholeFile[i], -1)
				for i := range res {
					tempVar.Url = GetStringFromQouteLine(strings.Join(res[i], ""))
				}

			}
		}

		ls := strings.Split(wholeFile[i], " ")
		for i, command := range ls {
			if strings.Contains(command, requestBody) {
				pos := IsFinal(ls[i+1], i)
				tempVar.Objects = append(tempVar.Objects, ls[pos])
			}
			if strings.Contains(command, pathVariable) {
				pos := IsFinal(ls[i+1], i)
				params := notations.Params{
					Name:    strings.Replace(ls[i+2], ")", "", 1),
					VarType: ls[pos],
				}
				if IsArrayType(ls[pos]) {
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

func IsFinal(s string, i int) int {
	if strings.Contains(s, finalTag) {
		return i + 2
	}
	return i + 1
}

func FetchMethodFromMapping(s string) string {
	for _, m := range methods {
		if strings.Contains(strings.ToLower(s), strings.ToLower(m)) {
			return strings.ToLower(m)
		}
	}
	return ""
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
	return strings.ToLower(str[s : s+e])
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
