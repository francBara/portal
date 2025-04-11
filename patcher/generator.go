package patcher

import (
	"bytes"
	"html/template"
	"log"
	"math/rand"
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

func GenerateDashboard(variables parser.PortalVariables) string {
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

	tmpl, err := template.ParseFiles("patcher/static/dashboard.html", "patcher/static/setters/slider.html", "patcher/static/setters/textField.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	data := DashboardData{
		Components: components,
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Fatal("Error executing template:", err)
	}

	return buf.String()
}
