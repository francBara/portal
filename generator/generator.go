package generator

import (
	"html/template"
	"log"
	"math/rand"
	"os"
	"portal/parser"
	"strconv"
)

type DashboarComponents struct {
	Sliders    []Slider
	TextFields []TextField
}

type DashboardData struct {
	Components DashboarComponents
	PatcherUrl string
}

func Generate(variables parser.PortalVariables) {
	var components DashboarComponents

	for _, numVar := range variables.Number {
		components.Sliders = append(components.Sliders, Slider{
			Id:           strconv.Itoa(rand.Intn(1000000)),
			InitialValue: numVar.Value,
			Min:          numVar.Min,
			Max:          numVar.Max,
			Step:         numVar.Step,
			Name:         numVar.Name,
		})
	}

	for _, stringVar := range variables.String {
		components.TextFields = append(components.TextFields, TextField{
			Id:           strconv.Itoa(rand.Intn(1000000)),
			Name:         stringVar.Name,
			InitialValue: stringVar.Value,
		})
	}

	tmpl, err := template.ParseFiles("generator/templates/dashboard.html", "generator/templates/setters/slider.html", "generator/templates/setters/textField.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	outFile, err := os.Create("out.html")
	if err != nil {
		log.Fatal("Error creating output file:", err)
	}
	defer outFile.Close()

	data := DashboardData{
		Components: components,
		PatcherUrl: "http://localhost:8080",
	}

	err = tmpl.Execute(outFile, data)
	if err != nil {
		log.Fatal("Error executing template:", err)
	}
}
