package swagger

type PathItem struct {
	// Ref must be refed to another PathItem object
	Ref string `json:"$ref,omitempty"`
	// Summary an optional, string summary, intended to apply to all operations in this path.
	Summary string `json:"summary,omitempty"`
	// Description CommanMark syntax is acceptable
	Description string `json:"description,omitempty"`
	// Get get operation
	Get string `json:"get"`
}

// Operation
//  callback defined in OAS is currently not available, maybe support it in future

type Operation struct {
	// Tags logical group of operation
	Tags []string `json:"tags"`
	// Summary short description of what this operation does
	Summary string `json:"summary"`
	// ExternalDocs
	ExternalDocs ExternalDocs `json:"externalDocs"`
	// OperationId case-sensitive unique across global
	OperationId string `json:"operationId"`
	// Parameters parameters in header cookie path or query
	Parameters []ParameterOrRef `json:"parameters"`
	// RequestBody parameter in body
	RequestBody RequestBodyOrRef `json:"requestBody"`
	// Responses describe responses
	Responses Responses `json:"responses"`
	// Deprecated declares this operation should be deprecated
	Deprecated bool `json:"deprecated"`
	// Servers override default servers
	Servers []Server `json:"servers"`
	// Security as defined in oas
	Security []map[string][]string `json:"security"`
}

type ParameterOrRef struct {
	*Reference `json:",omitempty"`
	*Parameter `json:",omitempty"`
}

// Parameter object
//  content field would be supported in the future
type Parameter struct {
	// Name name of current parameter
	Name string `json:"name"`
	// In "query", "header", "path" or "cookie"
	//  if value is "path", path field must be provided
	//  if value is "header", and name is "Accept","Content-Type","Authorization" than this parameter will be ignored
	In string `json:"in"`
	// Description brief description of parameter, CommonMark syntax is acceptable
	Description string `json:"description"`
	// 	Required must be presented
	Required bool `json:"required"`
	// Deprecated as field name
	Deprecated bool `json:"deprecated"`
	// Style how value will be serialized according to different in value
	//  "query" -> "form";"path" -> "simple";"header" -> "simple";"cookie" -> "form"
	// value can be:
	//  "matrix": RFC6570
	//  "label": RFC6570
	//  "form": RFC6570
	//  "simple": RFC6570
	//  "spaceDelimited": array value would be space delimited like "ssv"
	//  "pipeDelimited": array value would be pipe separated
	//  "deepObject": object using form parameters to rendering nested objects
	//   examples can be found at "https://swagger.io/specification/#parameter-object"
	Style string `json:"style"`
	// 	Explode
	//	when style == "form", the default value is true
	Explode bool `json:"explode"`
	// Schema parameter definition
	Schema SchemaOrRef `json:"schema"`
	// Example the value must match schema
	Example interface{} `json:"example,omitempty"`
	// Examples the value must match schema field
	Examples map[string]interface{} `json:"examples,omitempty"`
}

// Encoding not supported currently
type Encoding struct {
}

type MediaType struct {
	// Schema describe content type
	Schema SchemaOrRef `json:"schema"`
	// Example value must match schema
	Example interface{} `json:"example"`
	// Examples value must match schema
	Examples map[string]interface{} `json:"examples"`
	// Encoding map between property name and its encoding information
	Encoding map[string]Encoding `json:"encoding"`
}

type RequestBodyOrRef struct {
	*Reference   `json:",omitempty"`
	*RequestBody `json:",omitempty"`
}

type RequestBody struct {
	// Description CommonMark syntax can be used
	Description string `json:"description"`
	// Content content of request body
	Content map[string]MediaType `json:"content"`
	// Required tells whether this request body must be presented
	Required bool `json:"required"`
}

// Header
//  not exactly as defined in OAS, just assume the value to be string or integer
type Header struct {
	// Description brief description of parameter, CommonMark syntax is acceptable
	Description string `json:"description"`
	// 	Required must be presented
	Required bool `json:"required"`
	// Deprecated as field name
	Deprecated bool `json:"deprecated"`
	// Schema parameter definition
	Schema SchemaOrRef `json:"schema"`
	// Example the value must match schema
	Example interface{} `json:"example,omitempty"`
}

// Response
//  link object defined in OAS is omitted, maybe support it in future
type Response struct {
	// Description CommonMark syntax can be used
	Description string `json:"description"`
	// Headers in response header
	Headers map[string]Header `json:"headers"`
	// Content describe the potential response payload, key should be media type or media type range
	Content map[string]MediaType `json:"content"`
}

// Responses
//  fixed key would be "default",other key should be http status code like "200"
type Responses map[string]Response
