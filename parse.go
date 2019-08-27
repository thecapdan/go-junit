// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package junit

import (
	"encoding/xml"
	"io"
)

// parse unmarshalls the given XML data into a graph of nodes, and then returns
// a slice of all top-level nodes.
func parse(r io.Reader) ([]xmlNode, error) {
	var (
		dec  = xml.NewDecoder(r)
		root xmlNode
	)

	if err := dec.Decode(&root); err != nil {
		return nil, err
	}

	return root.Nodes, nil
}
