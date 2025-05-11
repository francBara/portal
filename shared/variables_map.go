package shared

import (
	"encoding/json"
	"io"
)

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

// ToMap converts PortalVariables struct to a hash map containing variables as final values and keys hierarchy: file -> group -> variable name.
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

	for _, uiVar := range variables.UI {
		mappedVariables.add(uiVar.PortalVariable, uiVar.UINode.ToMap())
	}

	return mappedVariables
}

// JsonToVariablesMap takes a raw json and outputs a VariablesMap representation.
func JsonToVariablesMap(varsJson io.Reader) (VariablesMap, error) {
	var variablesMap *VariablesMap

	decoder := json.NewDecoder(varsJson)
	decoder.UseNumber()

	err := decoder.Decode(&variablesMap)
	if err != nil {
		return VariablesMap{}, err
	}

	for viewName, groups := range *variablesMap {
		for groupName, variables := range groups {
			for varName, variable := range variables {
				if num, ok := variable["value"].(json.Number); variable["type"] == "integer" && ok {
					newValue, err := num.Int64()
					if err != nil {
						return VariablesMap{}, err
					}
					(*variablesMap)[viewName][groupName][varName]["value"] = int(newValue)
				}
			}
		}
	}

	return *variablesMap, nil
}
