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
	Name     string
	Group    string
	FilePath string
}

type NumberVariable struct {
	PortalVariable
	Value int
	Max   int
	Min   int
	Step  int
}

type StringVariable struct {
	PortalVariable
	Value string
}

type PortalVariables struct {
	Number     map[string]NumberVariable
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

	for key, value := range variablesPatch {
		if _, ok := variables.Number[key]; ok {
			currentVar := variables.Number[key]
			currentVar.Value, err = strconv.Atoi(value)
			if err != nil {
				return variables, errors.New("could not patch variables")
			}
			variables.Number[key] = currentVar
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
	return len(variables.Number) > 0 || len(variables.String) > 0
}

func mergeMaps[K comparable, v any](map1 map[K]v, map2 map[K]v) map[K]v {
	newMap := make(map[K]v)

	maps.Copy(newMap, map1)
	maps.Copy(newMap, map2)

	return newMap
}

func (variables PortalVariables) Merge(newVariables PortalVariables) PortalVariables {
	var merged PortalVariables

	merged.Number = mergeMaps(variables.Number, newVariables.Number)
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
