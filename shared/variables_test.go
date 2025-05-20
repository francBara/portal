package shared

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestVariablesMap(t *testing.T) {
	variables := PortalVariables{
		"index.js": FileVariables{
			Integer: map[string]IntVariable{
				"numero": {
					PortalVariable: PortalVariable{
						Name:        "numero",
						Group:       "gruppettino",
						DisplayName: "numerino",
						View:        "index",
					},
					Max:   100,
					Min:   0,
					Value: 2,
				},
			},
			String: map[string]StringVariable{
				"testo": {
					PortalVariable: PortalVariable{
						Name:        "testo",
						Group:       "gruppettino",
						DisplayName: "testino",
						View:        "index",
					},
					Value: "contenuto",
				},
			},
		},
	}

	expected := VariablesMap{
		"index": {
			"gruppettino": {
				"numero": map[string]any{
					"displayName": "numerino",
					"max":         100,
					"min":         0,
					"step":        0,
					"value":       2,
					"type":        "integer",
				},
				"testo": map[string]any{
					"displayName": "testino",
					"value":       "contenuto",
					"type":        "string",
				},
			},
		},
	}

	mapped := variables.ToMap()

	if diff := cmp.Diff(expected, mapped); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestUINode(t *testing.T) {
	root1 := UINode{
		Type: "div",
		Properties: []struct {
			Prefix string "json:\"prefix\""
			Value  string "json:\"value\""
		}{
			{
				Prefix: "mt",
				Value:  "4",
			},
		},
		Children: []*UINode{
			{
				Type: "arg",
				Properties: []struct {
					Prefix string "json:\"prefix\""
					Value  string "json:\"value\""
				}{
					{
						Prefix: "mb",
						Value:  "5",
					},
				},
			},
		},
	}

	root2 := UINode{
		Type: "div",
		Properties: []struct {
			Prefix string "json:\"prefix\""
			Value  string "json:\"value\""
		}{
			{
				Prefix: "mt",
				Value:  "4",
			},
		},
		Children: []*UINode{
			{
				Type: "arg",
				Properties: []struct {
					Prefix string "json:\"prefix\""
					Value  string "json:\"value\""
				}{
					{
						Prefix: "mb",
						Value:  "5",
					},
				},
			},
		},
	}

	if !root1.isEqual(root2) {
		t.Error("nodes not equal")
	}

	root2.Children[0].Type = "usd"

	if root1.isEqual(root2) {
		t.Error("nodes equal")
	}
}
