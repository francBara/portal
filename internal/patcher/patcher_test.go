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

	newVariables := shared.FileVariables{
		Integer: map[string]shared.IntVariable{
			"b:4": {
				PortalVariable: shared.PortalVariable{
					Name:       "b",
					LineNumber: 4,
				},
				Value: 100,
			},
		},
		String: map[string]shared.StringVariable{
			"asd:7": {
				PortalVariable: shared.PortalVariable{
					Name:       "asd",
					LineNumber: 7,
				},
				Value: "qwerty",
			},
		},
	}

	patched, err := PatchFile(content, newVariables)
	if err != nil {
		t.Error("Error patching file", err.Error())
	}

	newContent := strings.Split(patched, "\n")

	if newContent[3] != "const b = 100;" {
		t.Errorf("Wrong patched line %s", newContent[3])
	}

	if newContent[6] != "let asd = \"qwerty\";" {
		t.Errorf("Wrong patched line %s", newContent[6])
	}
}
