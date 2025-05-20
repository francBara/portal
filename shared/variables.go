package shared

import (
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

func (root1 UINode) isEqual(root2 UINode) bool {
	if root1.Type != root2.Type {
		return false
	}

	if len(root1.Properties) != len(root2.Properties) {
		return false
	}

	for i := range root1.Properties {
		if root1.Properties[i] != root2.Properties[i] {
			return false
		}
	}

	if len(root1.Children) != len(root2.Children) {
		return false
	}

	for i := range root1.Children {
		if !root1.Children[i].isEqual(*root2.Children[i]) {
			return false
		}
	}

	return true
}

type UIVariable struct {
	PortalVariable
	UINode
}

// FileVariables retains all annotated variables in a project, along with view, group and files data.
type FileVariables struct {
	Integer map[string]IntVariable    `json:"integer"`
	Float   map[string]FloatVariable  `json:"float"`
	String  map[string]StringVariable `json:"string"`
	UI      map[string]UIVariable     `json:"ui"`
}

type PortalVariables map[string]FileVariables

// Init allocates PortalVariables inner maps.
func (variables *FileVariables) Init() {
	variables.Integer = make(map[string]IntVariable)
	variables.Float = make(map[string]FloatVariable)
	variables.String = make(map[string]StringVariable)
	variables.UI = make(map[string]UIVariable)
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

// Collect returns FileVariables containing all variables across all files in PortalVariables.
func (variables PortalVariables) Collect() FileVariables {
	var mergedFileVars FileVariables

	for _, fileVariables := range variables {
		mergedFileVars = mergedFileVars.Merge(fileVariables)
	}

	return mergedFileVars
}

// GetPatch returns a new PortalVariables instance, where values are updated with VariablesMap values
func (variables PortalVariables) GetPatch(varsMap VariablesMap) (PortalVariables, error) {
	for _, groups := range varsMap {
		for _, groupVars := range groups {
			for varName, variable := range groupVars {
				fileVariables := variables[variable["filePath"].(string)]

				if _, ok := fileVariables.Integer[varName]; ok {
					value, ok := variable["value"].(int)
					if !ok {
						return PortalVariables{}, fmt.Errorf("variable %s is not int: %v %T", varName, variable["value"], variable["value"])
					}

					currVar := fileVariables.Integer[varName]

					if value == currVar.Value {
						continue
					}

					currVar.Value = value
					fileVariables.Integer[varName] = currVar
				} else if _, ok := fileVariables.Float[varName]; ok {
					value, ok := variable["value"].(float32)
					if !ok {
						return PortalVariables{}, errors.New("value is not float32")
					}

					currVar := fileVariables.Float[varName]

					if value == currVar.Value {
						continue
					}

					currVar.Value = value
					fileVariables.Float[varName] = currVar
				} else if _, ok := fileVariables.String[varName]; ok {
					value, ok := variable["value"].(string)
					if !ok {
						return PortalVariables{}, errors.New("value is not string")
					}

					currVar := fileVariables.String[varName]

					if value == currVar.Value {
						continue
					}

					currVar.Value = value
					fileVariables.String[varName] = currVar
				} else if _, ok := fileVariables.UI[varName]; ok {
					marshaled, err := json.Marshal(variable)
					if err != nil {
						panic(err)
					}

					var root UINode

					err = json.Unmarshal(marshaled, &root)
					if err != nil {
						return PortalVariables{}, errors.New("value is not UI node")
					}

					currVar := fileVariables.UI[varName]

					if currVar.UINode.isEqual(root) {
						continue
					}

					currVar.UINode = root
					fileVariables.UI[varName] = currVar
				}

				variables[variable["filePath"].(string)] = fileVariables
			}
		}
	}

	return variables, nil
}

func mergeMaps[K comparable, v any](map1 map[K]v, map2 map[K]v) map[K]v {
	newMap := make(map[K]v)

	maps.Copy(newMap, map1)
	maps.Copy(newMap, map2)

	return newMap
}

// Merge merges two FileVariables instances into a single one.
func (variables FileVariables) Merge(newVariables FileVariables) FileVariables {
	var merged FileVariables

	merged.Integer = mergeMaps(variables.Integer, newVariables.Integer)
	merged.Float = mergeMaps(variables.Float, newVariables.Float)
	merged.String = mergeMaps(variables.String, newVariables.String)
	merged.UI = mergeMaps(variables.UI, newVariables.UI)

	return merged
}

// Length returns the total number of variables in FileVariables.
func (fileVariables FileVariables) Length() int {
	return len(fileVariables.Integer) + len(fileVariables.Float) + len(fileVariables.String) + len(fileVariables.UI)
}

// Length returns the total number of variables in PortalVariables.
func (variables PortalVariables) Length() int {
	totalLength := 0

	for _, fileVariables := range variables {
		totalLength += fileVariables.Length()
	}

	return totalLength
}
