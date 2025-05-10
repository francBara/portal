package preview

import (
	"os"
	"os/exec"
)

func generatePreview(filePath string) error {
	cmd := exec.Command("node", "tools/generatePreview.js", filePath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
