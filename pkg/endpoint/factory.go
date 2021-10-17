package endpoint

import (
	"errors"
	"net/http"
	gopath "path"
	"strings"

	"github.com/CodeNoobKc2/goendpoint/pkg/request"
	"github.com/CodeNoobKc2/goendpoint/pkg/response"
)

var (
	errHandlerIsNotAsExpected    = errors.New("err handler should be an function with signature func(ctx context.Context,req ReqObj)(resp RespObj)")
	errGroupPathShouldNotBeEmpty = errors.New("err group path should be not be empty")
)

func MustEndpoint(ep Endpoint, err error) Endpoint {
	if err != nil {
		panic(err)
	}
	return ep
}

func MustGroup(group Group, err error) Group {
	if err != nil {
		panic(err)
	}
	return group
}

var _ Creator = &Factory{}

// Factory initialize endpoint correctly
type Factory struct {
	// ReqBinderBuilder request binder is responsible for bind
	ReqBinderBuilder request.Binder
	// ResponseWriter response writer is responsible for writer response
	ResponseWriter response.Writer
	// CrashHandler crash handler on handler http request failed
	CrashHandler func(writer http.ResponseWriter, recover interface{})
}

func (f *Factory) CreateGroup(path string) (Group, error) {
	path = strings.Trim(path, "/")
	if len(path) == 0 {
		return nil, errGroupPathShouldNotBeEmpty
	}
	return &group{factory: f, parent: nil, relPath: path, absPath: gopath.Join("/", path)}, nil
}

func (f *Factory) CreateEndpoint(cfg Config) (Endpoint, error) {
	cfg.factory = f
	return newEndpoint(cfg)
}

var _ Group = &group{}

type group struct {
	factory *Factory
	// parent group
	parent Group
	// relPath to parent group, if Parent is nil, then this means absPath
	relPath string
	// absPath
	absPath string

	subGroups []Group
	endpoints []Endpoint
}

func (g *group) CreateGroup(relPath string) (Group, error) {
	relPath = strings.Trim(relPath, "/")
	if len(relPath) == 0 {
		return nil, errGroupPathShouldNotBeEmpty
	}

	sg := &group{factory: g.factory, parent: g, relPath: relPath, absPath: gopath.Join(g.absPath, relPath)}
	g.subGroups = append(g.subGroups, sg)
	return sg, nil
}

func (g *group) CreateEndpoint(cfg Config) (Endpoint, error) {
	cfg.factory = g.factory
	cfg.parentGroup = g
	ep, err := newEndpoint(cfg)
	if err != nil {
		return nil, err
	}
	g.endpoints = append(g.endpoints, ep)
	return ep, nil
}

func (g *group) Parent() Group {
	return g.parent
}

func (g *group) Endpoints() []Endpoint {
	return g.endpoints
}

func (g *group) SubGroups() []Group {
	return g.subGroups
}

func (g *group) GetAbsPath() string {
	return g.absPath
}

func (g *group) GetRelPath() string {
	return g.relPath
}

func (g *group) RouterGroup(path string) *group {
	return &group{factory: g.factory, parent: g, relPath: path, absPath: gopath.Join(g.absPath, path)}
}
