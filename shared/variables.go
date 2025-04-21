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

type FloatVariable struct {
	PortalVariable
	Value float32
	Max   int
	Min   int
	Step  int
}

type StringVariable struct {
	PortalVariable
	Value string
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
	var err error

	//TODO: Modularize following code with interfaces

	for key, value := range variablesPatch {
		if _, ok := variables.Integer[key]; ok {
			currentVar := variables.Integer[key]
			currentVar.Value, err = strconv.Atoi(value)
			if err != nil {
				return variables, errors.New("could not patch variables")
			}
			variables.Integer[key] = currentVar
		}

		if _, ok := variables.Float[key]; ok {
			currentVar := variables.Integer[key]
			currentVar.Value, err = strconv.Atoi(value)
			if err != nil {
				return variables, errors.New("could not patch variables")
			}
			variables.Integer[key] = currentVar
		}

		if _, ok := variables.String[key]; ok {
			currentVar := variables.String[key]
			currentVar.Value = value
			variables.String[key] = currentVar
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
