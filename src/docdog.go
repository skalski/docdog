package main

import (
	"encoding/json"
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

var endPointIdentifier = "@DD:ENDPOINT"
var paramIdentifier = "@DD:PARAM"
var payloadIdentifier = "@DD:PAYLOAD"
var descriptionIdentifier = "@DD:DESCRIPTION"
var notNullIdentifier = "@DD:NOTNULL"
var typeIdentifier = "@DD:TYPE"

const version = "0.1 ALPHA"
const fileReadIssue = "File Structure was changes during run or we run into a permission issue. Exit."

type Objects struct {
	name     string
	variable []Variable
}

type Variable struct {
	name        string
	description string
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
	outputPath = flag.String("out", "out.rml", "set file/path of the output file.")
	fileType = flag.String("lang", ".java", "Limit the type of file example: .java (.php||.go||.rust)")
	verbose = flag.Bool("verbose", true, "Debug true/false")
	flag.Parse()

	fmt.Printf("✓ Set filetype to: %s \n", *fileType)
	fmt.Printf("✓ Set path to: %s \n", *sourcePath)
	fmt.Printf("✓ Set output to: %s \n", *outputPath)
}

func ScanFiles() {
	fmt.Println("- start gathering of Files ... please stand by!")
	err := filepath.Walk(*sourcePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			GenerateModel(path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	fmt.Println("✓ finished with gathering of Files.")
}

func GenerateModel(path string) {
	files, _ := ioutil.ReadDir(path)
	for _, fileFromDir := range files {
		if !fileFromDir.IsDir() && strings.Contains(fileFromDir.Name(), *fileType) {
			file, err := os.Open(path + fileFromDir.Name())
			if err != nil {
				fmt.Println(fileReadIssue)
				log.Fatal(err)
			}
			AnalyseFile(file)
		}
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
				tempEndpointList = append(tempEndpointList, CreateApiEndpoint(i, temp))
			}
		}
	} else {
		var variableList []Variable
		for i, s := range temp {
			if IsJavaVariable(s) {
				variable, err := CreateVariableStruct(s, i, temp)
				if err == nil {
					variableList = append(variableList, variable)
				}
			}
		}
		objectList = append(objectList, Objects{
			name:     file.Name(),
			variable: variableList,
		})
	}
}

func IsJavaVariable(line string) bool {
	temp := strings.Split(strings.ReplaceAll(line, "\t", ""), " ")
	return strings.Contains(line, "private") && len(temp) == 2 || strings.Contains(line, "private") && len(temp) == 2
}

func IsController(line []byte) bool {
	return IsEndpointNotation(string(line[:]))
}

func IsDocDogNotation(line string) bool {
	return IsEndpointNotation(line) || IsParamNotation(line) || IsPayloadNotation(line) || IsNotNullNotation(line) || IsDescriptionNotation(line)
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

func CommentEndTag(line string) bool {
	return strings.Contains(line, "*/")
}

func CreateVariableStruct(line string, index int, wholeFile []string) (Variable, error) {
	temp := SeparateLineByTags(line)
	tempVar := Variable{
		name:        "",
		description: "",
		notnull:     false,
	}
	if !strings.Contains(line, "(") {
		tempVar.name = strings.ReplaceAll(temp[3], ";", "")
		i := 1
		for i <= 3 {
			if IsDocDogNotation(wholeFile[index-i]) {
				if IsNotNullNotation(wholeFile[index-i]) {
					tempVar.notnull = true
				}
				if IsDescriptionNotation(wholeFile[index-i]) {
					tempVar.description = GetStringFromQouteLine(wholeFile[i])
				}
			}
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
		if IsDescriptionNotation(wholeFile[i]) {
			tempVar.description = GetStringFromQouteLine(wholeFile[i])
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
	b, _ := json.Marshal(tempVar)
	InfoLog("Produced Endpoint :", string(b))
	return tempVar
}

func CreateFromInstructionTag(line string) Params {
	InfoLog("Produced Param for current Endpoint :", line)
	temp := SeparateLineByTags(line)
	fmt.Println(strings.Join(temp, ", "))
	param := &Params{
		name:        temp[1],
		description: GetStringFromQouteLine(strings.TrimSpace(line)),
		notnull:     false,
	}

	if IsNotNullNotation(line) {
		param.notnull = true
	}
	b, _ := json.Marshal(param)
	InfoLog("Produced Param for current Endpoint :", string(b))
	return *param
}

func SeparateLineByTags(line string) []string {
	return strings.Split(strings.TrimSpace(line), " ")
}
func GenerateOutput() {
	fmt.Println("- start creating of API structure ... please stand by!")
	GenerateArrayStructure()

	fmt.Println("✓ finished with creating. Thanks for using DogDoc")
}

func GenerateArrayStructure() {
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
			CreateObjectJSON()
			endpoint.objects = append(endpoint.objects, object)
		}
		endpointList = append(endpointList, endpoint)
	}
}

func CreateObjectJSON() {
	fmt.Printf("endpoints: %q\n", tempEndpointList)
	fmt.Printf("objects: %q\n", objectList)
	fmt.Printf("objectVars: %q\n", objectList[0].variable)
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
