package shared

import "regexp"

var AnnotationRegex = regexp.MustCompile(`//\s*@portal(?:\s+(.*))?`)

var AnnotationArgsRegex = regexp.MustCompile(`(\w+)\s*=\s*(".*?"|\S+)`)

var VariableRegex = regexp.MustCompile(`(let|const|var)\s+(\w+)\s*=\s*(.+)`)
