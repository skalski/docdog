package notations

import (
	"strings"
)

var endPointIdentifier = "@DD:ENDPOINT"
var paramIdentifier = "@DD:PARAM"
var payloadIdentifier = "@DD:PAYLOAD"
var descriptionIdentifier = "@DD:DESCRIPTION"
var notNullIdentifier = "@DD:NOTNULL"
var typeIdentifier = "@DD:TYPE"
var ignoreIdentifier = "@DD:IGNORE"
var responseIdentifier = "@DD:RESPONSE"

var springEndpointController = "@RestController"
var springMappingTagIdentifier = "Mapping("

type Objects struct {
	Name        string
	PackageName string
	Variable    []Variable
	Implements  []string
	Imports     []string
}

type Abstract struct {
	PackageName string
	Name        string
	Variable    []Variable
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
	VarType     string
	Notnull     bool
	IsArray     bool
}

type Response struct {
	HttpCode   string
	Type       string
	ObjectType string
}

type Endpoint struct {
	Url         string
	Description string
	HttpType    string
	Params      []Params
	Variable    []Variable
	Response    []Response
	Objects     []string
}

type TempEndpoint struct {
	Url         string
	Description string
	HttpType    string
	Params      []Params
	Response    []Response
	Objects     []string
}

func IsController(line []byte) bool {
	return IsEp(string(line[:])) || strings.Contains(string(line[:]), springEndpointController)
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

func IsEp(line string) bool {
	return strings.Contains(line, endPointIdentifier) || strings.Contains(line, springMappingTagIdentifier)
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

func IsResponseNotation(line string) bool {
	return strings.Contains(line, responseIdentifier)
}
