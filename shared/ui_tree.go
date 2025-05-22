package shared

import "encoding/json"

// UINode is a tree representing an html tree, retaining tailwind properties.
type UINode struct {
	Type       string `json:"type"`
	Id         int    `json:"id"`
	Properties []struct {
		Prefix string `json:"prefix"`
		Value  string `json:"value"`
	} `json:"properties"`
	Children []*UINode `json:"children"`
}

func (node UINode) ToMap() map[string]any {
	data, err := json.Marshal(node)
	if err != nil {
		panic(err)
	}

	var result map[string]any

	err = json.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func (root1 UINode) isEqual(root2 UINode) bool {
	if root1.Type != root2.Type {
		return false
	}

	if len(root1.Properties) != len(root2.Properties) {
		return false
	}

	for i := range root1.Properties {
		if root1.Properties[i] != root2.Properties[i] {
			return false
		}
	}

	if len(root1.Children) != len(root2.Children) {
		return false
	}

	for i := range root1.Children {
		if !root1.Children[i].isEqual(*root2.Children[i]) {
			return false
		}
	}

	return true
}

// GenerateIDs generates a unique incremental ID for every node of the tree
func (root *UINode) GenerateIDs(rootId *int) {
	root.Id = *rootId

	*rootId++

	for _, child := range root.Children {
		child.GenerateIDs(rootId)
	}
}
