package parser

import (
	"portal/internal/parser/annotation"
	"portal/shared"
	"strconv"
	"strings"
)

// GetVariableType returns the variable type basing on how its value is defined.
func GetVariableType(value string) string {
	if value[0] == '"' && value[len(value)-1] == '"' || value[0] == '\'' && value[len(value)-1] == '\'' {
		return "string"
	}

	if _, err := strconv.Atoi(value); err == nil {
		return "integer"
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return "float"
	}

	return ""
}

func numberVariableFactory(name string, value string, filePath string, lineNumber int, options annotation.PortalAnnotation) (shared.IntVariable, error) {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return shared.IntVariable{}, err
	}

	return shared.IntVariable{
		PortalVariable: options.GetPortalVariable(name, filePath, lineNumber),
		Value:          parsedValue,
		Max:            options.Max,
		Min:            options.Min,
		Step:           options.Step,
	}, nil
}

func floatVariableFactory(name string, value string, filePath string, lineNumber int, options annotation.PortalAnnotation) (shared.FloatVariable, error) {
	parsedValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return shared.FloatVariable{}, err
	}

	return shared.FloatVariable{
		PortalVariable: options.GetPortalVariable(name, filePath, lineNumber),
		Value:          float32(parsedValue),
		Max:            options.Max,
		Min:            options.Min,
		Step:           options.Step,
	}, nil
}

func stringVariableFactory(name string, value string, filePath string, lineNumber int, options annotation.PortalAnnotation) shared.StringVariable {
	value = strings.Trim(value, "\"'")

	return shared.StringVariable{
		PortalVariable: options.GetPortalVariable(name, filePath, lineNumber),
		Value:          value,
	}
}
