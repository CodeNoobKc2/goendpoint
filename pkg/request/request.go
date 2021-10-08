package request

type BindContext struct {
	// In value should be "header" "query" "body" "path"
	In string
	// Key related unique key to different in type
	Key string
	// Style how to parse object or array
	Style string
	// Required is current field required
	Required bool
}

type ParsedRequestObject struct {
	// FieldToBindContext map struct field name to BindContext
	FieldToBindContext map[string]BindContext
}

func NewParsedRequestObject() *ParsedRequestObject {
	return &ParsedRequestObject{FieldToBindContext: map[string]BindContext{}}
}
