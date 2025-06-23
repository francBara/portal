package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"portal/shared"
	"regexp"
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

func (args portalArguments) getString(key string) string {
	var value string
	var ok bool

	if value, ok = args[key]; ok {
		return value
	}
	return ""
}

// parseAnnotationArguments parses a @portal annotation optional arguments, returning the corresponding struct.
func parseAnnotationArguments(arguments string) portalArguments {
	mappedArguments := make(map[string]string)

	// Positional arguments
	splitArguments := strings.Fields(arguments)
	for _, arg := range splitArguments {
		if arg == "all" {
			mappedArguments["all"] = "true"
		} else if arg == "ui" {
			mappedArguments["ui"] = "true"
		}
	}

	// Named arguments
	matches := shared.AnnotationArgsRegex.FindAllStringSubmatch(arguments, -1)

	for _, m := range matches {
		mappedArguments[m[1]] = strings.Trim(m[2], "\"")
	}

	return mappedArguments
}

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

// getPortalVariable generates PortalVariable base data basing on arguments, name and filePath.
func (arguments portalArguments) getPortalVariable(name string, filePath string) shared.PortalVariable {
	group := arguments.getString("group")
	if group == "" {
		group = "Default"
	}

	displayName := arguments.getString("name")
	if displayName == "" {
		displayName = name
	}

	view := arguments.getString("view")
	if view == "" {
		view = strings.Split(filepath.Base(filePath), ".")[0]
	}

	return shared.PortalVariable{
		Name:        name,
		DisplayName: displayName,
		View:        view,
		Group:       group,
	}
}

func numberVariableFactory(name string, value string, filePath string, options portalArguments) (shared.IntVariable, error) {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return shared.IntVariable{}, err
	}

	maxValue, _ := options.getNum("max")
	minValue, _ := options.getNum("min")
	step, _ := options.getNum("step")

	return shared.IntVariable{
		PortalVariable: options.getPortalVariable(name, filePath),
		Value:          parsedValue,
		Max:            maxValue,
		Min:            minValue,
		Step:           step,
	}, nil
}

func floatVariableFactory(name string, value string, filePath string, options portalArguments) (shared.FloatVariable, error) {
	parsedValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return shared.FloatVariable{}, err
	}

	maxValue, _ := options.getNum("max")
	minValue, _ := options.getNum("min")
	step, _ := options.getNum("step")

	return shared.FloatVariable{
		PortalVariable: options.getPortalVariable(name, filePath),
		Value:          float32(parsedValue),
		Max:            maxValue,
		Min:            minValue,
		Step:           step,
	}, nil
}

func stringVariableFactory(name string, value string, filePath string, options portalArguments) shared.StringVariable {
	value = strings.Trim(value, "\"'")

	return shared.StringVariable{
		PortalVariable: options.getPortalVariable(name, filePath),
		Value:          value,
	}
}

// uiVariablesFactory uses generateTree.js tool to get a tree representation of the html tree.
func uiVariablesFactory(basePath string, filePath string, options portalArguments) (map[string]shared.UIVariable, error) {
	cmd := exec.Command("node", "tools/generateTree.js", fmt.Sprintf("%s/%s", basePath, filePath))

	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var roots map[string]shared.UIVariable

	if err := json.Unmarshal(out.Bytes(), &roots); err != nil {
		fmt.Println("JSON parse error:", err)
		return nil, err
	}

	// PortalVariable data and IDs are added after tree generation basing on arguments
	for rootName, root := range roots {
		root.PortalVariable = options.getPortalVariable(rootName, filePath)
		roots[rootName] = root
	}

	return roots, nil
}
