package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func ExecuteTool(tool string, input any) (bytes.Buffer, error) {
	var inputBuff bytes.Buffer
	err := json.NewEncoder(&inputBuff).Encode(input)
	if err != nil {
		return bytes.Buffer{}, err
	}

	cmd := exec.Command("node", fmt.Sprintf("tools/%s.js", tool))

	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Stdin = &inputBuff

	err = cmd.Run()
	if err != nil {
		return bytes.Buffer{}, err
	}

	return out, nil
}
