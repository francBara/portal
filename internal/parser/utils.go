package parser

import (
	"portal/shared"
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

func (args portalArguments) getGroup() string {
	group, _ := args.getString("group")
	if group != "" {
		return group
	}
	return "Default"
}

func parseAnnotationArguments(arguments string) portalArguments {
	matches := shared.AnnotationArgsRegex.FindAllStringSubmatch(arguments, -1)

	mappedArguments := make(map[string]string)

	for _, m := range matches {
		mappedArguments[m[1]] = strings.Trim(m[2], "\"")
	}

	return mappedArguments
}

func GetVariableType(value string) string {
	if value[0] == '"' && value[len(value)-1] == '"' || value[0] == '\'' && value[len(value)-1] == '\'' {
		return "string"
	} else if strings.Contains(value, ".") {
		return "float"
	} else {
		return "integer"
	}
}

//TODO: Modularize factories code

func numberVariableFactory(name string, value string, filePath string, options portalArguments) (shared.IntVariable, error) {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return shared.IntVariable{}, err
	}

	group := options.getGroup()

	maxValue, _ := options.getNum("max")
	minValue, _ := options.getNum("min")
	step, _ := options.getNum("step")

	displayName, _ := options.getString("name")
	if displayName == "" {
		displayName = name
	}

	return shared.IntVariable{
		PortalVariable: shared.PortalVariable{
			Name:        name,
			Group:       group,
			DisplayName: displayName,
			FilePath:    filePath,
		},
		Value: parsedValue,
		Max:   maxValue,
		Min:   minValue,
		Step:  step,
	}, nil
}

func floatVariableFactory(name string, value string, filePath string, options portalArguments) (shared.FloatVariable, error) {
	parsedValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return shared.FloatVariable{}, err
	}

	group := options.getGroup()

	maxValue, _ := options.getNum("max")
	minValue, _ := options.getNum("min")
	step, _ := options.getNum("step")

	displayName, _ := options.getString("name")
	if displayName == "" {
		displayName = name
	}

	return shared.FloatVariable{
		PortalVariable: shared.PortalVariable{
			Name:        name,
			DisplayName: displayName,
			Group:       group,
			FilePath:    filePath,
		},
		Value: float32(parsedValue),
		Max:   maxValue,
		Min:   minValue,
		Step:  step,
	}, nil
}

func stringVariableFactory(name string, value string, filePath string, options portalArguments) shared.StringVariable {
	value = strings.Trim(value, "\"'")

	group := options.getGroup()

	displayName, _ := options.getString("name")
	if displayName == "" {
		displayName = name
	}

	return shared.StringVariable{
		PortalVariable: shared.PortalVariable{
			Name:        name,
			DisplayName: displayName,
			Group:       group,
			FilePath:    filePath,
		},
		Value: value,
	}
}
