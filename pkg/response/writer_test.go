package response

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/CodeNoobKc2/goendpoint/pkg/response/testdata"
)

func TestResponseWriter(t *testing.T) {
	var (
		testAlwaysToken  = "always-token"
		testSuccessToken = "success-token"
	)

	testcases := []struct {
		input            testdata.Response
		expectedResponse testdata.MockResponseWriter
	}{
		{
			input: testdata.Response{
				AlwaysToken:  &testAlwaysToken,
				SuccessToken: &testSuccessToken,
				Body:         &testdata.SuccessResponseBody{Foo: "foo", Bar: "bar"},
			},
			expectedResponse: testdata.MockResponseWriter{
				HttpHeader: http.Header{
					"A-Token":      []string{testAlwaysToken},
					"S-Token":      []string{testSuccessToken},
					"Content-Type": []string{applicationJSON},
				},
				Body:       `{"foo":"foo","bar":"bar"}`,
				StatusCode: 200,
			},
		},
		{
			input: testdata.Response{
				AlwaysToken: &testAlwaysToken,
				ErrMsg:      errors.New("test response 400"),
			},
			expectedResponse: testdata.MockResponseWriter{
				HttpHeader: http.Header{
					"A-Token":      []string{testAlwaysToken},
					"Content-Type": []string{textPlain},
				},
				Body:       `test response 400`,
				StatusCode: 400,
			},
		},
		{
			input: testdata.Response{
				AlwaysToken: &testAlwaysToken,
				Forbidden:   testdata.Forbidden{Msg: "you shall not pass"},
			},
			expectedResponse: testdata.MockResponseWriter{
				HttpHeader: http.Header{
					"A-Token":      []string{testAlwaysToken},
					"Content-Type": []string{textPlain},
				},
				Body:       `you shall not pass`,
				StatusCode: 403,
			},
		},
		{
			input: testdata.Response{
				ErrMsg: errors.New("test response 400"),
			},
			expectedResponse: testdata.MockResponseWriter{
				HttpHeader: http.Header{
					"Content-Type": []string{textPlain},
				},
				Body:       `test response 400`,
				StatusCode: 400,
			},
		},
		{
			input: testdata.Response{
				Body:   &testdata.SuccessResponseBody{},
				ErrMsg: errors.New("test response 400"),
			},
			expectedResponse: testdata.MockResponseWriter{
				HttpHeader: http.Header{
					"Content-Type": []string{textPlain},
				},
				Body:       errMultipleStatusCodeWouldBeWritten.Error(),
				StatusCode: 500,
			},
		},
		{
			input: testdata.Response{},
			expectedResponse: testdata.MockResponseWriter{
				HttpHeader: http.Header{
					"Content-Type": []string{textPlain},
				},
				Body:       errNoResponseToWritten.Error(),
				StatusCode: 500,
			},
		},
		{
			input: testdata.Response{
				Body:         &testdata.SuccessResponseBody{},
				ConflictBody: &testdata.SuccessResponseBody{},
			},
			expectedResponse: testdata.MockResponseWriter{
				HttpHeader: http.Header{
					"Content-Type": []string{textPlain},
				},
				Body:       errMultipleBodyWouldBeWritten.Error(),
				StatusCode: 500,
			},
		},
	}

	writer := WriterBuilder{}.Build()
	for i, testcase := range testcases {
		computed := testdata.NewMockResponseWriter()
		writer.Write(computed, testcase.input)

		if !reflect.DeepEqual(*computed, testcase.expectedResponse) {
			t.Errorf("test case '%v' failed: expected '%v' computed '%v'", i, testcase.expectedResponse, *computed)
			continue
		}
	}
}
