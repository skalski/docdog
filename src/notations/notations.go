package notations

import "strings"

var endPointIdentifier = "@DD:ENDPOINT"
var paramIdentifier = "@DD:PARAM"
var payloadIdentifier = "@DD:PAYLOAD"
var descriptionIdentifier = "@DD:DESCRIPTION"
var notNullIdentifier = "@DD:NOTNULL"
var typeIdentifier = "@DD:TYPE"
var ignoreIdentifier = "@DD:IGNORE"

var springEndpointController = "@RestController"

type Objects struct {
	Name     string
	Variable []Variable
}

type Variable struct {
	Name        string
	Description string
	Typ         string
	Notnull     bool
	IsArray     bool
}

type Params struct {
	Name        string
	Description string
	Notnull     bool
}

type Endpoint struct {
	Url         string
	Description string
	HttpType    string
	Params      []Params
	Variable    []Variable
	Objects     []string
}

type TempEndpoint struct {
	Url         string
	Description string
	HttpType    string
	Params      []Params
	Objects     []string
}

func IsController(line []byte) bool {
	return IsEndpointNotation(string(line[:]))
}

func IsIgnoreNotation(line string) bool {
	return strings.Contains(line, ignoreIdentifier)
}

func HasIgnoreNotation(index int, wholeFile []string) bool {
	i := 1
	for i <= 3 {
		if IsIgnoreNotation(wholeFile[index-i]) {
			return true
		}
		i++
	}
	return false
}

func IsEndpointNotation(line string) bool {
	return strings.Contains(line, endPointIdentifier) || strings.Contains(line, springEndpointController)
}

func IsParamNotation(line string) bool {
	return strings.Contains(line, paramIdentifier)
}

func IsPayloadNotation(line string) bool {
	return strings.Contains(line, payloadIdentifier)
}

func IsNotNullNotation(line string) bool {
	return strings.Contains(line, notNullIdentifier)
}

func IsDescriptionNotation(line string) bool {
	return strings.Contains(line, descriptionIdentifier)
}

func IsConnectionMethodNotation(line string) bool {
	return strings.Contains(line, typeIdentifier)
}

func CommentEndTag(line string) bool {
	return strings.Contains(line, "*/")
}
