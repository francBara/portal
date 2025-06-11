package build

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"portal/internal/server/github"
	"portal/shared"
)

type componentMock struct {
	ComponentName string         `json:"componentName"`
	Mock          map[string]any `json:"mock"`
}

func scanComponent(componentFilePath string) (mock componentMock, err error) {
	file, err := os.ReadFile(filepath.Join(github.RepoFolderName, componentFilePath))
	if err != nil {
		return componentMock{}, err
	}

	out, err := shared.ExecuteTool("scanComponentPreview", map[string]any{
		"sourceCode": string(file),
	})
	if err != nil {
		return componentMock{}, err
	}

	var result componentMock

	if err = json.NewDecoder(&out).Decode(&result); err != nil {
		return componentMock{}, err
	}

	if result.ComponentName == "" {
		return componentMock{}, errors.New("no portal component found")
	}

	return result, nil
}

func makeEntryPoint(component componentMock, componentFilePath string) error {
	relPath, err := filepath.Rel("component-preview/src", filepath.Join("component-preview/src/components", componentFilePath))
	if err != nil {
		return err
	}

	variableDeclarations := ""
	componentProps := ""

	for name, value := range component.Mock {
		marshaledValue, err := json.Marshal(value)
		if err != nil {
			return err
		}

		variableDeclarations += fmt.Sprintf("const %s = %s;\n", name, string(marshaledValue))
		componentProps += fmt.Sprintf("%s={%s} ", name, name)
	}

	fileContent := fmt.Sprintf(`import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import "./index.css";
import %s from './%s';

%s
const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
	<React.StrictMode>
		<BrowserRouter>
			<div className="min-h-screen flex items-center justify-center">
				<div className="h-200 w-64">
					<%s %s/>
				</div>
			</div>
		</BrowserRouter>
	</React.StrictMode>
);
`, component.ComponentName, relPath, variableDeclarations, component.ComponentName, componentProps)

	if err = os.WriteFile("component-preview/src/index.jsx", []byte(fileContent), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile("component-preview/src/index.css", []byte("@tailwind base;\n@tailwind components;\n@tailwind utilities;\n"), os.ModePerm)

}
