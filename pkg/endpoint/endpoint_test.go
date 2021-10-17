package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/CodeNoobKc2/goendpoint/pkg/request"
	"github.com/CodeNoobKc2/goendpoint/pkg/response"
	"github.com/CodeNoobKc2/goendpoint/utils/reflects/converter"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/CodeNoobKc2/goendpoint/pkg/endpoint/testdata"
	resptestdata "github.com/CodeNoobKc2/goendpoint/pkg/response/testdata"
)

func TestEndpointHandler(t *testing.T) {
	testcases := []struct {
		req         func() *http.Request
		handler     interface{}
		expectedRes resptestdata.MockResponseWriter
	}{
		{
			req: func() *http.Request {
				body := bytes.NewBuffer(nil)
				_ = json.NewEncoder(body).Encode(&testdata.AuthReq{Username: "foo", Passwd: "bar"})
				req, _ := http.NewRequest("POST", "/", body)
				return req
			},
			handler: testdata.MockAuth,
			expectedRes: resptestdata.MockResponseWriter{
				HttpHeader: http.Header{
					"Content-Type": {"text/plain"},
				},
				Body:       `success`,
				StatusCode: 200,
			},
		},
		{
			req: func() *http.Request {
				body := bytes.NewBuffer(nil)
				_ = json.NewEncoder(body).Encode(&testdata.AuthReq{Username: "foo", Passwd: "world"})
				req, _ := http.NewRequest("POST", "/", body)
				return req
			},
			handler: testdata.MockAuth,
			expectedRes: resptestdata.MockResponseWriter{
				HttpHeader: http.Header{
					"Content-Type": {"text/plain"},
				},
				Body:       `auth failed`,
				StatusCode: 403,
			},
		},
		{
			req: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				query := url.Values{
					"pageSize": []string{"2"},
					"pageNo":   []string{"0"},
				}
				req.URL.RawQuery = query.Encode()
				return req
			},
			handler: testdata.MockListUsers,
			expectedRes: resptestdata.MockResponseWriter{
				HttpHeader: http.Header{
					"Content-Type": {"application/json;charset-UTF8"},
				},
				Body:       `{"tot":2,"list":[{"userId":0,"username":"user-0"},{"userId":1,"username":"user-1"}]}`,
				StatusCode: 200,
			},
		},
		{
			req: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				query := url.Values{
					"pageSize":     []string{"2"},
					"pageNo":       []string{"0"},
					"blurUsername": []string{"foo"},
				}
				req.URL.RawQuery = query.Encode()
				return req
			},
			handler: testdata.MockListUsers,
			expectedRes: resptestdata.MockResponseWriter{
				HttpHeader: http.Header{
					"Content-Type": {"application/json;charset-UTF8"},
				},
				Body:       `{"tot":2,"list":[{"userId":0,"username":"foo-0"},{"userId":1,"username":"foo-1"}]}`,
				StatusCode: 200,
			},
		},
	}

	factory := &Factory{
		ReqBinderBuilder: (&request.BinderBuilder{
			Converter: converter.Builder{}.Build(),
		}).Build(),
		ResponseWriter: (&response.WriterBuilder{}).Build(),
	}

	for idx, testcase := range testcases {
		ep := MustEndpoint(factory.CreateEndpoint(Config{Path: "test", Handler: testcase.handler}))
		writer := resptestdata.NewMockResponseWriter()
		ep.Handle(context.Background(), writer, testcase.req())

		if !reflect.DeepEqual(*writer, testcase.expectedRes) {
			t.Errorf("test case '%v' failed: expected '%v' computed '%v'", idx, testcase.expectedRes, *writer)
			continue
		}
	}
}

func TestEndpointPath(t *testing.T) {
	factory := &Factory{
		ReqBinderBuilder: (&request.BinderBuilder{}).Build(),
		ResponseWriter:   (&response.WriterBuilder{}).Build(),
	}

	testcases := []struct {
		propPath    func() PropPath
		expectedAbs string
		expectedRel string
	}{
		{
			propPath: func() PropPath {
				ep := MustEndpoint(factory.CreateEndpoint(Config{Path: "foo", Handler: testdata.MockAuth}))
				return ep
			},
			expectedAbs: "/foo",
			expectedRel: "foo",
		},
		{
			propPath: func() PropPath {
				group := MustGroup(factory.CreateGroup("v1"))
				ep := MustEndpoint(group.CreateEndpoint(Config{Path: "/foo", Handler: testdata.MockAuth}))
				return ep
			},
			expectedAbs: "/v1/foo",
			expectedRel: "foo",
		},
		{
			propPath: func() PropPath {
				group0 := MustGroup(factory.CreateGroup("v1"))
				group1 := MustGroup(group0.CreateGroup("user"))
				ep := MustEndpoint(group1.CreateEndpoint(Config{Path: "auth", Handler: testdata.MockAuth}))
				return ep
			},
			expectedAbs: "/v1/user/auth",
			expectedRel: "auth",
		},
		{
			propPath: func() PropPath {
				group0 := MustGroup(factory.CreateGroup("v1"))
				group1 := MustGroup(group0.CreateGroup("user"))
				return group1
			},
			expectedAbs: "/v1/user",
			expectedRel: "user",
		},
	}

	for i, testcase := range testcases {
		computed := testcase.propPath()

		if testcase.expectedRel != computed.GetRelPath() {
			t.Errorf("test case '%v' failed: expected rel '%v' computed '%v'", i, testcase.expectedRel, computed.GetRelPath())
		}

		if testcase.expectedAbs != computed.GetAbsPath() {
			t.Errorf("test case '%v' failed: expected abs '%v' computed '%v'", i, testcase.expectedAbs, computed.GetAbsPath())
		}
	}
}
