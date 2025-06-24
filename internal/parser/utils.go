package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"portal/internal/parser/annotation"
	"portal/shared"
	"regexp"
	"strconv"
	"strings"
)

// GetVariableType returns the variable type basing on how its value is defined.
func GetVariableType(value string) string {
	if value[0] == '"' && value[len(value)-1] == '"' || value[0] == '\'' && value[len(value)-1] == '\'' {
		return "string"
	} else if strings.Contains(value, ".") {
		return "float"
	} else {
		return "integer"
	}
}

func IsTailwindLine(line string) bool {
	return strings.Contains(line, "-")
}

// Parses a tailwind line returning parameter name and numeric value
func ParseTailwindLine(line string) (string, string) {
	line = strings.TrimSpace(line)
	valueIdx := strings.LastIndex(line, "-")

	varName := line[:valueIdx]
	value := regexp.MustCompile(`\D`).ReplaceAllString(line[valueIdx+1:], "")

	return varName, value
}

func numberVariableFactory(name string, value string, filePath string, options annotation.PortalAnnotation) (shared.IntVariable, error) {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return shared.IntVariable{}, err
	}

	return shared.IntVariable{
		PortalVariable: options.GetPortalVariable(name, filePath),
		Value:          parsedValue,
		Max:            options.Max,
		Min:            options.Min,
		Step:           options.Step,
	}, nil
}

func floatVariableFactory(name string, value string, filePath string, options annotation.PortalAnnotation) (shared.FloatVariable, error) {
	parsedValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return shared.FloatVariable{}, err
	}

	return shared.FloatVariable{
		PortalVariable: options.GetPortalVariable(name, filePath),
		Value:          float32(parsedValue),
		Max:            options.Max,
		Min:            options.Min,
		Step:           options.Step,
	}, nil
}

func stringVariableFactory(name string, value string, filePath string, options annotation.PortalAnnotation) shared.StringVariable {
	value = strings.Trim(value, "\"'")

	return shared.StringVariable{
		PortalVariable: options.GetPortalVariable(name, filePath),
		Value:          value,
	}
}

// uiVariablesFactory uses generateTree.js tool to get a tree representation of the html tree.
func uiVariablesFactory(basePath string, filePath string) (map[string]shared.UIVariable, error) {
	cmd := exec.Command("node", "tools/generateTree.js", filepath.Join(basePath, filePath))

	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var result struct {
		Components map[string]shared.UIVariable `json:"components"`
		Props      map[string][]string          `json:"props"`
		Comments   map[string][]string          `json:"comments"`
	}

	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		fmt.Println("JSON parse error:", err)
		return nil, err
	}

	// PortalVariable data and IDs are added after tree generation basing on arguments
	for rootName, root := range result.Components {
		ann, err := annotation.ParseAnnotation(strings.Join(result.Comments[rootName], " "))
		if err != nil {
			return nil, err
		}

		if len(result.Props[rootName]) < len(ann.Mocks) {
			return nil, fmt.Errorf("too many mocks, mocks: %d, props: %d", len(ann.Mocks), len(result.Props[rootName]))
		}

		root.PropsMocks = make(map[string]string)

		for i := range ann.Mocks {
			root.PropsMocks[result.Props[rootName][i]] = ann.Mocks[i]
		}

		root.Box = ann.Box

		root.PortalVariable = ann.GetPortalVariable(rootName, filePath)
		result.Components[rootName] = root
	}

	return result.Components, nil
}
