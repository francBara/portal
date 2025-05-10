package shared

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestVariablesMap(t *testing.T) {
	variables := PortalVariables{
		Integer: map[string]IntVariable{
			"numero": {
				PortalVariable: PortalVariable{
					Name:        "numero",
					Group:       "gruppettino",
					DisplayName: "numerino",
					FilePath:    "index.js",
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
					FilePath:    "index.js",
					View:        "index",
				},
				Value: "contenuto",
			},
		},
	}

	expected := VariablesMap{
		"index": {
			"gruppettino": {
				"numero": map[string]any{
					"displayName": "numerino",
					"filePath":    "index.js",
					"max":         100,
					"min":         0,
					"step":        0,
					"value":       2,
					"type":        "integer",
				},
				"testo": map[string]any{
					"displayName": "testino",
					"filePath":    "index.js",
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
