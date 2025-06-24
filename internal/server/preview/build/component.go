package build

import (
	"fmt"
	"os"
	"path/filepath"
	"portal/shared"
)

func makeEntryPoint(variable shared.UIVariable, componentFilePath string) error {
	relPath, err := filepath.Rel("component-preview/src", filepath.Join("component-preview/src/components", componentFilePath))
	if err != nil {
		return err
	}

	variableDeclarations := ""
	componentProps := ""

	for name, value := range variable.PropsMocks {
		variableDeclarations += fmt.Sprintf("const %s = %s;\n", name, value)
		componentProps += fmt.Sprintf("%s={%s} ", name, name)
	}

	boxString := ""

	if variable.Box.Height != 0 {
		boxString += fmt.Sprintf("h-%d ", variable.Box.Height)
	} else {
		boxString += "h-full "
	}
	if variable.Box.Width != 0 {
		boxString += fmt.Sprintf("w-%d", variable.Box.Width)
	} else {
		boxString += "w-full"
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
				<div className="%s">
					<%s %s/>
				</div>
			</div>
		</BrowserRouter>
	</React.StrictMode>
);
`, variable.Name, relPath, variableDeclarations, boxString, variable.Name, componentProps)

	if err = os.WriteFile("component-preview/src/index.jsx", []byte(fileContent), os.ModePerm); err != nil {
		return err
	}

	cssPath := seekFiles([]string{"index.css", "src/index.css"})

	if cssPath == "" {
		return os.WriteFile("component-preview/src/index.css", []byte("@tailwind base;\n@tailwind components;\n@tailwind utilities;\n"), os.ModePerm)
	}

	if err = copyFile(cssPath, "component-preview/src/index.css"); err != nil {
		return err
	}

	indexPath := seekFiles([]string{"index.html", "src/index.html"})

	return copyFile(indexPath, "component-preview/index.html")
}
