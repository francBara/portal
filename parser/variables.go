package parser

import (
	"encoding/json"
	"os"
)

type NumberVariable struct {
	Name  string
	Value int
	Max   int
	Min   int
	Step  int
}

type StringVariable struct {
	Name  string
	Value string
}

type PortalVariables struct {
	Number []NumberVariable
	String []StringVariable
}

func (pv PortalVariables) Concat(newPv PortalVariables) PortalVariables {
	var concatenatedPv PortalVariables
	concatenatedPv.Number = append(pv.Number, newPv.Number...)
	concatenatedPv.String = append(pv.String, newPv.String...)
	return concatenatedPv
}

func (variables PortalVariables) Dump() {
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
