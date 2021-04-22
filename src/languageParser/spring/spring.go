package spring

import (
	"docdog/src/notations"
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

var methods = [4]string{"post", "get", "delete", "put"}

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
			return m
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
