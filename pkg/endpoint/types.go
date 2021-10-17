package endpoint

import (
	"context"
	"net/http"
)

// Creator generate endpoint or group
type Creator interface {
	CreateGroup(relPath string) (Group, error)
	CreateEndpoint(cfg Config) (Endpoint, error)
}

type PropPath interface {
	// GetRelPath relative path under group
	GetRelPath() string
	// GetAbsPath absolute path
	GetAbsPath() string
}

// Group
type Group interface {
	PropPath
	Creator
	// Parent get parent group if any else nil
	Parent() Group
	// Endpoints endpoints under current group
	Endpoints() []Endpoint
	// SubGroups sub groups
	SubGroups() []Group
}

// Endpoint handle http request
type Endpoint interface {
	PropPath
	// Group router group if any, else nil
	Group() Group
	// Handle http request
	Handle(ctx context.Context, writer http.ResponseWriter, req *http.Request)
}
