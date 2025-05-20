package patcher

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"portal/shared"
	"regexp"
	"strconv"
	"strings"
)

func UpdateTailwindLine(line string, newValue int) string {
	valueIdx := strings.LastIndex(line, "-")

	newValueStr := regexp.MustCompile(`\d+`).ReplaceAllString(line[valueIdx:], strconv.Itoa(newValue))

	return line[:valueIdx] + newValueStr
}

// patchUI applies UI variables to a file, using updateTree.js tool.
func patchUI(content string, roots map[string]shared.UIVariable) (newContent string, err error) {
	if len(roots) == 0 {
		return content, nil
	}

	// updateTree.js accepts source code and variables via stdin

	payload := map[string]any{
		"sourceCode": content,
		"components": roots,
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("node", "tools/updateTree.js")

	cmd.Env = append(os.Environ(), "NODE_PATH=/usr/local/lib/node_modules")

	var stdoutBuf bytes.Buffer

	cmd.Stdout = &stdoutBuf
	cmd.Stderr = os.Stderr
	cmd.Stdin = &buf

	err = cmd.Run()

	return string(stdoutBuf.String()), err
}
