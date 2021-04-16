package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var sourcePath *string
var outputPath *string
var fileType *string
var verbose *bool
var tabset *int

var endPointIdentifier = "@DD:ENDPOINT"
var paramIdentifier = "@DD:PARAM"
var payloadIdentifier = "@DD:PAYLOAD"
var descriptionIdentifier = "@DD:DESCRIPTION"
var notNullIdentifier = "@DD:NOTNULL"
var typeIdentifier = "@DD:TYPE"
var ignoreIdentifier = "@DD:IGNORE"

const version = "0.1 ALPHA"
const fileReadIssue = "File Structure was changes during run or we run into a permission issue. Exit."

type Objects struct {
	name     string
	variable []Variable
}

type Variable struct {
	name        string
	description string
	typ         string
	notnull     bool
}

type Params struct {
	name        string
	description string
	notnull     bool
}

type Endpoint struct {
	url         string
	description string
	httpType    string
	params      []Params
	variable    []Variable
	objects     []string
}

type TempEndpoint struct {
	url         string
	description string
	httpType    string
	params      []Params
	objects     []string
}

var objectList []Objects
var paramsList []Params
var endpointList []Endpoint
var tempEndpointList []TempEndpoint

func main() {
	WelcomeMsg()
	SetEnvironment()
	ScanFiles()
	GenerateOutput()
	GoodbyeMsg()
}

func SetEnvironment() {
	sourcePath = flag.String("path", "./", "set path of source.")
	outputPath = flag.String("out", "out.raml", "set file/path of the output file.")
	fileType = flag.String("lang", ".java", "Limit the type of file example: .java (.php||.go||.rust)")
	verbose = flag.Bool("verbose", true, "Debug true/false")
	tabset = flag.Int("tabs", 4, "lenght of tabs")
	flag.Parse()

	fmt.Printf("✓ Set filetype to: %s \n", *fileType)
	fmt.Printf("✓ Set path to: %s \n", *sourcePath)
	fmt.Printf("✓ Set output to: %s \n", *outputPath)
}

func ScanFiles() {
	fmt.Println("- start gathering of Files ... please stand by!")
	GenerateModel(*sourcePath)
	fmt.Println("✓ finished with gathering of Files.")
}

func GenerateModel(path string) {
	err := filepath.Walk(path,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.Contains(info.Name(), *fileType) {
				file, err := os.Open(filePath)
				if err != nil {
					fmt.Println(fileReadIssue)
					log.Fatal(err)
				}
				AnalyseFile(file)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func AnalyseFile(file *os.File) {
	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Cannot open File")
	}
	temp := BytesToStringArrayByLinebreaks(b)
	if IsController(b) {
		InfoLog("Found Controller:", file.Name())
		for i, s := range temp {
			if IsEndpointNotation(s) {
				tempEndpointList = append(tempEndpointList, CreateApiEndpoint(i-1, temp))
			}
		}
	} else {
		InfoLog("Found possible Object: ", file.Name())
		var variableList []Variable
		for i, s := range temp {
			if IsJavaVariableOrFunctionentry(s) {
				if !HasIgnoreNotation(i, temp) {
					variable, err := CreateVariableStruct(s, i, temp)
					if err == nil {
						variableList = append(variableList, variable)
					}
				}
			}
		}
		objectList = append(objectList, Objects{
			name:     file.Name(),
			variable: variableList,
		})
	}
}

func IsJavaVariableOrFunctionentry(line string) bool {
	return strings.Contains(line, "private") || strings.Contains(line, "public")
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
	return strings.Contains(line, endPointIdentifier)
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

func CreateVariableStruct(line string, index int, wholeFile []string) (Variable, error) {
	temp := SeparateLineByTags(line)
	tempVar := Variable{
		name:        "",
		description: "",
		typ:         "",
		notnull:     false,
	}
	if !strings.Contains(line, "(") && !strings.Contains(line, "class") && !strings.Contains(line, "{") {
		tempVar.name = strings.ReplaceAll(temp[2], ";", "")
		tempVar.typ = temp[1]
		InfoLog("Found Variable: ", tempVar.name)
		InfoLog("Variable has type: ", tempVar.typ)
		i := 1
		for i <= 3 {
			if IsNotNullNotation(wholeFile[index-i]) {
				InfoLog("Found notNullTag for Variable: ", tempVar.name)
				tempVar.notnull = true
			}
			if IsDescriptionNotation(wholeFile[index-i]) {
				InfoLog("Found description for Variable: ", tempVar.name)
				tempVar.description = GetStringFromQouteLine(wholeFile[i])
			}

			i++
		}
		return tempVar, nil
	}
	return tempVar, errors.New("is function")
}

func CreateApiEndpoint(index int, wholeFile []string) TempEndpoint {
	tempVar := TempEndpoint{
		url:         "",
		description: "",
		httpType:    "",
		params:      nil,
		objects:     nil,
	}
	i := index + 1

	for {
		if CommentEndTag(wholeFile[i]) {
			break
		}
		if IsEndpointNotation(wholeFile[i]) {
			InfoLog("Found Url at :", wholeFile[i])
			tempVar.url = GetStringFromQouteLine(wholeFile[i])
		}
		if IsDescriptionNotation(wholeFile[i]) {
			tempVar.description = GetStringFromQouteLine(wholeFile[i])
		}
		if IsConnectionMethodNotation(wholeFile[i]) {
			InfoLog("Found Connection Type at :", wholeFile[i])
			tempPayload := SeparateLineByTags(wholeFile[i])
			tempVar.httpType = tempPayload[1]
		}
		if IsParamNotation(wholeFile[i]) {
			InfoLog("Found Param at :", wholeFile[i])
			tempVar.params = append(tempVar.params, CreateFromInstructionTag(wholeFile[i]))
		}
		if IsPayloadNotation(wholeFile[i]) {
			InfoLog("Found Payload at :", wholeFile[i])
			tempPayload := SeparateLineByTags(wholeFile[i])
			tempVar.objects = append(tempVar.objects, tempPayload[1])
		}
		i++
	}
	return tempVar
}

func CreateFromInstructionTag(line string) Params {
	InfoLog("Produced Param for current Endpoint :", line)
	temp := SeparateLineByTags(line)
	param := &Params{
		name:        temp[1],
		description: GetStringFromQouteLine(strings.TrimSpace(line)),
		notnull:     false,
	}

	if IsNotNullNotation(line) {
		param.notnull = true
	}
	return *param
}

func SeparateLineByTags(line string) []string {
	return strings.Split(strings.TrimSpace(line), " ")
}
func GenerateOutput() {
	fmt.Println("- start creating of API structure ... please stand by!")
	GenerateEndpointsArrayStructure()
	fileBuilderRAML()
	fmt.Println("✓ finished with creating. Thanks for using DogDoc")
}

func GenerateEndpointsArrayStructure() {
	for _, tempEndpoint := range tempEndpointList {
		endpoint := Endpoint{
			url:         tempEndpoint.url,
			description: tempEndpoint.description,
			httpType:    tempEndpoint.httpType,
			params:      nil,
			variable:    nil,
			objects:     nil,
		}
		for _, params := range tempEndpoint.params {
			endpoint.params = append(endpoint.params, params)
		}
		for _, object := range tempEndpoint.objects {
			endpoint.objects = append(endpoint.objects, object)
		}
		endpointList = append(endpointList, endpoint)
	}
}

func GetStringFromQouteLine(str string) (result string) {
	s := strings.Index(str, "'")
	if s == -1 {
		return
	}
	s += len("'")
	e := strings.Index(str[s:], "'")
	if e == -1 {
		return
	}
	return str[s : s+e]
}

func BytesToStringArrayByLinebreaks(data []byte) []string {
	return strings.Split(strings.ReplaceAll(string(data[:]), "\r\n", "\n"), "\n")
}

func InfoLog(msg string, source string) {
	if !*verbose {
		log.Println(msg + source)
	}
}

func fileBuilderRAML() {
	outputData := []string{
		"#RAML1.0",
	}
	outputData = append(outputData, BuildRAMLObjects()...)
	outputData = append(outputData, BuildEndpoints()...)

	file, err := os.OpenFile(*outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)
	for _, data := range outputData {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	file.Close()
}

func BuildRAMLObjects() []string {
	outputData := []string{}
	if len(objectList) != 0 {
		outputData = append(outputData, "types:")
		for _, object := range objectList {
			outputData = append(outputData, StringWithTabs(1, ObjectNameBuilder(object.name)+":"))
			outputData = append(outputData, StringWithTabs(2, "type: object"))
			for _, vars := range object.variable {
				outputData = append(outputData, StringWithTabs(3, vars.name+":"))
				outputData = append(outputData, StringWithTabs(4, "type:"+vars.typ))
				if vars.notnull {
					outputData = append(outputData, StringWithTabs(4, "required: true"))
				}
			}
		}
	}
	return outputData
}

func BuildEndpoints() []string {
	outputData := []string{}
	for _, endpoint := range endpointList {
		outputData = append(outputData, endpoint.url+":")
		outputData = append(outputData, StringWithTabs(1, "description: "+endpoint.description))
		outputData = append(outputData, StringWithTabs(2, endpoint.httpType+":"))
		if len(endpoint.params) != 0 {
			outputData = append(outputData, StringWithTabs(3, "queryParameters:"))
			for _, param := range endpoint.params {
				outputData = append(outputData, StringWithTabs(4, param.name+":"))
				outputData = append(outputData, StringWithTabs(5, "description: "+param.description))
				if param.notnull {
					outputData = append(outputData, StringWithTabs(5, "required: true"))
				}
			}
		}
	}
	return outputData
}

func ObjectNameBuilder(filename string) string {
	parts := strings.Split(strings.TrimSpace(filename), "/")
	if len(parts) <= 2 {
		parts = strings.Split(strings.TrimSpace(filename), "\\")
	}
	objectName := parts[len(parts)-1]
	return strings.Replace(objectName, *fileType, "", -1)
}

func Tab() string {
	if *tabset == 0 {
		return "\t"
	}
	tabs := []string{}
	i := 0
	for i <= *tabset-1 {
		tabs = append(tabs, " ")
		i++
	}
	return strings.Join(tabs, "")
}

func StringWithTabs(count int, text string) string {
	tabs := []string{}
	i := 0
	for i <= count-1 {
		tabs = append(tabs, Tab())
		i++
	}
	return strings.Join(tabs, "") + text
}

func WelcomeMsg() {
	fmt.Println("     --- DocDog ---")
	fmt.Printf("      version:%s\n", version)
	fmt.Print(" written by Swen Kalski\n\n\n")
}

func GoodbyeMsg() {
	fmt.Println("")
	fmt.Println("^..^      /")
	fmt.Println("/_/\\_____/")
	fmt.Println("   /\\   /\\")
	fmt.Println("  /  \\ /  \\")
	fmt.Println("")
	fmt.Println("Thanks for using DocDog")
}
