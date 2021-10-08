package request

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/CodeNoobKc2/goendpoint/utils/reflects"
	"github.com/CodeNoobKc2/goendpoint/utils/reflects/converter"
)

const (
	hasQueryTag  = 1
	hasBodyTag   = 1 << 1
	hasPathTag   = 1 << 2
	hasHeaderTag = 1 << 3
)

type BinderBuilder struct {
	// PathTag tag for path Key, default would be "path"
	PathTag *string
	// QueryTag tag for query Key, default would be "query"
	QueryTag *string
	// HeaderTag tag for query Key, default would be "header"
	HeaderTag *string
	// HeaderTag tag for body Key, default would be "body"
	BodyTag *string

	Converter *converter.Converter
}

func (b *BinderBuilder) Build(pathTemplate string) *binder {
	pathKeyToPos := make(map[string]int)
	for idx, chunk := range strings.Split(strings.TrimPrefix(pathTemplate, "/"), "/") {
		if strings.HasPrefix(chunk, "{") && strings.HasSuffix(chunk, "}") {
			pathKeyToPos[chunk[1:len(chunk)-1]] = idx
		}
	}

	return &binder{pathKeyToPos: pathKeyToPos, lock: &sync.RWMutex{}, parsedRequestObjects: map[string]ParsedRequestObject{}, BinderBuilder: b}
}

func (b *BinderBuilder) GetQueryTag() string {
	if b.QueryTag != nil {
		return *b.QueryTag
	}
	return "query"
}

func (b *BinderBuilder) GetPathTag() string {
	if b.PathTag != nil {
		return *b.PathTag
	}
	return "path"
}

func (b *BinderBuilder) GetHeaderTag() string {
	if b.HeaderTag != nil {
		return *b.HeaderTag
	}
	return "header"
}

func (b *BinderBuilder) GetBodyTag() string {
	if b.BodyTag != nil {
		return *b.BodyTag
	}
	return "body"
}

type binder struct {
	*BinderBuilder
	lock                 *sync.RWMutex
	parsedRequestObjects map[string]ParsedRequestObject
	pathKeyToPos         map[string]int
}

func (b binder) ParsePath(uri string) map[string][]string {
	// todo: we need to validate whether uri match the pattern of pathTemplate
	if len(b.pathKeyToPos) == 0 {
		return nil
	}
	ret := make(map[string][]string)
	splited := strings.Split(strings.TrimPrefix(uri, "/"), "/")

	for key, i := range b.pathKeyToPos {
		ret[key] = []string{splited[i]}
	}

	return ret
}

func (b binder) bindFromParsed(ctx BindContext, in string, m map[string][]string, to reflect.Value) error {
	values := m[ctx.Key]
	switch {
	// validate would not be preformed here
	case len(values) == 0:
		return nil
		// TODO: parse list according to style field in the future, during pof this is not that important
	case len(values) > 1 && reflects.IsList(to.Type()):
		if err := b.Converter.ConvertValue(reflect.ValueOf(values), to); err != nil {
			return fmt.Errorf("failed to bind %v '%v': %v", in, ctx.Key, err)
		}
		return nil
	case len(values) > 1 && !reflects.IsList(to.Type()):
		return fmt.Errorf("multiple value occurred on %v Key '%v'", in, ctx.Key)
	default:
		if err := b.Converter.ConvertValue(reflect.ValueOf(values[0]), to); err != nil {
			return fmt.Errorf("failed to bind %v '%v': %v", in, ctx.Key, err)
		}
		return nil
	}
}

func (b binder) bindHeader(ctx BindContext, req *http.Request, to reflect.Value) error {
	return b.bindFromParsed(ctx, "header", req.Header, to)
}

func (b binder) bindQuery(ctx BindContext, req *http.Request, to reflect.Value) error {
	return b.bindFromParsed(ctx, "query", req.URL.Query(), to)
}

func (b binder) bindPath(ctx BindContext, req *http.Request, to reflect.Value) error {
	return b.bindFromParsed(ctx, "path", b.ParsePath(req.URL.Path), to)
}

func (b binder) bindBody(ctx BindContext, req *http.Request, to reflect.Value) error {
	// validate would not be performed here
	if req.Body == nil {
		return nil
	}
	receiver := reflect.New(to.Type())
	switch {
	case ctx.Key == "json":
		if err := json.NewDecoder(req.Body).Decode(receiver.Interface()); err != nil {
			return fmt.Errorf("error decode json body: %v", err)
		}
		to.Set(receiver.Elem())
		return nil
	case ctx.Key == "xml":
		if err := xml.NewDecoder(req.Body).Decode(receiver.Interface()); err != nil {
			return fmt.Errorf("error decode xml body: %v", err)
		}
		to.Set(receiver.Elem())
		return nil
	case ctx.Key == "form":
		panic("not implemented")
	default:
		return fmt.Errorf("unknown body kind '%v'", ctx.Key)
	}
}

func (b binder) parseRequestObject(tobj reflect.Type, parsed *ParsedRequestObject) error {
	var has bool
	ref := tobj.PkgPath() + "." + tobj.Name()
	func() {
		b.lock.RLock()
		defer b.lock.RUnlock()

		if _, has = b.parsedRequestObjects[ref]; has {
			*parsed = b.parsedRequestObjects[ref]
		}
	}()

	if has {
		return nil
	}

	parsing := NewParsedRequestObject()
	for i := 0; i < tobj.NumField(); i++ {
		field := tobj.Field(i)
		if !field.IsExported() {
			continue
		}

		// this is an embedded struct
		if len(field.Tag) == 0 && field.Anonymous {
			embedded := NewParsedRequestObject()
			if err := b.parseRequestObject(field.Type, embedded); err != nil {
				return err
			}
			for name, context := range embedded.FieldToBindContext {
				parsing.FieldToBindContext[name] = context
			}
			continue
		}

		tagFlag := 0
		if _, ok := field.Tag.Lookup(b.GetQueryTag()); ok {
			tagFlag += hasQueryTag
		}

		if _, ok := field.Tag.Lookup(b.GetPathTag()); ok {
			tagFlag += hasPathTag
		}

		if _, ok := field.Tag.Lookup(b.GetHeaderTag()); ok {
			tagFlag += hasHeaderTag
		}

		if _, ok := field.Tag.Lookup(b.GetBodyTag()); ok {
			tagFlag += hasBodyTag
		}

		bindCtx := BindContext{}
		switch {
		case tagFlag^hasQueryTag == 0:
			bindCtx.In = "query"
			tags := strings.Split(field.Tag.Get(b.GetQueryTag()), ",")
			if len(tags[0]) == 0 {
				bindCtx.Key = strings.ToLower(field.Name[:1]) + field.Name[1:]
			} else {
				bindCtx.Key = tags[0]
			}
		case tagFlag^hasPathTag == 0:
			bindCtx.In = "path"
			tags := strings.Split(field.Tag.Get(b.GetPathTag()), ",")
			if len(tags[0]) == 0 {
				bindCtx.Key = strings.ToLower(field.Name[:1]) + field.Name[1:]
			} else {
				bindCtx.Key = tags[0]
			}
		case tagFlag^hasHeaderTag == 0:
			bindCtx.In = "header"
			tags := strings.Split(field.Tag.Get(b.GetHeaderTag()), ",")
			if len(tags[0]) == 0 {
				bindCtx.Key = field.Name
			} else {
				bindCtx.Key = tags[0]
			}
		case tagFlag^hasBodyTag == 0:
			bindCtx.In = "body"
			tags := strings.Split(field.Tag.Get(b.GetHeaderTag()), ",")
			if len(tags[0]) == 0 {
				bindCtx.Key = "json"
			} else {
				bindCtx.Key = tags[0]
				if bindCtx.Key != "json" && bindCtx.Key != "xml" && bindCtx.Key != "form" {
					return fmt.Errorf("unknown body kind '%v'", tags[0])
				}
			}
		default:
			return fmt.Errorf("struct field '%v' may has multiple in tag", field.Name)
		}
		parsing.FieldToBindContext[field.Name] = bindCtx
	}

	func() {
		b.lock.Lock()
		defer b.lock.Unlock()

		if obj, ok := b.parsedRequestObjects[ref]; ok {
			*parsed = obj
		} else {
			b.parsedRequestObjects[ref] = *parsing
			*parsed = *parsing
		}

	}()
	return nil
}

func (b binder) bind(req *http.Request, obj reflect.Value) (err error) {
	parsed := NewParsedRequestObject()
	if err := b.parseRequestObject(obj.Type(), parsed); err != nil {
		return err
	}

	for name, context := range parsed.FieldToBindContext {
		// inject request uri
		field, err := reflects.GetStructMember(obj, name, true)
		if err != nil {
			return err
		}

		switch context.In {
		case "header":
			if err := b.bindHeader(context, req, *field); err != nil {
				return err
			}
		case "query":
			if err := b.bindQuery(context, req, *field); err != nil {
				return err
			}
		case "path":
			if err := b.bindPath(context, req, *field); err != nil {
				return err
			}
		case "body":
			if err := b.bindBody(context, req, *field); err != nil {
				return err
			}
		default:
			panic(fmt.Sprintf("in '%v' is not supported for binding", context.In))
		}
	}

	return nil
}

// Bind would fill in corresponding from user request to given struct
// obj should be a pointer to struct
func (b binder) Bind(req *http.Request, obj interface{}) error {
	v := reflect.ValueOf(obj)
	if err := reflects.ShouldBeKind(v.Type(), reflect.Ptr); err != nil {
		return err
	}

	if err := reflects.ShouldBeKind(v.Type().Elem(), reflect.Struct); err != nil {
		return err
	}

	return b.bind(req, v.Elem())
}
