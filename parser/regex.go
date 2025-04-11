package parser

import "regexp"

var AnnotationRegex = regexp.MustCompile(`//\s*@portal\s+(.*)`)
var VariableRegex = regexp.MustCompile(`(let|const|var)\s+(\w+)\s*=\s*(.+?);`)
