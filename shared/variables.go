package shared

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"maps"
	"os"
	"strconv"
)

type PortalVariable struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Group       string `json:"group"`
	FilePath    string `json:"filePath"`
}

type IntVariable struct {
	PortalVariable
	Value int `json:"value"`
	Max   int `json:"max"`
	Min   int `json:"min"`
	Step  int `json:"step"`
}

func (variable IntVariable) update(value string) (IntVariable, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return variable, errors.New("could not patch variables")
	}
	variable.Value = parsed
	return variable, nil
}

type FloatVariable struct {
	PortalVariable
	Value float32 `json:"value"`
	Max   int     `json:"max"`
	Min   int     `json:"min"`
	Step  int     `json:"step"`
}

func (variable FloatVariable) update(value string) (FloatVariable, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return variable, errors.New("could not patch variables")
	}
	variable.Value = float32(parsed)
	return variable, nil
}

type StringVariable struct {
	PortalVariable
	Value string `json:"value"`
}

func (variable StringVariable) update(value string) StringVariable {
	variable.Value = value
	return variable
}

type PortalVariables struct {
	Integer    map[string]IntVariable    `json:"integer"`
	Float      map[string]FloatVariable  `json:"float"`
	String     map[string]StringVariable `json:"string"`
	FileHashes map[string]string         `json:"fileHashes"`
}

func (variables PortalVariables) DumpVariables() {
	file, err := os.Create("variables.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(variables); err != nil {
		panic(err)
	}
}

func (variables PortalVariables) UpdateVariables(variablesPatch map[string]string) (PortalVariables, error) {
	for key, value := range variablesPatch {
		if _, ok := variables.Integer[key]; ok {
			newVar, err := variables.Integer[key].update(value)
			if err != nil {
				return variables, errors.New("could not patch variables")
			}
			variables.Integer[key] = newVar
		}

		if _, ok := variables.Float[key]; ok {
			newVar, err := variables.Float[key].update(value)
			if err != nil {
				return variables, errors.New("could not patch variables")
			}
			variables.Float[key] = newVar
		}

		if _, ok := variables.String[key]; ok {
			variables.String[key] = variables.String[key].update(value)
		}
	}

	return variables, nil
}

func (variables PortalVariables) HasVariables() bool {
	return len(variables.Integer) > 0 || len(variables.String) > 0 || len(variables.Float) > 0
}

func mergeMaps[K comparable, v any](map1 map[K]v, map2 map[K]v) map[K]v {
	newMap := make(map[K]v)

	maps.Copy(newMap, map1)
	maps.Copy(newMap, map2)

	return newMap
}

func (variables PortalVariables) Merge(newVariables PortalVariables) PortalVariables {
	var merged PortalVariables

	merged.Integer = mergeMaps(variables.Integer, newVariables.Integer)
	merged.Float = mergeMaps(variables.Float, newVariables.Float)
	merged.String = mergeMaps(variables.String, newVariables.String)
	merged.FileHashes = mergeMaps(variables.FileHashes, newVariables.FileHashes)

	return merged
}

func (variables PortalVariables) HasFileChanged(fileContent string, filePath string) bool {
	hasher := sha256.New()
	hasher.Write([]byte(fileContent))
	hashString := hex.EncodeToString(hasher.Sum(nil))

	return hashString != variables.FileHashes[filePath]
}

func (variables PortalVariables) Length() int {
	return len(variables.Integer) + len(variables.Float) + len(variables.String)
}

type VariablesMap map[string]map[string]any

func (varsMap *VariablesMap) add(group string, name string, value any) {
	if _, ok := (*varsMap)[group]; !ok {
		(*varsMap)[group] = make(map[string]any)
	}

	(*varsMap)[group][name] = value
}

// Converts PortalVariables struct to a hash map containing groups as keys, variable names as subkeys and variables as values
func (variables PortalVariables) ToMap() VariablesMap {
	mappedVariables := make(VariablesMap)

	for _, intVar := range variables.Integer {
		mappedVariables.add(intVar.Group, intVar.Name, map[string]any{
			"displayName": intVar.DisplayName,
			"value":       intVar.Value,
			"max":         intVar.Max,
			"min":         intVar.Min,
			"step":        intVar.Step,
			"filePath":    intVar.FilePath,
			"type":        "integer",
		})
	}

	for _, floatVar := range variables.Float {
		mappedVariables.add(floatVar.Group, floatVar.Name, map[string]any{
			"displayName": floatVar.DisplayName,
			"value":       floatVar.Value,
			"max":         floatVar.Max,
			"min":         floatVar.Min,
			"step":        floatVar.Step,
			"filePath":    floatVar.FilePath,
			"type":        "float",
		})
	}

	for _, stringVar := range variables.String {
		mappedVariables.add(stringVar.Group, stringVar.Name, map[string]any{
			"displayName": stringVar.DisplayName,
			"value":       stringVar.Value,
			"filePath":    stringVar.FilePath,
			"type":        "string",
		})
	}

	return mappedVariables
}
