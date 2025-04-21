package generator

import (
	"bytes"
	"html/template"
	"io/fs"
	"log"
	"path/filepath"
	"portal/shared"
)

type dashboardData struct {
	VariablesGroups []DashboardGroup
	UserName        string
}

func GenerateDashboard(variables shared.PortalVariables, userName string) string {
	dashboardGroups := make(dashboardComponents)

	for _, numVar := range variables.Integer {
		currentFields := dashboardGroups[numVar.Group]

		if numVar.Max == numVar.Min {
			currentFields.NumberFields = append(currentFields.NumberFields, NumberField{
				Id:           numVar.Name,
				Name:         numVar.DisplayName,
				InitialValue: numVar.Value,
			})
		} else {
			currentFields.Sliders = append(currentFields.Sliders, Slider{
				Id:           numVar.Name,
				InitialValue: numVar.Value,
				Min:          numVar.Min,
				Max:          numVar.Max,
				Step:         numVar.Step,
				Name:         numVar.DisplayName,
			})
		}

		dashboardGroups[numVar.Group] = currentFields
	}

	//TODO: Add float generator

	for _, stringVar := range variables.String {
		currentFields := dashboardGroups[stringVar.Group]

		currentFields.TextFields = append(currentFields.TextFields, TextField{
			Id:           stringVar.Name,
			Name:         stringVar.Name,
			InitialValue: stringVar.Value,
		})

		dashboardGroups[stringVar.Group] = currentFields
	}

	tmpl, err := template.ParseFiles(getTemplates()...)
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	data := dashboardData{
		VariablesGroups: dashboardGroups.toSlice(),
		UserName:        userName,
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Fatal("Error executing template:", err)
	}

	return buf.String()
}

func getTemplates() []string {
	var templates []string

	filepath.WalkDir("static/dashboard", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".html" {
			templates = append(templates, path)
		}
		return nil
	})

	return templates
}
