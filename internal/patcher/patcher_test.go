package patcher

import (
	"os"
	"portal/internal/parser"
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

func TestTailwindPatcher(t *testing.T) {
	content := `<div/ className={"
//@portal
border-round-2

//@portal
duration-[1000ms]
"}
`
	newVariables := shared.FileVariables{
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

	patched, err := PatchFile(content, newVariables)
	if err != nil {
		t.Error("Error patching file", err.Error())
	}

	newContent := strings.Split(patched, "\n")

	if newContent[2] != "border-round-100" {
		t.Errorf("Wrong patched line %s", newContent[2])
	}

	if newContent[5] != "duration-[247ms]" {
		t.Errorf("Wrong patched line %s", newContent[5])
	}
}

func TestUiPatcher(t *testing.T) {
	if err := os.Chdir("../.."); err != nil {
		panic(err)
	}

	variables, _, err := parser.ParseFile("internal/parser/tests", "ui.jsx", parser.ParseOptions{})
	if err != nil {
		panic(err)
	}

	fileContent, err := os.ReadFile("internal/parser/tests/ui.jsx")
	if err != nil {
		panic(err)
	}

	content := string(fileContent)

	variables.UI["CardLanding"].Children[0].Children[1].Properties[0] = struct {
		Prefix string "json:\"prefix\""
		Value  string "json:\"value\""
	}{
		Prefix: "cursor",
		Value:  "puntatore",
	}

	patched, err := PatchFile(content, variables)
	if err != nil {
		panic(err)
	}

	if !strings.Contains(patched, "className=\"cursor-puntatore px-5 pb-5\"") {
		t.Error("bad patch")
	}
}
