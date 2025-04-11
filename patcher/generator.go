package patcher

import (
	"bytes"
	"html/template"
	"log"
	"portal/parser"
)

type DashboardData struct {
	Components DashboarComponents
	UserName   string
}

func GenerateDashboard(variables parser.PortalVariables, userName string) string {
	var components DashboarComponents

	for _, numVar := range variables.Number {
		if numVar.Max == numVar.Min {
			components.NumberFields = append(components.NumberFields, NumberField{
				Id:           numVar.Name,
				Name:         numVar.Name,
				InitialValue: numVar.Value,
			})
		} else {
			components.Sliders = append(components.Sliders, Slider{
				Id:           numVar.Name,
				InitialValue: numVar.Value,
				Min:          numVar.Min,
				Max:          numVar.Max,
				Step:         numVar.Step,
				Name:         numVar.Name,
			})
		}
	}

	for _, stringVar := range variables.String {
		components.TextFields = append(components.TextFields, TextField{
			Id:           stringVar.Name,
			Name:         stringVar.Name,
			InitialValue: stringVar.Value,
		})
	}

	tmpl, err := template.ParseFiles("patcher/static/dashboard.html", "patcher/static/setters/slider.html", "patcher/static/setters/textField.html", "patcher/static/setters/numberField.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	data := DashboardData{
		Components: components,
		UserName:   userName,
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Fatal("Error executing template:", err)
	}

	return buf.String()
}
