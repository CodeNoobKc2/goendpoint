package request

import (
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
				req, _ := http.NewRequest("GET", "/foo", nil)
				url := url.Values{
					"q": []string{"foo"},
				}
				req.URL.RawQuery = url.Encode()
				return req
			},
			expected: struct {
				Query string `query:"q"`
			}{
				Query: "foo",
			},
		},
		{
			pathTemplate: "/foo/{int}",
			request: func() *http.Request {
				req, _ := http.NewRequest("GET", "/foo/1", nil)
				return req
			},
			expected: struct {
				Int int `path:"int"`
			}{
				Int: 1,
			},
		},
	}

	binderBuilder := &BinderBuilder{Converter: converter.Builder{}.Build()}
	binder := binderBuilder.Build()

	for idx, testcase := range testcases {
		out := reflect.New(reflect.TypeOf(testcase.expected))
		if err := binder.WithPath(testcase.pathTemplate)(testcase.request(), out.Interface()); err != nil {
			t.Errorf("test case '%v' failed: %v", idx, err)
			continue
		}

		if !reflect.DeepEqual(out.Elem().Interface(), testcase.expected) {
			t.Errorf("test case '%v' failed: expected %v;computed %v", idx, testcase.expected, out.Elem().Interface())
			continue
		}
	}
}
