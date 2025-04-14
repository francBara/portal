package patcher

import (
	"portal/shared"
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

	newVariables := shared.PortalVariables{
		Number: map[string]shared.NumberVariable{
			"b": {
				PortalVariable: shared.PortalVariable{
					Name: "b",
				},
				Value: 100,
			},
		},
		String: map[string]shared.StringVariable{
			"asd": {
				PortalVariable: shared.PortalVariable{
					Name: "asd",
				},
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
