package parser

import (
	"os"
	"portal/shared"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSubfolders(t *testing.T) {
	variables, err := ParseProject("tests/togiftit", ParseOptions{Verbose: false})
	if err != nil {
		t.Errorf("Error parsing project")
	}

	if len(variables.Integer) != 1 {
		t.Errorf("Bad integer variables length: %d", len(variables.Integer))
	}
	if len(variables.String) != 1 {
		t.Errorf("Bad string variables length: %d", len(variables.String))
	}

	if variables.Integer["maxChats"].Name != "maxChats" || variables.Integer["maxChats"].Value != 24 {
		t.Error("Bad number variable")
	}
	if variables.String["chatName"].Name != "chatName" || variables.String["chatName"].Value != "My Chat" {
		t.Error("Bad string variable")
	}
}

func TestJavascript(t *testing.T) {
	variables, err := parseFile("tests", "simple_javascript.js", ParseOptions{Verbose: false})
	if err != nil {
		t.Errorf("error parsing project: " + err.Error())
	}

	file, err := os.Open("tests/simple_javascript.js")
	if err != nil {
		panic(err)
	}
	fileHash := getFileHash(file)
	file.Close()

	expected := shared.PortalVariables{
		Integer: map[string]shared.IntVariable{
			"a": {
				PortalVariable: shared.PortalVariable{
					Name:        "a",
					DisplayName: "a",
					View:        "simple_javascript.js",
					FilePath:    "simple_javascript.js",
					Group:       "gruppaccio",
				},
				Max:   1234,
				Min:   -100,
				Value: 2,
			},
		},
		Float: map[string]shared.FloatVariable{},
		String: map[string]shared.StringVariable{
			"b": {
				PortalVariable: shared.PortalVariable{
					Name:        "b",
					DisplayName: "ecco il bel nome",
					View:        "simple_javascript.js",
					FilePath:    "simple_javascript.js",
					Group:       "altro",
				},
				Value: "ciao",
			},
			"c": {
				PortalVariable: shared.PortalVariable{
					Name:        "c",
					DisplayName: "nome",
					View:        "simple_javascript.js",
					FilePath:    "simple_javascript.js",
					Group:       "Default",
				},
				Value: "eccoci",
			},
		},
		FileHashes: map[string]string{
			"simple_javascript.js": fileHash,
		},
	}

	if diff := cmp.Diff(expected, variables); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestTailwind(t *testing.T) {
	variables, err := parseFile("tests", "tailwind.js", ParseOptions{Verbose: false})
	if err != nil {
		t.Errorf("error parsing project: " + err.Error())
	}

	if variables.Integer["duration"].Value != 1000 {
		t.Errorf("bad duration value %d", variables.Integer["duration"].Value)
	}

	if variables.Integer["hover:scale"].Value != 92 {
		t.Errorf("bad hover:scale value %d", variables.Integer["hover:scale"].Value)
	}

	if variables.Integer["bg-color-red"].Value != 500 {
		t.Errorf("bad bg-color value %d", variables.Integer["bg-color-red"].Value)
	}
}
