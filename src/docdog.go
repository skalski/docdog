package main

import (
	"bufio"
	"docdog/src/constants"
	"docdog/src/helper"
	javaLang "docdog/src/languageParser/javalang"
	"docdog/src/languageParser/spring"
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
var verboseOutput *bool
var isSpringBoot = false
var wl = 0
var cd string

var objectList []notations.Objects
var abstractList []notations.Abstract
var endpointList []notations.Endpoint
var tempEndpointList []notations.TempEndpoint

func main() {
	defer RecoverPanic("parse panic")
	SetEnvironment()
	if !*verboseOutput {
		WelcomeMsg()
	}
	ScanFiles()
	Output()
	if !*verboseOutput {
		GoodbyeMsg()
	}
}

func SetEnvironment() {
	sourcePath = flag.String("path", "./", "set path of source.")
	outputPath = flag.String("out", "out.raml", "set file/path of the output file.")
	fileType = flag.String("lang", ".java", "Limit the type of file example: .java (spring||.php||.go||.rust)")
	verbose = flag.Bool("verbose", true, "Debug true/false")
	verboseOutput = flag.Bool("print", false, "Direct screen output for piping")

	tabset = flag.Int("tabs", 4, "lenght of tabs")
	flag.Parse()
	if *fileType == "spring" {
		isSpringBoot = true
		*fileType = ".java"
	}
	if !*verboseOutput {
		fmt.Printf("✓ Set filetype to: %s \n", *fileType)
		fmt.Printf("✓ Set path to: %s \n", *sourcePath)
		fmt.Printf("✓ Set output to: %s \n", *outputPath)
	}
}

func ScanFiles() {
	if !*verboseOutput {
		fmt.Println(constants.StartGatheringFiles)
	}
	GenerateModel(*sourcePath)
	if !*verboseOutput {
		fmt.Println(constants.FinishedGatheringFiles)
	}
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
	temp := helper.BytesToStringArrayByLinebreaks(b)
	cd = file.Name()
	// test if it's a possible controller or a object
	if notations.IsController(b) {
		helper.InfoLog(constants.LogMsgFoundController, file.Name(), verbose)
		for i, s := range temp {
			if notations.IsEp(s) {
				if !isSpringBoot {
					tempEndpointList = append(tempEndpointList, javaLang.CreateApiEp(i-1, temp, verbose, isSpringBoot, &wl))
				} else {
					tempEndpointList = append(tempEndpointList, spring.CreateApiEndpoint(i-1, temp))
				}
			}
		}
	} else {
		// if object is an interface - we ignore it
		if javaLang.IsItrf(temp) && *fileType == ".java" {
			return
		}

		helper.InfoLog(constants.LogMsgFoundPossibleObject, file.Name(), verbose)
		var variableList []notations.Variable
		for i, s := range temp {
			if IsJavaVariableOrFunctionEntry(s) {
				if !notations.HasIgnoreNotation(i, temp) {
					variable, err := CreateVariableStruct(s, i, temp)
					if err == nil {
						variableList = append(variableList, variable)
					}
				}
			}
		}
		if javaLang.IsAbstrc(temp) {
			// when Object is an abstract class, we store it in a special List
			helper.InfoLog(constants.LogMsgFoundAbstrct, file.Name(), verbose)
			abstractList = append(abstractList, notations.Abstract{
				Name:        ObjectNameBuilder(file.Name()),
				Variable:    variableList,
				PackageName: javaLang.PackgName(temp),
			})
		} else {
			objectList = append(objectList, notations.Objects{
				Name:        ObjectNameBuilder(file.Name()),
				Variable:    variableList,
				Implements:  javaLang.ChckImpl(temp),
				PackageName: javaLang.PackgName(temp),
				Imports:     javaLang.Imp(temp),
			})
		}
	}
}

func IsJavaVariableOrFunctionEntry(line string) bool {
	return strings.Contains(line, javaLang.Private) || strings.Contains(line, javaLang.Public)
}

func CreateVariableStruct(line string, index int, wholeFile []string) (notations.Variable, error) {
	if *fileType == ".java" && !isSpringBoot {
		return javaLang.JavaVariableHandler(line, index, wholeFile)
	}
	if isSpringBoot {
		return spring.SBVariableHandler(line, index, wholeFile)
	}
	return notations.Variable{}, errors.New(constants.NoMatchingLanguageMsg)
}

func Output() {
	if !*verboseOutput {
		fmt.Println(constants.StartCreatingAPIStructure)
	}
	GenEndpointsArrayStruc()
	fileBuilderRAML()
	if !*verboseOutput {
		fmt.Println(constants.FinishedCreatingAPIStructure)
	}
}

func GenEndpointsArrayStruc() {
	for _, e := range tempEndpointList {
		es := notations.Endpoint{
			Url:         e.Url,
			Description: e.Description,
			HttpType:    e.HttpType,
			Params:      nil,
			Variable:    nil,
			Objects:     nil,
		}
		es.Params = append(es.Params, e.Params...)
		es.Objects = append(es.Objects, e.Objects...)
		es.Response = append(es.Response, e.Response...)
		endpointList = append(endpointList, es)
	}
}

func fileBuilderRAML() {
	outputData := []string{
		"#%RAML 1.0",
		"---",
	}
	outputData = append(outputData, BuildRAMLObjects()...)
	outputData = append(outputData, BuildEndpoints()...)

	file, err := os.OpenFile(*outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	if !*verboseOutput {
		datawriter := bufio.NewWriter(file)
		for _, data := range outputData {
			_, _ = datawriter.WriteString(data + "\n")
		}

		_ = datawriter.Flush()
		_ = file.Close()
		return
	}
	for _, data := range outputData {
		fmt.Println(data)
	}

}

func BuildRAMLObjects() []string {
	outputData := []string{}
	if len(objectList) != 0 {

		outputData = append(outputData, constants.Types)
		for _, object := range objectList {
			for _, a := range abstractList {
				_, found := helper.Find(object.Imports, a.PackageName)
				// here we bring the additional Variables from the Abstracts back to the implementing Object
				for _, imp := range object.Implements {
					if a.Name == imp && found {
						object.Variable = append(object.Variable, a.Variable...)
					}
				}
			}

			outputData = append(outputData, StringWithTabs(1, object.Name+constants.Colon))
			outputData = append(outputData, StringWithTabs(2, "properties:"))
			for _, vars := range object.Variable {
				outputData = append(outputData, StringWithTabs(3, vars.Name+constants.Colon))
				if vars.IsArray {
					outputData = append(outputData, StringWithTabs(4, constants.TypeTag+"array"))
					outputData = append(outputData, StringWithTabs(4, "items:"))
					outputData = append(outputData, StringWithTabs(5, constants.TypeTag+vars.Typ))
				} else {
					outputData = append(outputData, StringWithTabs(4, constants.TypeTag+vars.Typ))
				}
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
		outputData = append(outputData, strings.Replace("/"+endpoint.Url+constants.Colon, "//", "/", 1))
		outputData = append(outputData, StringWithTabs(1, constants.DescriptionTag+endpoint.Description))
		outputData = append(outputData, StringWithTabs(1, endpoint.HttpType+constants.Colon))
		if len(endpoint.Params) != 0 {
			outputData = append(outputData, StringWithTabs(2, constants.QueryParamsTag))
			for _, param := range endpoint.Params {
				outputData = append(outputData, StringWithTabs(3, param.Name+constants.Colon))
				outputData = append(outputData, StringWithTabs(4, constants.DescriptionTag+param.Description))
				if param.Notnull {
					outputData = append(outputData, StringWithTabs(4, constants.RequiredTagTrue))
				}
				if param.IsArray {
					outputData = append(outputData, StringWithTabs(4, constants.ItemTag+"array"))
					outputData = append(outputData,
						StringWithTabs(5,
							constants.TypeTag+strings.Replace(param.VarType, constants.ArrayIdentifier, "", 1)),
					)
				} else {
					outputData = append(outputData, StringWithTabs(4, constants.ItemTag+"object"))
					outputData = append(outputData, StringWithTabs(5, constants.TypeTag+param.VarType))
				}
			}
		}
		if len(endpoint.Objects) != 0 {
			outputData = append(outputData, StringWithTabs(2, constants.BodyTag))
			outputData = append(outputData, StringWithTabs(3, constants.ApplicationJsonTag))
			for _, param := range endpoint.Objects {
				outputData = append(outputData, StringWithTabs(4, constants.SchemaTag+param))
			}
		}
		if len(endpoint.Response) != 0 {
			outputData = append(outputData, StringWithTabs(2, constants.Rsp))
			for _, rsp := range endpoint.Response {
				outputData = append(outputData, StringWithTabs(3, rsp.HttpCode+":"))
				outputData = append(outputData, StringWithTabs(4, constants.BodyTag))
				if strings.Contains(rsp.Type, constants.Txt) {
					outputData = append(outputData, StringWithTabs(5, constants.ApplicationDiverse))
				} else {
					outputData = append(outputData, StringWithTabs(5, constants.ApplicationJsonTag))
					outputData = append(outputData, StringWithTabs(6, constants.SchemaTag+rsp.ObjectType))
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
	fmt.Println(constants.DonationTag)
}

func RecoverPanic(msg string) {
	if r := recover(); r != nil {
		if msg != "" {
			if wl == 0 {
				fmt.Printf("There was a error while parsing on\n %s \n\n", cd)
			} else {
				fmt.Printf("There was a error while parsing at line %d on\n %s\n\n", wl, cd)
			}
		}
	}
}
