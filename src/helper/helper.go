package helper

import (
	"log"
	"strings"
)

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

func InfoLog(msg string, source string, verbose *bool) {
	if !*verbose {
		log.Println(msg + source)
	}
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if strings.Contains(item, val) {
			return i, true
		}
	}
	return -1, false
}

func LineEnd(s string) bool {
	return strings.Contains(s, "{")
}
