package helper

import "strings"

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

func SeparateLineByTags(line string) []string {
	return strings.Split(strings.TrimSpace(line), " ")
}

func BytesToStringArrayByLinebreaks(data []byte) []string {
	return strings.Split(strings.ReplaceAll(string(data[:]), "\r\n", "\n"), "\n")
}
