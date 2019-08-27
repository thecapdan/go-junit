// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package junit

import (
	"bytes"
	"encoding/xml"
	"io"
)

// reparentXML will wrap the given data (which is assumed to be valid XML), in
// a fake root nodeAlias.
//
// This action is useful in the event that the original XML document does not
// have a single root nodeAlias, which is required by the XML specification.
// Additionally, Go's XML parser will silently drop all nodes after the first
// that is encountered, which can lead to data loss from a parser perspective.
// This function also enables the ingestion of blank XML files, which would
// normally cause a parsing error.
func reparentXML(r io.Reader) io.Reader {
	var (
		pre  = bytes.NewReader([]byte("<fake-root>"))
		post = bytes.NewReader([]byte("</fake-root>"))
	)
	return io.MultiReader(pre, r, post)
}

// parse unmarshalls the given XML data into a graph of nodes, and then returns
// a slice of all top-level nodes.
func parse(r io.Reader) ([]xmlNode, error) {
	var (
		dec  = xml.NewDecoder(reparentXML(r))
		root xmlNode
	)

	if err := dec.Decode(&root); err != nil {
		return nil, err
	}

	return root.Nodes, nil
}
