package patcher

import (
	"portal/parser"
	"strings"
	"testing"
)

func TestPatcher(t *testing.T) {
	content := `let a = 2;

//@portal number
const b = 2;

//@portal string
let asd = "asdf";
`

	newVariables := parser.PortalVariables{
		Number: map[string]parser.NumberVariable{
			"b": {
				Name:  "b",
				Value: 100,
			},
		},
		String: map[string]parser.StringVariable{
			"asd": {
				Name:  "asd",
				Value: "qwerty",
			},
		},
	}

	newContent := strings.Split(PatchFile(content, newVariables), "\n")

	if newContent[3] != "const b = 100;" {
		t.Errorf("Wrong patched line %s", newContent[3])
	}

	if newContent[6] != "let asd = \"qwerty\";" {
		t.Errorf("Wrong patched line %s", newContent[6])
	}
}
