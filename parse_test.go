// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package junit

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		title       string
		input       []byte
		expected    []xmlNode
		expectedErr error
	}{
		{
			title:       "nil input",
			expectedErr: io.EOF,
		},
		{
			title:       "empty input",
			input:       []byte(``),
			expectedErr: io.EOF,
		},
		{
			title:       "plaintext input",
			input:       []byte(`This is some data that does not look like xml.`),
			expectedErr: io.EOF,
		},
		{
			title:       "json input",
			input:       []byte(`{"This is some data": "that looks like json"}`),
			expectedErr: io.EOF,
		},
		{
			title: "single xml node",
			input: []byte(`<t><this-is-a-tag/></t>`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
				},
			},
		},
		{
			title: "multiple xml nodes",
			input: []byte(`
				<t>
				<this-is-a-tag/>
				<this-is-also-a-tag/>
				</t>
			`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
				},
				{
					XMLName: xml.Name{
						Local: "this-is-also-a-tag",
					},
				},
			},
		},
		{
			title: "single xml node with content",
			input: []byte(`<t><this-is-a-tag>This is some content.</this-is-a-tag></t>`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
					Content: []byte("This is some content."),
				},
			},
		},
		{
			title: "single xml node with attributes",
			input: []byte(`<t><this-is-a-tag name="my name" status="passed"></this-is-a-tag></t>`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
					Attrs: map[string]string{
						"name":   "my name",
						"status": "passed",
					},
				},
			},
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("#%d - %s", index+1, test.title)

		t.Run(name, func(t *testing.T) {
			actual, err := parse(bytes.NewReader(test.input))

			if test.expectedErr != nil {
				require.Equal(t, test.expectedErr, err)
			} else {
				require.Nil(t, err)
			}

			assert.Equal(t, test.expected, actual)
		})
	}
}
