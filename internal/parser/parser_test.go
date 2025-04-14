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

	if len(variables.Number) != 1 {
		t.Errorf("Bad number variables length: %d", len(variables.Number))
	}
	if len(variables.String) != 1 {
		t.Errorf("Bad string variables length: %d", len(variables.String))
	}

	if variables.Number["maxChats"].Name != "maxChats" || variables.Number["maxChats"].Value != 24 {
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
		Number: map[string]shared.NumberVariable{
			"a": {
				PortalVariable: shared.PortalVariable{
					Name:     "a",
					FilePath: "simple_javascript.js",
					Group:    "gruppaccio",
				},
				Max:   1234,
				Min:   -100,
				Value: 2,
			},
		},
		String: map[string]shared.StringVariable{
			"b": {
				PortalVariable: shared.PortalVariable{
					Name:     "b",
					FilePath: "simple_javascript.js",
					Group:    "altro",
				},
				Value: "ciao",
			},
			"c": {
				PortalVariable: shared.PortalVariable{
					Name:     "c",
					FilePath: "simple_javascript.js",
					Group:    "Default",
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
