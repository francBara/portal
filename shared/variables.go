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
	Name        string
	DisplayName string
	Group       string
	FilePath    string
}

type IntVariable struct {
	PortalVariable
	Value int
	Max   int
	Min   int
	Step  int
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
	Value float32
	Max   int
	Min   int
	Step  int
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
	Value string
}

func (variable StringVariable) update(value string) StringVariable {
	variable.Value = value
	return variable
}

type PortalVariables struct {
	Integer    map[string]IntVariable
	Float      map[string]FloatVariable
	String     map[string]StringVariable
	FileHashes map[string]string
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
	//TODO: Modularize following code with interfaces

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
