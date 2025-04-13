package parser

import (
	"strconv"
	"strings"
)

type portalArguments map[string]string

func (args portalArguments) getNum(key string) (int, bool) {
	var parsedValue int
	var err error

	if strValue, ok := args[key]; ok {
		parsedValue, err = strconv.Atoi(strValue)
		if err != nil {
			return 0, false
		}
		return parsedValue, true
	}

	return 0, false
}

func (args portalArguments) getString(key string) (string, bool) {
	var value string
	var ok bool

	if value, ok = args[key]; ok {
		return value, true
	}
	return "", false
}

func parseAnnotationArguments(arguments string) portalArguments {
	matches := AnnotationArgsRegex.FindAllStringSubmatch(arguments, -1)

	mappedArguments := make(map[string]string)

	for _, m := range matches {
		mappedArguments[m[1]] = m[2]
	}

	return mappedArguments
}

func getVariableType(value string) string {
	if value[0] == '"' && value[len(value)-1] == '"' || value[0] == '\'' && value[len(value)-1] == '\'' {
		return "string"
	} else {
		return "number"
	}
}

func numberVariableFactory(name string, value string, filePath string, options portalArguments) (NumberVariable, error) {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return NumberVariable{}, err
	}

	group, _ := options.getString("group")

	maxValue, _ := options.getNum("max")
	minValue, _ := options.getNum("min")
	step, _ := options.getNum("step")

	return NumberVariable{
		PortalVariable: PortalVariable{
			Name:     name,
			Group:    group,
			FilePath: filePath,
		},
		Value: parsedValue,
		Max:   maxValue,
		Min:   minValue,
		Step:  step,
	}, nil
}

func stringVariableFactory(name string, value string, options portalArguments) StringVariable {
	value = strings.Trim(value, "\"'")

	group, _ := options.getString("group")

	return StringVariable{
		PortalVariable: PortalVariable{
			Name:  name,
			Group: group,
		},
		Value: value,
	}
}
