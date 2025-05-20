package shared

import (
	"encoding/json"
	"io"
)

type VariablesMap map[string]map[string]map[string]map[string]any

func (varsMap *VariablesMap) add(variable PortalVariable, filePath string, value map[string]any) {
	if _, ok := (*varsMap)[variable.View]; !ok {
		(*varsMap)[variable.View] = make(map[string]map[string]map[string]any)
	}
	if _, ok := (*varsMap)[variable.View][variable.Group]; !ok {
		(*varsMap)[variable.View][variable.Group] = make(map[string]map[string]any)
	}

	value["displayName"] = variable.DisplayName
	value["filePath"] = filePath
	(*varsMap)[variable.View][variable.Group][variable.Name] = value

}

// ToMap converts PortalVariables struct to a hash map containing variables as final values and keys hierarchy: file -> group -> variable name.
func (variables PortalVariables) ToMap() VariablesMap {
	mappedVariables := make(VariablesMap)

	for filePath, fileVars := range variables {
		for _, intVar := range fileVars.Integer {
			mappedVariables.add(intVar.PortalVariable, filePath, map[string]any{
				"value": intVar.Value,
				"max":   intVar.Max,
				"min":   intVar.Min,
				"step":  intVar.Step,
				"type":  "integer",
			})
		}

		for _, floatVar := range fileVars.Float {
			mappedVariables.add(floatVar.PortalVariable, filePath, map[string]any{
				"value": floatVar.Value,
				"max":   floatVar.Max,
				"min":   floatVar.Min,
				"step":  floatVar.Step,
				"type":  "float",
			})
		}

		for _, stringVar := range fileVars.String {
			mappedVariables.add(stringVar.PortalVariable, filePath, map[string]any{
				"value": stringVar.Value,
				"type":  "string",
			})
		}

		for _, uiVar := range fileVars.UI {
			mappedVariables.add(uiVar.PortalVariable, filePath, uiVar.UINode.ToMap())
		}
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
