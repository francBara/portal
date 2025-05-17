package build

import (
	"bufio"
	"fmt"
	"os"
	"portal/internal/server/github"
	"strings"
)

type envFile map[string]string

func parseEnvFile(path string) (envFile, error) {
	file, err := os.Open(fmt.Sprintf("%s/%s", github.RepoFolderName, path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	envMap := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

		envMap[key] = fmt.Sprintf("\"%s\"", value)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return envMap, nil
}

func (env envFile) toReactVite() envFile {
	newEnv := make(envFile)

	for k, v := range env {
		newEnv[fmt.Sprintf("import.meta.env.%s", k)] = v
	}

	return newEnv
}
