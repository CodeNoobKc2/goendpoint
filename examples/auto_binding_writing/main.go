package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/CodeNoobKc2/goendpoint/pkg/endpoint"
)

type Payload struct {
	Hello string `json:"hello"`
}

type RequestObject struct {
	RequestId *string `header:"X-Request-Id"`
	Payload   `body:"json"`
}

type ResponseObject struct {
	*Payload `body:"json"`
	Error    error `body:"text"`
}

func Echo(ctx context.Context, req RequestObject) (resp ResponseObject) {
	if req.RequestId == nil {
		resp.Error = errors.New("invalid request")
		return
	}

	resp.Payload = &Payload{Hello: fmt.Sprintf("%v : %v", *req.RequestId, req.Hello)}
	return
}

func startHttpServer() <-chan error {
	ep := endpoint.MustEndpoint(endpoint.CreateEndpoint(endpoint.Config{Path: "echo", Handler: Echo}))
	http.Handle(ep.GetAbsPath(), endpoint.GoHttpHandler(ep))

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	done := make(chan error)
	go func() {
		defer listener.Close()
		done <- http.Serve(listener, nil)
	}()

	port := listener.(*net.TCPListener).Addr().(*net.TCPAddr).Port
	fmt.Printf("listening at port: %v\n", port)
	// you would get '{"hello":"foo : world"}'
	fmt.Printf(`try out cmd: curl -X POST localhost:%v/echo -H 'Content-Type: application/json' -H 'X-Request-ID: foo' -d '{"hello":"world"}'`+"\n", port)
	// you would get 'invalid request%'
	fmt.Printf(`try out cmd: curl -X POST localhost:%v/echo -H 'Content-Type: application/json' -d '{"hello":"world"}'`+"\n", port)
	return done
}

func main() {
	err := <-startHttpServer()
	if err != nil {
		panic(err)
		return
	}
}
