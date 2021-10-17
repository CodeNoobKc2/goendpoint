package endpoint

import (
	"context"
	"fmt"
	"net/http"
	gopath "path"
	"reflect"
	"strings"

	"github.com/CodeNoobKc2/goendpoint/utils/reflects"
)

var _ Endpoint = &endpoint{}

type Config struct {
	factory     *Factory
	parentGroup Group
	// Description description of current endpoint
	Description string
	// Path rel api path under certain api group
	Path string
	// Tags current endpoint tag
	Tags []string
	// Handler the underlying Handler value must be a function with signature func(ctx context.Context,req ReqObject)RespObject
	Handler interface{}
}

func (cfg Config) absPath() string {
	if cfg.parentGroup != nil {
		return gopath.Join(cfg.parentGroup.GetAbsPath(), cfg.Path)
	} else {
		return "/" + cfg.Path
	}
}

func newEndpoint(c Config) (*endpoint, error) {
	c.Path = strings.Trim(c.Path, "/")
	if c.Handler == nil {
		return nil, errHandlerIsNotAsExpected
	}

	handler, err := c.makeHttpHandler(c.Handler)
	if err != nil {
		return nil, err
	}

	return &endpoint{
		group:   c.parentGroup,
		absPath: c.absPath(),
		relPath: c.Path,
		handle:  handler,
	}, nil
}

func (cfg Config) validateHandler(t reflect.Type) error {
	if t.Kind() != reflect.Func {
		return errHandlerIsNotAsExpected
	}

	if t.NumIn() != 2 || t.NumOut() != 1 {
		return errHandlerIsNotAsExpected
	}

	if reflects.RepresentType(t.In(0)) != "context.Context" {
		return errHandlerIsNotAsExpected
	}

	if err := reflects.ShouldBeKind(t.In(1), reflect.Struct); err != nil {
		return errHandlerIsNotAsExpected
	}

	if err := reflects.ShouldBeKind(t.Out(0), reflect.Struct); err != nil {
		return errHandlerIsNotAsExpected
	}

	return nil
}

func (cfg Config) makeHttpHandler(handler interface{}) (func(ctx context.Context, req *http.Request, writer http.ResponseWriter), error) {
	v := reflect.ValueOf(handler)
	if err := cfg.validateHandler(v.Type()); err != nil {
		return nil, err
	}

	newReqObject := func() reflect.Value {
		return reflect.New(v.Type().In(1))
	}
	bindFunc := cfg.factory.ReqBinderBuilder.WithPath(gopath.Join(cfg.absPath(), cfg.Path))

	crashHandler := cfg.factory.CrashHandler
	if crashHandler == nil {
		crashHandler = func(writer http.ResponseWriter, recover interface{}) {
			writer.Write([]byte(fmt.Sprintf("%s", recover)))
			writer.WriteHeader(500)
		}
	}

	return func(ctx context.Context, req *http.Request, writer http.ResponseWriter) {
		defer func() {
			if r := recover(); r != nil {
				crashHandler(writer, r)
			}
		}()

		// bind request object
		vptrReqObj := newReqObject()
		if err := bindFunc(req, vptrReqObj.Interface()); err != nil {
			panic(err)
			return
		}
		// call handler
		rtns := v.Call([]reflect.Value{reflect.ValueOf(ctx), vptrReqObj.Elem()})
		// write response
		cfg.factory.ResponseWriter.Write(writer, rtns[0].Interface())
	}, nil
}

// endpoint config contains all the necessary info to generate swagger doc as well as handler request properly
type endpoint struct {
	group   Group
	absPath string
	relPath string
	handle  func(ctx context.Context, req *http.Request, writer http.ResponseWriter)
}

func (e *endpoint) Group() Group {
	return e.group
}

// AbsPath absolute path of current endpoint
func (e *endpoint) GetAbsPath() string {
	return e.absPath
}

func (e *endpoint) GetRelPath() string {
	return e.relPath
}

func (e *endpoint) Handle(ctx context.Context, writer http.ResponseWriter, req *http.Request) {
	e.handle(ctx, req, writer)
}
