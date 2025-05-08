package shared

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"maps"
	"os"
)

type PortalVariable struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	View        string `json:"view"`
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

type FloatVariable struct {
	PortalVariable
	Value float32 `json:"value"`
	Max   int     `json:"max"`
	Min   int     `json:"min"`
	Step  int     `json:"step"`
}

type StringVariable struct {
	PortalVariable
	Value string `json:"value"`
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
	if err := encoder.Encode(variables.ToMap()); err != nil {
		panic(err)
	}
}

func (variables PortalVariables) UpdateVariables(varsMap VariablesMap) (PortalVariables, error) {
	for _, groups := range varsMap {
		for _, groupVars := range groups {
			for varName, variable := range groupVars {
				if _, ok := variables.Integer[varName]; ok {
					value := variable["value"].(int)
					if !ok {
						return PortalVariables{}, errors.New("value is not int")
					}

					currVar := variables.Integer[varName]
					currVar.Value = value
					variables.Integer[varName] = currVar
				}

				if _, ok := variables.Float[varName]; ok {
					value := variable["value"].(float32)
					if !ok {
						return PortalVariables{}, errors.New("value is not float32")
					}

					currVar := variables.Float[varName]
					currVar.Value = value
					variables.Float[varName] = currVar
				}

				if _, ok := variables.String[varName]; ok {
					value := variable["value"].(string)
					if !ok {
						return PortalVariables{}, errors.New("value is not string")
					}

					currVar := variables.String[varName]
					currVar.Value = value
					variables.String[varName] = currVar
				}
			}
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

type VariablesMap map[string]map[string]map[string]map[string]any

func (varsMap *VariablesMap) add(variable PortalVariable, value map[string]any) {
	if _, ok := (*varsMap)[variable.View]; !ok {
		(*varsMap)[variable.View] = make(map[string]map[string]map[string]any)
	}
	if _, ok := (*varsMap)[variable.View][variable.Group]; !ok {
		(*varsMap)[variable.View][variable.Group] = make(map[string]map[string]any)
	}

	(*varsMap)[variable.View][variable.Group][variable.Name] = value
}

// Converts PortalVariables struct to a hash map containing variables as final values and keys hierarchy: file -> group -> variable name
func (variables PortalVariables) ToMap() VariablesMap {
	mappedVariables := make(VariablesMap)

	for _, intVar := range variables.Integer {
		mappedVariables.add(intVar.PortalVariable, map[string]any{
			"displayName": intVar.DisplayName,
			"filePath":    intVar.FilePath,
			"value":       intVar.Value,
			"max":         intVar.Max,
			"min":         intVar.Min,
			"step":        intVar.Step,
			"type":        "integer",
		})
	}

	for _, floatVar := range variables.Float {
		mappedVariables.add(floatVar.PortalVariable, map[string]any{
			"displayName": floatVar.DisplayName,
			"filePath":    floatVar.FilePath,
			"value":       floatVar.Value,
			"max":         floatVar.Max,
			"min":         floatVar.Min,
			"step":        floatVar.Step,
			"type":        "float",
		})
	}

	for _, stringVar := range variables.String {
		mappedVariables.add(stringVar.PortalVariable, map[string]any{
			"displayName": stringVar.DisplayName,
			"filePath":    stringVar.FilePath,
			"value":       stringVar.Value,
			"type":        "string",
		})
	}

	return mappedVariables
}
