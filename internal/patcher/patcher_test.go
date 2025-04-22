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
		Integer: map[string]shared.IntVariable{
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

func TestTailwindPatcher(t *testing.T) {
	content := `<div/ className={"
//@portal
border-round-2

//@portal
duration-[1000ms]
"}
`
	newVariables := shared.PortalVariables{
		Integer: map[string]shared.IntVariable{
			"border-round": {
				PortalVariable: shared.PortalVariable{
					Name: "border-round",
				},
				Value: 100,
			},
			"duration": {
				PortalVariable: shared.PortalVariable{
					Name: "duration",
				},
				Value: 247,
			},
		},
	}

	newContent := strings.Split(PatchFile(content, newVariables), "\n")

	if newContent[2] != "border-round-100" {
		t.Errorf("Wrong patched line %s", newContent[2])
	}

	if newContent[5] != "duration-[247ms]" {
		t.Errorf("Wrong patched line %s", newContent[5])
	}
}
