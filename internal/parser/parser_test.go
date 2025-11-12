package parser

import (
	"fmt"
	"portal/shared"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSubfolders(t *testing.T) {
	variables, err := ParseProject("tests/togiftit", ParseOptions{Verbose: true})
	if err != nil {
		t.Errorf("Error parsing project")
	}

	for k := range variables {
		fmt.Println(k)
	}

	fileVars := variables["chats/Chat.tsx"]

	if len(fileVars.Integer) != 1 {
		t.Errorf("Bad integer fileVars length: %d", len(fileVars.Integer))
	}
	if len(fileVars.String) != 1 {
		t.Errorf("Bad string fileVars length: %d", len(fileVars.String))
	}

	if fileVars.Integer["maxChats"].Name != "maxChats" || fileVars.Integer["maxChats"].Value != 24 {
		t.Error("Bad number variable")
	}
	if fileVars.String["chatName"].Name != "chatName" || fileVars.String["chatName"].Value != "My Chat" {
		t.Error("Bad string variable")
	}
}

func TestJavascript(t *testing.T) {
	variables, err := ParseFile("tests", "simple_javascript.js", ParseOptions{Verbose: true})
	if err != nil {
		t.Errorf("error parsing project: " + err.Error())
	}

	expected := shared.FileVariables{
		Integer: map[string]shared.IntVariable{
			"a:2": {
				PortalVariable: shared.PortalVariable{
					Name:        "a",
					DisplayName: "a",
					View:        "simple_javascript",
					Group:       "gruppaccio",
					LineNumber:  2,
				},
				Max:   1234,
				Min:   -100,
				Value: 2,
			},
		},
		Float: map[string]shared.FloatVariable{},
		String: map[string]shared.StringVariable{
			"b:5": {
				PortalVariable: shared.PortalVariable{
					Name:        "b",
					DisplayName: "ecco il bel nome",
					View:        "simple_javascript",
					Group:       "altro",
					LineNumber:  5,
				},
				Value: "ciao",
			},
			"c:9": {
				PortalVariable: shared.PortalVariable{
					Name:        "c",
					DisplayName: "nome",
					View:        "simple_javascript",
					Group:       "Default",
					LineNumber:  9,
				},
				Value: "eccoci",
			},
		},
	}

	if diff := cmp.Diff(expected, variables); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestJavascriptAll(t *testing.T) {
	variables, err := ParseFile("tests", "simple_javascript_all.js", ParseOptions{Verbose: true})
	if err != nil {
		t.Errorf("error parsing project: " + err.Error())
	}

	if variables.Length() != 3 {
		t.Errorf("Wrong number of variables %d, expected 3", variables.Length())
	}

	if _, ok := variables.Integer["d:14"]; !ok {
		t.Error("int variable d not present")
	}
	if _, ok := variables.String["b:8"]; !ok {
		t.Error("string variable b not present")
	}
	if _, ok := variables.String["c:11"]; !ok {
		t.Error("string variable c not present")
	}

	if variables.Integer["d:14"].Value != 24 {
		t.Errorf("int variable d wrong value %d", variables.Integer["d:14"].Value)
	}
}
