package response

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"sync"

	"github.com/CodeNoobKc2/goendpoint/utils/reflects"
	"github.com/CodeNoobKc2/goendpoint/utils/sets"
)

const (
	headerContentType = "Content-Type"
	applicationJSON   = "application/json;charset-UTF8"
	applicationXML    = "application/xml;charset-UTF8"
	textPlain         = "text/plain"
)

var (
	errNoResponseToWritten              = errors.New("no response to written")
	errMultipleStatusCodeWouldBeWritten = errors.New("multiple status code would be written")
	errMultipleBodyWouldBeWritten       = errors.New("multiple body would be written")
	defaultMapShortContentTypeToFull    = map[string]string{
		"text": textPlain,
		"json": applicationJSON,
		"xml":  applicationXML,
	}
)

type WriterBuilder struct {
	// CodeTag default is "code"
	CodeTag *string
	// BodyTag default is "body"
	BodyTag *string
	// HeaderTag default is "header"
	HeaderTag *string
	// DefaultShortContentType default is "json"
	DefaultShortContentType *string
	// ShortContentTypeToFull default is
	// json->application/json;charset=UTF-8
	// xml->application/xml;charset=UTF-8
	ShortContentTypeToFull map[string]string
	// DefaultNonErrHttpStatusCode if no code tag is given and type is not error,then use this code
	DefaultNonErrHttpStatusCode *int
	// DefaultErrHttpStatusCode if no code tag is given and type is error, then use this code
	DefaultErrHttpStatusCode *int
	// OnFailedToWriteResponse if write response failed, then handle the failure
	OnFailedToWriteResponse func(writer http.ResponseWriter, err error)
}

func (w WriterBuilder) GetDefaultShortContentType() string {
	if w.DefaultShortContentType == nil {
		return "json"
	}
	return *w.DefaultShortContentType
}

func (w WriterBuilder) GetCodeTag() string {
	if w.CodeTag == nil {
		return "code"
	}
	return *w.CodeTag
}

func (w WriterBuilder) GetBodyTag() string {
	if w.BodyTag == nil {
		return "body"
	}
	return *w.BodyTag
}

func (w WriterBuilder) GetHeaderTag() string {
	if w.HeaderTag == nil {
		return "header"
	}
	return *w.HeaderTag
}

func (w WriterBuilder) GetDefaultNonErrHttpStatusCode() int {
	if w.DefaultNonErrHttpStatusCode == nil {
		return 200
	}
	return *w.DefaultNonErrHttpStatusCode
}

func (w WriterBuilder) GetDefaultErrHttpStatusCode() int {
	if w.DefaultErrHttpStatusCode == nil {
		return 400
	}
	return *w.DefaultNonErrHttpStatusCode
}

func (w WriterBuilder) Build() *writer {
	shortContentTypeToFull := w.ShortContentTypeToFull
	if len(shortContentTypeToFull) == 0 {
		shortContentTypeToFull = defaultMapShortContentTypeToFull
	}

	onFailedToWriteResponse := w.OnFailedToWriteResponse
	if onFailedToWriteResponse == nil {
		onFailedToWriteResponse = func(writer http.ResponseWriter, err error) {
			writer.Write([]byte(err.Error()))
			writer.Header().Set(headerContentType, textPlain)
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}

	return &writer{
		lock:                        &sync.RWMutex{},
		parsedResponse:              map[string]ParsedResponse{},
		shortContentTypeToFull:      shortContentTypeToFull,
		bodyTag:                     w.GetBodyTag(),
		headerTag:                   w.GetHeaderTag(),
		codeTag:                     w.GetCodeTag(),
		defaultShortContentType:     w.GetDefaultShortContentType(),
		defaultErrHttpStatusCode:    w.GetDefaultErrHttpStatusCode(),
		defaultNonErrHttpStatusCode: w.GetDefaultNonErrHttpStatusCode(),
		onFailedToWriteResponse:     onFailedToWriteResponse,
	}
}

type Writer interface {
	Write(writer http.ResponseWriter, response interface{})
}

type writer struct {
	lock                        *sync.RWMutex
	parsedResponse              map[string]ParsedResponse
	shortContentTypeToFull      map[string]string
	bodyTag                     string
	headerTag                   string
	codeTag                     string
	defaultShortContentType     string
	defaultErrHttpStatusCode    int
	defaultNonErrHttpStatusCode int
	onFailedToWriteResponse     func(writer http.ResponseWriter, err error)
}

func (w writer) GetFullContentType(short string) string {
	full, has := w.shortContentTypeToFull[short]
	if !has {
		return defaultMapShortContentTypeToFull[short]
	}
	return full
}

func (w *writer) parseField(tfield reflect.StructField) (*ParsedField, error) {
	contentType, inBody := tfield.Tag.Lookup(w.bodyTag)
	headerKey, inHeader := tfield.Tag.Lookup(w.headerTag)

	if inBody && inHeader {
		return nil, fmt.Errorf("response object field '%v' has both body tag and header tag", tfield.Name)
	}

	if !inBody && !inHeader {
		return nil, fmt.Errorf("response object field '%v' underterminted body or header", tfield.Name)
	}

	// get code
	code := w.defaultNonErrHttpStatusCode
	if strCode := tfield.Tag.Get(w.codeTag); len(strCode) != 0 {
		// for all response
		if strCode == "default" {
			code = -1
		} else {
			int64Code, err := strconv.ParseInt(strCode, 0, 0)
			if err != nil {
				return nil, fmt.Errorf("cannot format code as int: %v", err)
			}
			if int64Code < 0 {
				return nil, fmt.Errorf("code is '%v'. it should not < 0", int64Code)
			}
			code = int(int64Code)
		}
	} else if tfield.Type.String() == "error" {
		code = w.defaultErrHttpStatusCode
	}
	parsedField := ParsedField{Field: tfield.Name}
	parsedField.Code = code

	// parse header or body info
	switch {
	case inBody:
		if len(contentType) == 0 {
			contentType = w.defaultShortContentType
		}
		parsedField.In = "body"
		parsedField.ContentType = &contentType
	case inHeader:
		// TODO: support multiple type
		if tfield.Type.String() != "*string" {
			return nil, fmt.Errorf("response object header field '%v' should be type *string", tfield.Name)
		}

		if len(headerKey) == 0 {
			headerKey = tfield.Name
		}
		parsedField.In = "header"
		parsedField.Key = &headerKey
	}

	return &parsedField, nil
}

func (w writer) doParseResp(response reflect.Type) (*ParsedResponse, error) {
	parsing := NewParsedResponse()
	for i := 0; i < response.NumField(); i++ {
		tfield := response.Field(i)

		if len(tfield.Tag) == 0 && tfield.Anonymous {
			parsedField := NewParsedResponse()
			if err := w.parseResponse(tfield.Type, parsedField); err != nil {
				return nil, err
			}
			// merge with current parsing
			for code, fields := range parsedField.CodeToFields {
				if parsing.CodeToFields[code] == nil {
					parsing.CodeToFields[code] = fields
				} else {
					parsing.CodeToFields[code] = append(parsing.CodeToFields[code], fields...)
				}
			}
			continue
		}

		parsedField, err := w.parseField(tfield)
		if err != nil {
			return nil, err
		}

		if parsing.CodeToFields[parsedField.Code] == nil {
			parsing.CodeToFields[parsedField.Code] = []ParsedField{*parsedField}
		} else {
			parsing.CodeToFields[parsedField.Code] = append(parsing.CodeToFields[parsedField.Code], *parsedField)
		}
	}

	return parsing, nil
}

func (w *writer) parseResponse(response reflect.Type, out *ParsedResponse) error {
	ref := response.PkgPath() + "." + response.Name()
	has := false

	func() {
		w.lock.RLock()
		defer w.lock.RUnlock()

		if parsed, ok := w.parsedResponse[ref]; ok {
			has = ok
			*out = parsed
		}
	}()

	if has {
		return nil
	}

	parsing, err := w.doParseResp(response)
	if err != nil {
		return err
	}

	func() {
		w.lock.Lock()
		defer w.lock.Unlock()

		if parsed, ok := w.parsedResponse[ref]; ok {
			*out = parsed
		} else {
			*out = *parsing
			w.parsedResponse[ref] = *parsing
		}
	}()

	return nil
}

type writeContext struct {
	code             int
	values           []interface{}
	fields           []ParsedField
	in               []string
	keyOrContentType []string
}

func (w writer) prepareWrite(writer http.ResponseWriter, resp reflect.Value, parsed *ParsedResponse) (*writeContext, error) {
	var (
		// first ensure only one code would be written in header
		codes = sets.NewGenericHashSet()
		// avoid iterate struct member twice
		cachedValue []interface{}
		// store only those field required to be written
		cachedParsed []ParsedField

		keyOrContentType []string
		in               []string
		// bodyCounter count how many body fields requires to be written, if multiple then return error
		bodyCounter = 0
	)

	for code, fields := range parsed.CodeToFields {
		for _, field := range fields {
			vfield, _ := reflects.GetStructMember(resp, field.Field, true)
			if !vfield.IsZero() {
				// default code -1 would not be written to httpStatusCode
				if code != -1 {
					codes.Insert(code)
				}
				cachedValue = append(cachedValue, vfield.Interface())
				cachedParsed = append(cachedParsed, field)
				switch field.In {
				case "header":
					keyOrContentType = append(keyOrContentType, *field.Key)
				case "body":
					keyOrContentType = append(keyOrContentType, *field.ContentType)
					bodyCounter++
				}
				in = append(in, field.In)
			}
		}

		if codes.Len() > 1 {
			return nil, errMultipleStatusCodeWouldBeWritten
		}

		if bodyCounter > 1 {
			return nil, errMultipleBodyWouldBeWritten
		}
	}

	if codes.Len() == 0 {
		return nil, errNoResponseToWritten
	}

	return &writeContext{
		code:   codes.UnsortedList()[0].(int),
		values: cachedValue,
		fields: cachedParsed,
		in:     in,

		keyOrContentType: keyOrContentType,
	}, nil
}

func (w writer) doWriteResponse(writer http.ResponseWriter, ctx *writeContext) error {
	for i, val := range ctx.values {
		keyOrContentType := ctx.keyOrContentType[i]
		in := ctx.in[i]
		switch in {
		case "body":
			var (
				fullContentType string
				raw             []byte
			)
			switch keyOrContentType {
			case "json":
				raw, _ = json.Marshal(val)
			case "xml":
				raw, _ = xml.Marshal(val)
			case "text":
				str := fmt.Sprintf("%v", val)
				raw = []byte(str)
			default:
				return fmt.Errorf("content-type '%v' is not supported", keyOrContentType)
			}
			fullContentType = w.shortContentTypeToFull[keyOrContentType]
			writer.Header().Set(headerContentType, fullContentType)
			writer.Write(raw)
		case "header":
			pstr := val.(*string)
			writer.Header().Set(keyOrContentType, *pstr)
		}
	}

	writer.WriteHeader(ctx.code)
	return nil
}

func (w writer) writeResponse(writer http.ResponseWriter, resp reflect.Value, parsed *ParsedResponse) error {
	writeCtx, err := w.prepareWrite(writer, resp, parsed)
	if err != nil {
		return err
	}
	return w.doWriteResponse(writer, writeCtx)
}

func (w writer) doOnFailedToWriteResponse(writer http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	w.onFailedToWriteResponse(writer, err)
}

func (w writer) Write(writer http.ResponseWriter, response interface{}) {
	var err error
	defer func() { w.doOnFailedToWriteResponse(writer, err) }()

	v := reflect.ValueOf(response)
	if err := reflects.ShouldBeKind(v.Type(), reflect.Struct); err != nil {
		return
	}
	parsed := NewParsedResponse()
	if err = w.parseResponse(v.Type(), parsed); err != nil {
		return
	}

	err = w.writeResponse(writer, v, parsed)
	return
}
