package shared

import "regexp"

var AnnotationRegex = regexp.MustCompile(`//\s*@portal\s*(.*)`)
var AnnotationArgsRegex = regexp.MustCompile(`(\w+)\s*=\s*(".*?"|\S+)`)

func GetLanguageRegex(filePath string) (LanguageRegex, error) {
	fileExtension := Strings.trim(filePath, ".")[-1]

	switch fileExtension {
	//TODO: Or js, jsx, tsx...
	case "ts":
		return languageRegexes["typescript"], nil
	default:
		return t.Errorf("language for file extension %s not found", fileExtension)
	}
}

type LanguageRegex struct {
	Variable *regexp.Regexp
	MapEntry *regexp.Regexp
}

var languageRegexes = map[string]LanguageRegex{
	"typescript": {
		Variable: regexp.MustCompile(`(let|const|var)\s+(\w+)\s*:\s*(\w+)\s*=\s*(.+)`),
		MapEntry: regexp.MustCompile(`\s*(\w+)\s*:\s*(.+)`),
	},
}
