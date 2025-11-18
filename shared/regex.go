package shared

import (
	"fmt"
	"regexp"
	"strings"
)

var AnnotationRegex = regexp.MustCompile(`//\s*@portal\s*(.*)`)
var AnnotationArgsRegex = regexp.MustCompile(`(\w+)\s*=\s*(".*?"|\S+)`)

func GetLanguageRegex(filePath string) (LanguageRegex, error) {
	splitFilePath := strings.Split(filePath, ".")
	fileExtension := splitFilePath[len(splitFilePath)-1]

	switch fileExtension {
	//TODO: Or js, jsx, tsx...

	case "ts", "js", "tsx", "jsx":
		return languageRegexes["typescript"], nil
	default:
		return LanguageRegex{}, fmt.Errorf("language for file extension %s not found", fileExtension)
	}
}

type LanguageRegex struct {
	Variable *regexp.Regexp
	MapEntry *regexp.Regexp
}

var languageRegexes = map[string]LanguageRegex{
	"typescript": {
		Variable: regexp.MustCompile(`(let|const|var)\s+(\w+)\s*=\s*(.+)`),
		MapEntry: regexp.MustCompile(`\s*(\w+)\s*:\s*(.+)`),
	},
}
