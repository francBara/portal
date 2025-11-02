package patcher

import (
	"regexp"
	"strconv"
	"strings"
)

func UpdateTailwindLine(line string, newValue int) string {
	valueIdx := strings.LastIndex(line, "-")

	newValueStr := regexp.MustCompile(`\d+`).ReplaceAllString(line[valueIdx:], strconv.Itoa(newValue))

	return line[:valueIdx] + newValueStr
}
