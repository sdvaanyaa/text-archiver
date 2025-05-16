package vlc

import "strings"

type DecodingTree struct {
	Value string
	Zero  *DecodingTree
	One   *DecodingTree
}

// DecodingTree builds a decoding tree based on the encoding table.
// Each rune is associated with a binary string (code), which defines a path in the tree.
// At the end of each path, the corresponding value is stored in the node.
func (et encodingTable) DecodingTree() DecodingTree {
	res := DecodingTree{}

	for ch, code := range et {
		res.Add(code, ch)
	}

	return res
}

func (dt *DecodingTree) Decode(str string) string {
	var b strings.Builder

	currNode := dt

	for _, ch := range str {
		switch ch {
		case '0':
			currNode = currNode.Zero
		case '1':
			currNode = currNode.One
		}

		if currNode.Value != "" {
			b.WriteString(currNode.Value)
			currNode = dt
		}
	}

	return b.String()
}

// Add inserts a value into the decoding tree by following the binary code string.
// For each '0' or '1' in the code, it descends to the corresponding child node.
// If a node does not exist, it creates a new one. The value is stored in the final node.
// 01000(0) <- value
func (dt *DecodingTree) Add(code string, value rune) {
	currNode := dt

	for _, ch := range code {
		switch ch {
		case '0':
			if currNode.Zero == nil {
				currNode.Zero = &DecodingTree{}
			}

			currNode = currNode.Zero
		case '1':
			if currNode.One == nil {
				currNode.One = &DecodingTree{}
			}

			currNode = currNode.One
		}
	}

	currNode.Value = string(value)
}
