package request

import (
	"bytes"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/CodeNoobKc2/goendpoint/pkg/request/testdata"
	"github.com/CodeNoobKc2/goendpoint/utils/reflects/converter"
)

func TestBinding(t *testing.T) {
	optionalStr := "bar"
	optionalInt := 2

	testcases := []struct {
		pathTemplate string
		request      func() *http.Request
		expected     interface{}
	}{
		// header
		{
			pathTemplate: "/foo",
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/foo", nil)
				req.Header.Set("Str", "foo")
				req.Header.Set("Optional-Str", optionalStr)
				req.Header.Set("Int", "1")
				return req
			},
			expected: testdata.BindObject{
				InHeader: testdata.InHeader{
					HeaderStr:         "foo",
					HeaderOptionalStr: &optionalStr,
					HeaderInt:         1,
				},
			},
		},
		{
			pathTemplate: "/foo",
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/foo", nil)
				req.URL.RawQuery = url.Values{
					"str":         []string{"foo"},
					"int":         []string{"1"},
					"optionalStr": []string{optionalStr},
					"optionalInt": []string{"2"},
					"strSlice":    []string{"foo", "bar"},
				}.Encode()
				return req
			},
			expected: testdata.BindObject{
				InQuery: testdata.InQuery{
					QueryStr:         "foo",
					QueryInt:         1,
					QueryOptionalStr: &optionalStr,
					QueryOptionalInt: &optionalInt,
					QueryStrSlice:    []string{"foo", "bar"},
				},
			},
		},
		{
			pathTemplate: "/foo/{str}/{int}",
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/foo/bar/1", nil)
				return req
			},
			expected: testdata.BindObject{
				InPath: testdata.InPath{
					PathInt: 1,
					PathStr: "bar",
				},
			},
		},
		{
			pathTemplate: "/foo",
			request: func() *http.Request {
				req, _ := http.NewRequest("POST", "/foo/bar", bytes.NewBuffer([]byte(`{"foo":"foo","bar":"bar"}`)))
				return req
			},
			expected: testdata.BindObject{
				RequestBody: testdata.RequestBody{
					Foo: "foo",
					Bar: "bar",
				},
			},
		},
	}

	binderBuilder := &BinderBuilder{Converter: converter.Builder{}.Build()}
	for idx, testcase := range testcases {
		binder := binderBuilder.Build()
		clean := &testdata.BindObject{}
		if err := binder.WithPath(testcase.pathTemplate)(testcase.request(), clean); err != nil {
			t.Errorf("test case '%v' failed: %v", idx, err)
			continue
		}

		if !reflect.DeepEqual(*clean, testcase.expected) {
			t.Errorf("test case '%v' failed: expected %v;computed %v", idx, testcase.expected, *clean)
			continue
		}
	}
}
