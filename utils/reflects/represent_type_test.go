package reflects

import (
	"reflect"
	"testing"

	"github.com/CodeNoobKc2/goendpoint/utils/reflects/converter"
)

const expectedConverterRepresent = "github.com/CodeNoobKc2/goendpoint/utils/reflects/converter.Converter"

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
			input:    map[string]converter.Converter{},
			expected: "map[string]" + expectedConverterRepresent,
		},
		{
			input:    map[string]*converter.Converter{},
			expected: "map[string]*" + expectedConverterRepresent,
		},
		{
			input:    [2]*converter.Converter{},
			expected: "[2]*" + expectedConverterRepresent,
		},
		{
			input:    converter.Converter{},
			expected: expectedConverterRepresent,
		},
		{
			input:    &converter.Converter{},
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
