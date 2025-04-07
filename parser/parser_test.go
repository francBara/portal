package parser

import "testing"

func TestParser(t *testing.T) {
	variables := ParseProject("tests/togiftit", ParseOptions{Verbose: false})

	if len(variables.Number) != 1 {
		t.Errorf("Bad number variables length: %d", len(variables.Number))
	}
	if len(variables.String) != 1 {
		t.Errorf("Bad string variables lenght: %d", len(variables.String))
	}

	if variables.Number[0].Name != "maxChats" || variables.Number[0].Value != 24 {
		t.Error("Bad number variable")
	}
	if variables.String[0].Name != "chatName" || variables.String[0].Value != "My Chat" {
		t.Error("Bad string variable")
	}
}
