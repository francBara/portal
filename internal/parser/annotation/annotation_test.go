package annotation

import (
	"fmt"
	"testing"
)

func TestTokenizer(t *testing.T) {
	expected := []string{"group", "=", "haha", "asd", "bsdc", "\"no wayyy === 21oin { {\"", "{asd: 2, \"bs gjk\"}", "23", "2", "{{{{{{hey: [\"{{ hahah   \"]}}}}}}"}

	annotation := ""

	for _, e := range expected {
		annotation += e
		annotation += "   "
	}

	tokens := tokenizeAnnotation(annotation)

	fmt.Println(tokens)

	for idx := range expected {
		if tokens[idx] != expected[idx] {
			t.Errorf("expected %s, got %s", expected[idx], tokens[idx])
		}
	}
}

func TestParser(t *testing.T) {
	tokens := []string{"all", "ui", "group", "=", "blabla", "view", "=", "sample", "mock", "\"Rocky Balboa\"", "2", "{asd: 2}"}

	annotation := parseTokens(tokens)

	if !annotation.All {
		t.Error("all is false")
	}
	if !annotation.UI {
		t.Error("UI is false")
	}

	if annotation.Group != "blabla" {
		t.Errorf("bad group, expected %s, got %s", "blabla", annotation.Group)
	}

	if annotation.View != "sample" {
		t.Errorf("bad view, expected %s, got %s", "sample", annotation.View)
	}

	if len(annotation.Mock) != 3 {
		t.Errorf("bad mock length, expected %d, got %d", 3, len(annotation.Mock))
	}

	for i, mock := range annotation.Mock {
		if mock != tokens[i+9] {
			t.Errorf("bad mock, expected %s, got %s", tokens[i+9], mock)
		}
	}
}

func TestAnnotation(t *testing.T) {
	annotationStr := "ui    view     =   lskjdf		name = \"lorem ipsum hahaha\" mock all user gift 2 "

	annotation := ParseAnnotation(annotationStr)

	if !annotation.UI {
		t.Error("UI is false")
	}
	if annotation.All {
		t.Error("All is true")
	}

	if annotation.Name != "lorem ipsum hahaha" {
		t.Errorf("bad name, expected lorem ipsum hahaha, got %s", annotation.Name)
	}

	if annotation.Group != "" {
		t.Error("group is not empty")
	}

	if annotation.View != "lskjdf" {
		t.Errorf("bad view, expected %s, got %s", "lskjdf", annotation.View)
	}

	if annotation.Mock[0] != "all" {
		t.Errorf("bad first mock")
	}
	if annotation.Mock[1] != "user" {
		t.Error("bad second mock")
	}
	if annotation.Mock[2] != "gift" {
		t.Error("bad third mock")
	}
	if annotation.Mock[3] != "2" {
		t.Error("bad fourth mock")
	}
}
