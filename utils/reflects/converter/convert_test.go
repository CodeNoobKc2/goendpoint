package converter

import (
	"reflect"
	"testing"
)

func TestConvert(t *testing.T) {
	s := "1"
	i := 1
	testcases := []struct {
		input    interface{}
		expected interface{}
	}{
		{
			input:    s,
			expected: i,
		},
		{
			input:    &s,
			expected: &i,
		},
		{
			input:    s,
			expected: &i,
		},
		{
			input:    "1.321",
			expected: float32(1.321),
		},
		{
			input:    "123",
			expected: uint(123),
		},
		{
			input:    int32(32),
			expected: float32(32),
		},
		{
			input:    []interface{}{"1", "2", "3"},
			expected: []int{1, 2, 3},
		},
		{
			input:    []interface{}{"1", "2", "3"},
			expected: [3]int{1, 2, 3},
		},
		{
			input:    [2]interface{}{"1", "2"},
			expected: [3]int{1, 2},
		},
		{
			input:    []string{"1.1", "1.2", "1.3"},
			expected: []float32{1.1, 1.2, 1.3},
		},
		{
			input:    map[interface{}]interface{}{"1": "1.1", "2": "1.2", "3": "1.3"},
			expected: map[int]float32{1: 1.1, 2: 1.2, 3: 1.3},
		},
		{
			input:    map[interface{}][]interface{}{"1": {&s, &s}, "2": {&s, &s}},
			expected: map[float32][]*int{1: {&i, &i}, 2: {&i, &i}},
		},
		{
			input:    map[string]float32{"1": 1.1, "2": 1.2, "3": 1.3},
			expected: map[interface{}]interface{}{"1": float32(1.1), "2": float32(1.2), "3": float32(1.3)},
		},
		{
			input: RawJSON(`{"foo":"bar"}`),
			expected: struct {
				Foo string `json:"foo"`
			}{Foo: "bar"},
		},
	}

	converter := Builder{}.Build()

	for _, testcase := range testcases {
		out := reflect.New(reflect.TypeOf(testcase.expected)).Interface()
		if err := converter.Convert(testcase.input, out); err != nil {
			t.Errorf("convert from '%v' to '%v' failed: %v", testcase.input, testcase.expected, err)
			continue
		}

		outVal := reflect.ValueOf(out).Elem().Interface()
		if !reflect.DeepEqual(outVal, testcase.expected) {
			t.Errorf("convert from '%v' to '%v' failed: got '%v'", testcase.input, testcase.expected, outVal)
			continue
		}
	}
}
