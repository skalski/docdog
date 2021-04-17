package main

import (
	"bufio"
	"docdog/src/constants"
	"docdog/src/helper"
	"docdog/src/javalang"
	"docdog/src/notations"
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

var objectList []notations.Objects
var paramsList []notations.Params
var endpointList []notations.Endpoint
var tempEndpointList []notations.TempEndpoint

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
	fileType = flag.String("javalang", ".java", "Limit the type of file example: .java (.php||.go||.rust)")
	verbose = flag.Bool("verbose", true, "Debug true/false")
	tabset = flag.Int("tabs", 4, "lenght of tabs")
	flag.Parse()

	fmt.Printf("✓ Set filetype to: %s \n", *fileType)
	fmt.Printf("✓ Set path to: %s \n", *sourcePath)
	fmt.Printf("✓ Set output to: %s \n", *outputPath)
}

func ScanFiles() {
	fmt.Println(constants.StartGatheringFiles)
	GenerateModel(*sourcePath)
	fmt.Println(constants.FinishedGatheringFiles)
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
					fmt.Println(constants.FileReadIssue)
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
		fmt.Println(constants.CannotOpenFiles)
	}
	temp := BytesToStringArrayByLinebreaks(b)
	if notations.IsController(b) {
		InfoLog(constants.LogMsgFoundController, file.Name())
		for i, s := range temp {
			if notations.IsEndpointNotation(s) {
				tempEndpointList = append(tempEndpointList, CreateApiEndpoint(i-1, temp))
			}
		}
	} else {
		InfoLog(constants.LogMsgFoundPossibleObject, file.Name())
		var variableList []notations.Variable
		for i, s := range temp {
			if IsJavaVariableOrFunctionentry(s) {
				if !notations.HasIgnoreNotation(i, temp) {
					variable, err := CreateVariableStruct(s, i, temp)
					if err == nil {
						variableList = append(variableList, variable)
					}
				}
			}
		}
		objectList = append(objectList, notations.Objects{
			Name:     file.Name(),
			Variable: variableList,
		})
	}
}

func IsJavaVariableOrFunctionentry(line string) bool {
	return strings.Contains(line, javalang.Private) || strings.Contains(line, javalang.Public)
}

func CreateVariableStruct(line string, index int, wholeFile []string) (notations.Variable, error) {
	if *fileType == ".java" {
		return javalang.JavaVariableHandler(line, index, wholeFile)
	}

	return notations.Variable{}, errors.New(constants.NoMatchingLanguageMsg)
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
			InfoLog(constants.LogMsgFoundUrl, wholeFile[i])
			tempVar.Url = helper.GetStringFromQouteLine(wholeFile[i])
		}
		if notations.IsDescriptionNotation(wholeFile[i]) {
			tempVar.Description = helper.GetStringFromQouteLine(wholeFile[i])
		}
		if notations.IsConnectionMethodNotation(wholeFile[i]) {
			InfoLog(constants.LogMsgFoundConnectioType, wholeFile[i])
			tempPayload := helper.SeparateLineByTags(wholeFile[i])
			tempVar.HttpType = tempPayload[1]
		}
		if notations.IsParamNotation(wholeFile[i]) {
			InfoLog(constants.LogMsgFoundParam, wholeFile[i])
			tempVar.Params = append(tempVar.Params, CreateFromInstructionTag(wholeFile[i]))
		}
		if notations.IsPayloadNotation(wholeFile[i]) {
			InfoLog(constants.LogMsgFoundPayload, wholeFile[i])
			tempPayload := helper.SeparateLineByTags(wholeFile[i])
			tempVar.Objects = append(tempVar.Objects, tempPayload[1])
		}
		i++
	}
	return tempVar
}

func CreateFromInstructionTag(line string) notations.Params {
	InfoLog(constants.LogMsgFoundParamForEndpoint, line)
	temp := helper.SeparateLineByTags(line)
	param := &notations.Params{
		Name:        temp[1],
		Description: helper.GetStringFromQouteLine(strings.TrimSpace(line)),
		Notnull:     false,
	}

	if notations.IsNotNullNotation(line) {
		param.Notnull = true
	}
	return *param
}

func GenerateOutput() {
	fmt.Println(constants.StartCreatingAPIStructure)
	GenerateEndpointsArrayStructure()
	fileBuilderRAML()
	fmt.Println(constants.FinishedCreatingAPIStructure)
}

func GenerateEndpointsArrayStructure() {
	for _, tempEndpoint := range tempEndpointList {
		endpoint := notations.Endpoint{
			Url:         tempEndpoint.Url,
			Description: tempEndpoint.Description,
			HttpType:    tempEndpoint.HttpType,
			Params:      nil,
			Variable:    nil,
			Objects:     nil,
		}
		for _, params := range tempEndpoint.Params {
			endpoint.Params = append(endpoint.Params, params)
		}
		for _, object := range tempEndpoint.Objects {
			endpoint.Objects = append(endpoint.Objects, object)
		}
		endpointList = append(endpointList, endpoint)
	}
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

	_ = datawriter.Flush()
	_ = file.Close()
}

func BuildRAMLObjects() []string {
	outputData := []string{}
	if len(objectList) != 0 {
		outputData = append(outputData, constants.Types)
		for _, object := range objectList {
			outputData = append(outputData, StringWithTabs(1, ObjectNameBuilder(object.Name)+constants.Colon))
			outputData = append(outputData, StringWithTabs(2, constants.TypeObject))
			for _, vars := range object.Variable {
				outputData = append(outputData, StringWithTabs(3, vars.Name+constants.Colon))
				outputData = append(outputData, StringWithTabs(4, constants.TypeTag+vars.Typ))
				if vars.Notnull {
					outputData = append(outputData, StringWithTabs(4, constants.RequiredTagTrue))
				}
			}
		}
	}
	return outputData
}

func BuildEndpoints() []string {
	outputData := []string{}
	for _, endpoint := range endpointList {
		outputData = append(outputData, endpoint.Url+constants.Colon)
		outputData = append(outputData, StringWithTabs(1, constants.DescriptionTag+endpoint.Description))
		outputData = append(outputData, StringWithTabs(2, endpoint.HttpType+constants.Colon))
		if len(endpoint.Params) != 0 {
			outputData = append(outputData, StringWithTabs(3, constants.QueryParamsTag))
			for _, param := range endpoint.Params {
				outputData = append(outputData, StringWithTabs(4, param.Name+constants.Colon))
				outputData = append(outputData, StringWithTabs(5, constants.DescriptionTag+param.Description))
				if param.Notnull {
					outputData = append(outputData, StringWithTabs(5, constants.RequiredTagTrue))
				}
			}
		}
		if len(endpoint.Objects) != 0 {
			outputData = append(outputData, StringWithTabs(3, constants.BodyTag))
			outputData = append(outputData, StringWithTabs(4, constants.ApplicationJsonTag))
			outputData = append(outputData, StringWithTabs(5, constants.AmfAdditionalProperties))
			for _, param := range endpoint.Objects {
				outputData = append(outputData, StringWithTabs(6, constants.TypeTag+param))
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
	fmt.Printf("      version:%s\n", constants.Version)
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
