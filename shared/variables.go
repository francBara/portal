package shared

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
)

// PortalVariable retains common variables data.
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

// UINode is a tree representing an html tree, retaining tailwind properties.
type UINode struct {
	Type       string `json:"type"`
	Properties []struct {
		Prefix string `json:"prefix"`
		Value  string `json:"value"`
	} `json:"properties"`
	Children []*UINode `json:"children"`
}

func (node UINode) ToMap() map[string]any {
	data, err := json.Marshal(node)
	if err != nil {
		panic(err)
	}

	var result map[string]any

	err = json.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}
	return result
}

type UIRoot struct {
	PortalVariable
	UINode
}

// PortalVariables retains all annotated variables in a project, along with view, group and files data.
type PortalVariables struct {
	Integer    map[string]IntVariable    `json:"integer"`
	Float      map[string]FloatVariable  `json:"float"`
	String     map[string]StringVariable `json:"string"`
	UI         map[string]UIRoot         `json:"ui"`
	FileHashes map[string]string         `json:"fileHashes"`
}

// Init allocates PortalVariables inner maps.
func (variables *PortalVariables) Init() {
	variables.FileHashes = make(map[string]string)

	variables.Integer = make(map[string]IntVariable)
	variables.Float = make(map[string]FloatVariable)

	variables.String = make(map[string]StringVariable)

	variables.UI = make(map[string]UIRoot)
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

// UpdateVariables returns new PortalVariables with values updated by varsMap.
func (variables PortalVariables) UpdateVariables(varsMap VariablesMap) (PortalVariables, error) {
	for _, groups := range varsMap {
		for _, groupVars := range groups {
			for varName, variable := range groupVars {
				if _, ok := variables.Integer[varName]; ok {
					value, ok := variable["value"].(int)
					if !ok {
						return PortalVariables{}, fmt.Errorf("variable %s is not int: %v %T", varName, variable["value"], variable["value"])
					}

					currVar := variables.Integer[varName]
					currVar.Value = value
					variables.Integer[varName] = currVar
				}

				if _, ok := variables.Float[varName]; ok {
					value, ok := variable["value"].(float32)
					if !ok {
						return PortalVariables{}, errors.New("value is not float32")
					}

					currVar := variables.Float[varName]
					currVar.Value = value
					variables.Float[varName] = currVar
				}

				if _, ok := variables.String[varName]; ok {
					value, ok := variable["value"].(string)
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

// HasVariables returns true if PortalVariables is not empty.
func (variables PortalVariables) HasVariables() bool {
	return len(variables.Integer) > 0 || len(variables.String) > 0 || len(variables.Float) > 0
}

func mergeMaps[K comparable, v any](map1 map[K]v, map2 map[K]v) map[K]v {
	newMap := make(map[K]v)

	maps.Copy(newMap, map1)
	maps.Copy(newMap, map2)

	return newMap
}

// Merge merges two PortalVariables instances into a single one.
func (variables PortalVariables) Merge(newVariables PortalVariables) PortalVariables {
	var merged PortalVariables

	merged.Integer = mergeMaps(variables.Integer, newVariables.Integer)
	merged.Float = mergeMaps(variables.Float, newVariables.Float)
	merged.String = mergeMaps(variables.String, newVariables.String)
	merged.UI = mergeMaps(variables.UI, newVariables.UI)

	merged.FileHashes = mergeMaps(variables.FileHashes, newVariables.FileHashes)

	return merged
}

func (variables PortalVariables) HasFileChanged(fileContent string, filePath string) bool {
	hasher := sha256.New()
	hasher.Write([]byte(fileContent))
	hashString := hex.EncodeToString(hasher.Sum(nil))

	return hashString != variables.FileHashes[filePath]
}

// Length returns the total number of variables in PortalVariables.
func (variables PortalVariables) Length() int {
	return len(variables.Integer) + len(variables.Float) + len(variables.String) + len(variables.UI)
}
