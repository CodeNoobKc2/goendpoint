package request

import (
	"errors"
	"reflect"
	"testing"
)

func TestPathTemplate(t *testing.T) {
	testcases := []struct {
		template  string
		url       string
		expectErr error
		parsed    map[string][]string
	}{
		{
			template:  "/v1/{ foo }/bar",
			url:       "/v1/hello/bar",
			expectErr: nil,
			parsed:    map[string][]string{"foo": {"hello"}},
		},
		{
			template:  "/v1/{ foo }/bar",
			url:       "/v1/hello/world",
			expectErr: errors.New("uri '/v1/hello/world' does not match pattern '/v1/{ foo }/bar'"),
		},
		{
			template: "/v1/{ foo}/foo/{bar }/bar",
			url:      "/v1/hello/foo/world/bar",
			parsed:   map[string][]string{"foo": {"hello"}, "bar": {"world"}},
		},
	}

	for idx, testcase := range testcases {
		tmpl, err := NewPathTemplate(testcase.template)
		if err != nil {
			return
		}
		computed, err := tmpl.Parse(testcase.url)
		switch {
		case testcase.expectErr == nil && err != nil:
			t.Errorf("test case '%v' failed: expected err 'nil' computed err '%v'", idx, err)
			continue
		case testcase.expectErr != nil && err == nil:
			t.Errorf("test case '%v' failed: expected err '%v' computed err 'nil'", idx, err)
			continue
		case testcase.expectErr != nil && err != nil:
			if testcase.expectErr.Error() != err.Error() {
				t.Errorf("test case '%v' failed: expected err '%v' computed err '%v'", idx, testcase.expectErr, err)
				continue
			}
		}

		if !reflect.DeepEqual(computed, testcase.parsed) {
			t.Errorf("test case '%v' failed: expected parsed '%v' computed parsed '%v'", idx, testcase.parsed, computed)
		}
	}
}
