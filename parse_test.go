// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package junit

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		title    string
		input    []byte
		expected []xmlNode
	}{
		{
			title: "nil input",
		},
		{
			title: "empty input",
			input: []byte(``),
		},
		{
			title: "plaintext input",
			input: []byte(`This is some data that does not look like xml.`),
		},
		{
			title: "json input",
			input: []byte(`{"This is some data": "that looks like json"}`),
		},
		{
			title: "single xml node",
			input: []byte(`<this-is-a-tag/>`),
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
				<this-is-a-tag/>
				<this-is-also-a-tag/>
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
			input: []byte(`<this-is-a-tag>This is some content.</this-is-a-tag>`),
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
			title: "single xml node with encoded content",
			input: []byte(`<this-is-a-tag>&lt;sender&gt;John Smith&lt;/sender&gt;</this-is-a-tag>`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
					Content: []byte("<sender>John Smith</sender>"),
				},
			},
		},
		{
			title: "single xml node with cdata content",
			input: []byte(`<this-is-a-tag><![CDATA[<sender>John Smith</sender>]]></this-is-a-tag>`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
					Content: []byte("<sender>John Smith</sender>"),
				},
			},
		},
		{
			title: "single xml node with attributes",
			input: []byte(`<this-is-a-tag name="my name" status="passed"></this-is-a-tag>`),
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
		{
			title: "single xml node with encoded attributes",
			input: []byte(`<this-is-a-tag name="&lt;sender&gt;John Smith&lt;/sender&gt;"></this-is-a-tag>`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
					Attrs: map[string]string{
						"name": "<sender>John Smith</sender>",
					},
				},
			},
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("#%d - %s", index+1, test.title)

		t.Run(name, func(t *testing.T) {
			actual, err := parse(bytes.NewReader(test.input))

			require.Nil(t, err)

			assert.Equal(t, test.expected, actual)
		})
	}
}
