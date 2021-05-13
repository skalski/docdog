package javaLang

import (
	"docdog/src/constants"
	"docdog/src/helper"
	"docdog/src/notations"
	"errors"
	"fmt"
	"strings"
)

const Private = "private"
const Public = "public"
const Abstract = "abstract"
const Packg = "package"
const ImpTag = "import"
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
		for i <= 4 {
			if notations.IsNotNullNotation(wholeFile[index-i]) {
				tempVar.Notnull = true
			}
			if notations.IsDescriptionNotation(wholeFile[index-i]) {
				tempVar.Description = helper.GetStringFromQouteLine(wholeFile[index-i])
			}

			i++
		}
		return tempVar, nil
	}

	return tempVar, errors.New("is function")
}

func CreateApiEp(index int, wholeFile []string, verbose *bool, isSpringBoot bool, wl *int) notations.TempEndpoint {
	CheckLangTag(wholeFile, isSpringBoot)

	tempVar := notations.TempEndpoint{
		Url:         "",
		Description: "",
		HttpType:    "",
		Params:      nil,
		Objects:     nil,
	}
	i := index + 1

	for {
		*wl = i
		if i == len(wholeFile) {
			break
		}
		if notations.CommentEndTag(wholeFile[i]) {
			break
		}
		if notations.IsEp(wholeFile[i]) {
			helper.InfoLog(constants.LogMsgFoundUrl, wholeFile[i], verbose)
			tempVar.Url = helper.GetStringFromQouteLine(wholeFile[i])
		}
		if notations.IsDescriptionNotation(wholeFile[i]) {
			tempVar.Description = helper.GetStringFromQouteLine(wholeFile[i])
		}
		if notations.IsConnectionMethodNotation(wholeFile[i]) {
			helper.InfoLog(constants.LogMsgFoundConnectioType, wholeFile[i], verbose)
			tempPayload := helper.SeparateLineByTags(wholeFile[i])
			tempVar.HttpType = tempPayload[1]
		}
		if notations.IsParamNotation(wholeFile[i]) {
			helper.InfoLog(constants.LogMsgFoundParam, wholeFile[i], verbose)
			tempVar.Params = append(tempVar.Params, CreatePrm(wholeFile[i], verbose))
		}
		if notations.IsResponseNotation(wholeFile[i]) {
			helper.InfoLog(constants.LogMsgFoundResponse, wholeFile[i], verbose)
			tempVar.Response = append(tempVar.Response, CreateRes(wholeFile[i]))
		}
		if notations.IsPayloadNotation(wholeFile[i]) {
			helper.InfoLog(constants.LogMsgFoundPayload, wholeFile[i], verbose)
			tempPayload := helper.SeparateLineByTags(wholeFile[i])
			tempVar.Objects = append(tempVar.Objects, tempPayload[1])
		}
		i++
	}
	return tempVar
}

func CheckLangTag(wholeFile []string, isSpringBoot bool) bool {
	if !isSpringBoot && strings.Contains(strings.Join(wholeFile, ""), "springframework") {
		fmt.Println("WARING: Found SpringBoot you may should use -lang=spring ")
		return false
	}
	return true
}

func IsArrayType(line string) bool {
	return strings.Contains(line, arrayIdentifier) || strings.Contains(line, listIdentifier)
}

func CreateArrayType(line string) string {
	return strings.Replace(strings.Replace(strings.Replace(line, ">", "", 1), listIdentifier, "", 1), arrayIdentifier, "", 1)
}

func ChckImpl(fls []string) []string {
	imp := []string{}
	hasImp := false
	for _, s := range fls {
		if strings.Contains(s, Implements) {
			temp := helper.SeparateLineByTags(s)
			for _, s := range temp {
				if helper.LineEnd(s) {
					return imp
				}
				if hasImp {
					imp = append(imp, strings.Replace(s, ",", "", 1))
				}
				if strings.Contains(s, Implements) {
					hasImp = true
				}
			}
		}
	}
	return nil
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

func PackgName(fls []string) string {
	for _, s := range fls {
		if strings.Contains(s, Packg) {
			t := helper.SeparateLineByTags(s)
			return strings.Replace(t[1], ";", "", 1)
		}
	}
	return ""
}

func Imp(fls []string) []string {
	var imps []string
	for _, s := range fls {
		if strings.Contains(s, ImpTag) {
			t := helper.SeparateLineByTags(s)
			imps = append(imps, strings.Replace(t[1], ";", "", 1))
		}
	}
	imps = append(imps, PackgName(fls))
	return imps
}

func CreatePrm(l string, verbose *bool) notations.Params {
	helper.InfoLog(constants.LogMsgFoundParamForEndpoint, l, verbose)
	t := helper.SeparateLineByTags(l)
	param := &notations.Params{
		Name:        t[2],
		VarType:     t[1],
		Description: helper.GetStringFromQouteLine(strings.TrimSpace(l)),
		Notnull:     false,
	}

	if IsArray(t[1]) {
		param.IsArray = true
	}

	if notations.IsNotNullNotation(l) {
		param.Notnull = true
	}
	return *param
}

func CreateRes(l string) notations.Response {
	t := helper.SeparateLineByTags(l)
	rsp := &notations.Response{HttpCode: t[1]}
	rsp.Type = t[2]
	if t[2] == constants.Jsn && len(t) > 2 {
		rsp.ObjectType = t[3]
	}
	return *rsp
}

func IsArray(line string) bool {
	return strings.Contains(line, constants.ArrayIdentifier)
}
