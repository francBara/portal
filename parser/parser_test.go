package parser

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	variables := ParseProject("tests/togiftit", ParseOptions{Verbose: false})

	fmt.Println(variables)

	if len(variables.Number) != 1 {
		t.Errorf("Bad number variables length: %d", len(variables.Number))
	}
	if len(variables.String) != 1 {
		t.Errorf("Bad string variables lenght: %d", len(variables.String))
	}

	if variables.Number["maxChats"].Name != "maxChats" || variables.Number["maxChats"].Value != 24 {
		t.Error("Bad number variable")
	}
	if variables.String["chatName"].Name != "chatName" || variables.String["chatName"].Value != "My Chat" {
		t.Error("Bad string variable")
	}
}
