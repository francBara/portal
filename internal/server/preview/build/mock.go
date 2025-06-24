package build

import (
	"os"
	"portal/shared"
)

func applyMocks(srcPath string, destPath string, mocks shared.FileMocks) error {
	rawContent, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	out, err := shared.ExecuteTool("mockComponent", map[string]any{
		"sourceCode": string(rawContent),
		"mocks":      mocks,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(destPath, out.Bytes(), os.ModePerm)
}
