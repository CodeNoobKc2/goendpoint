package reflects

import (
	"reflect"
	"testing"

	"github.com/CodeNoobKc2/goendpoint/pkg/swagger"
)

const expectedConverterRepresent = "github.com/CodeNoobKc2/goendpoint/pkg/swagger.Schema"

func TestRepresentType(t *testing.T) {
	testcases := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    "",
			expected: "string",
		},
		{
			input:    new(string),
			expected: "*string",
		},
		{
			input:    new(string),
			expected: "*string",
		},
		{
			input:    map[string]int{},
			expected: "map[string]int",
		},
		{
			input: struct {
				Foo   string `foo:"t"`
				Bar   []int  `bar:"b"`
				Hello int
			}{},
			expected: "struct { Foo string \"foo:\\\"t\\\"\"; Bar []int \"bar:\\\"b\\\"\"; Hello int }",
		},
		{
			input: struct {
				Foo         string `foo:"t"`
				Complicated swagger.Schema
				Bar         []int `bar:"b"`
			}{},
			expected: "struct { Foo string \"foo:\\\"t\\\"\"; Complicated github.com/CodeNoobKc2/goendpoint/pkg/swagger.Schema ; Bar []int \"bar:\\\"b\\\"\"}",
		},
		{
			input:    map[string]swagger.Schema{},
			expected: "map[string]" + expectedConverterRepresent,
		},
		{
			input:    map[string]*swagger.Schema{},
			expected: "map[string]*" + expectedConverterRepresent,
		},
		{
			input:    [2]*swagger.Schema{},
			expected: "[2]*" + expectedConverterRepresent,
		},
		{
			input:    swagger.Schema{},
			expected: expectedConverterRepresent,
		},
		{
			input:    &swagger.Schema{},
			expected: "*" + expectedConverterRepresent,
		},
	}

	for _, testcase := range testcases {
		computed := RepresentType(reflect.TypeOf(testcase.input))
		if computed != testcase.expected {
			t.Errorf("expected '%v' computed '%v'", testcase.expected, computed)
		}
	}
}
