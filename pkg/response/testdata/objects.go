package testdata

import "net/http"

var _ http.ResponseWriter = &MockResponseWriter{}

type MockResponseWriter struct {
	HttpHeader http.Header
	Body       string
	StatusCode int
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{HttpHeader: http.Header{}, Body: "", StatusCode: 0}
}

func (m *MockResponseWriter) Header() http.Header {
	return m.HttpHeader
}

func (m *MockResponseWriter) Write(bytes []byte) (int, error) {
	m.Body = string(bytes)
	return 0, nil
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.StatusCode = statusCode
}

type SuccessResponseBody struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

type Forbidden struct {
	Msg string `body:"text" code:"403"`
}

// Response example of response object
type Response struct {
	Forbidden
	// if code tag is "default" then by OAS definition, it would be presented by any kind of response
	AlwaysToken  *string              `header:"A-Token" code:"default"`
	SuccessToken *string              `header:"S-Token"`
	Body         *SuccessResponseBody `body:"json"`
	ConflictBody *SuccessResponseBody `body:"json"`
	ErrMsg       error                `body:"text"`
}
