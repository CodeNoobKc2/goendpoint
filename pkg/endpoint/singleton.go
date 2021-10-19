package endpoint

import (
	"github.com/CodeNoobKc2/goendpoint/pkg/request"
	"github.com/CodeNoobKc2/goendpoint/pkg/response"
	"github.com/CodeNoobKc2/goendpoint/utils/reflects/converter"
	"net/http"
)

var factoryInstance = &Factory{
	ReqBinderBuilder: (&request.BinderBuilder{Converter: converter.Builder{}.Build()}).Build(),
	ResponseWriter:   (&response.WriterBuilder{}).Build(),
	CrashHandler:     nil,
}

func CreateGroup(relPath string) (Group, error) {
	return factoryInstance.CreateGroup(relPath)
}

func CreateEndpoint(cfg Config) (Endpoint, error) {
	return factoryInstance.CreateEndpoint(cfg)
}

func GoHttpHandler(endpoint Endpoint) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		endpoint.Handle(r.Context(), writer, r)
	})
}
