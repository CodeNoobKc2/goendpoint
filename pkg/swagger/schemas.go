package swagger

// Schema oas schema object
//  extension fields are currently not supported
//  discriminator not supported
type Schema struct {
	// Schema
	// +readonly value will always been "https://json-schema.org/draft/2020-12/schema"
	Schema string `json:"$schema"`
	// Id reference to defination
	Id string `json:"$id,omitempty"`
	// Title the title of current schema
	Title string `json:"title"`
	// Description description of current schema
	Description string `json:"description"`
	// Nullable can be null
	Nullable bool
	// Deprecated alert user that this schema is deprecated
	Deprecated bool
	// Example example value for current schema
	Example  interface{} `json:"example"`
	*Integer `json:",omitempty"`
	*Number  `json:",omitempty"`
	*Object  `json:",omitempty"`
	*String  `json:",omitempty"`
	*Array   `json:",,omitempty"`
	// Enum valid values
	Enum []interface{} `json:"enum,omitempty"`
	// Const const value
	Const interface{} `json:"const,omitempty"`
	// AnyOf
	AnyOf []*Schema `json:"anyOf,omitempty"`
	// OneOf
	OneOf []*Schema `json:"oneOf,omitempty"`
}

type Null struct {
	// Type
	// +readonly will always be "null"
	Type string `json:"type"`
}

type Integer struct {
	// Type
	// +readonly will always be "integer"
	Type string `json:"type"`
	// MultipleOf value mod <multipleOf> = 0
	MultipleOf *int `json:"multipleOf"`
	// Format would be oneOf int32,int64 or omit
	Format *string `json:"format,omitempty"`
	// Minimum >=
	Minimum *int `json:"minimum,omitempty"`
	// ExclusiveMinimum >
	ExclusiveMinimum *int `json:"exclusiveMinimum,omitempty"`
	// Maximum <=
	Maximum *int `json:"maximum,omitempty"`
	// ExclusiveMaximum <
	ExclusiveMaximum *int `json:"exclusiveMaximum,omitempty"`
}

type Number struct {
	// Type
	// +readonly will always be "number"
	Type string `json:"type"`
}

type Boolean struct {
	// Type
	// +readonly will always be "boolean"
	Type string `json:"type"`
	// MultipleOf value mod <multipleOf> = 0
	MultipleOf *int `json:"multipleOf"`
	// Format would be oneOf double,float or omit
	Format *string `json:"format,omitempty"`
	// Minimum >=
	Minimum *int `json:"minimum,omitempty"`
	// ExclusiveMinimum >
	ExclusiveMinimum *int `json:"exclusiveMinimum,omitempty"`
	// Maximum <=
	Maximum *int `json:"maximum,omitempty"`
	// ExclusiveMaximum <
	ExclusiveMaximum *int `json:"exclusiveMaximum,omitempty"`
}

type Items struct {
	Type Schema `json:"type"`
}

type PrefixItem struct {
	Type *Schema `json:"type,omitempty"`
}

type Array struct {
	// Type
	// +readonly will always be "type"
	Type string `json:"type"`
	// Items
	Items *Items `json:"items,omitempty"`
	// PrefixItems tuple validation
	PrefixItems []PrefixItem `json:"prefixItems,omitempty"`
	// MinItems min length
	MinItems *int `json:"minItems,omitempty"`
	// MaxItems
	MaxItems *int `json:"maxItems,omitempty"`
	// UniqueItems
	UniqueItems bool `json:"uniqueItems,omitempty"`
}

type String struct {
	// Type
	// +readonly would always be "string"
	Type string `json:"type"`
	// MinLength minLength of current string
	MinLength *int `json:"minLength,omitempty"`
	// MaxLength maxLength of current string
	MaxLength *int `json:"maxLength,omitempty"`
	// Pattern regexp expr
	Pattern *string `json:"pattern"`
	// Format
	// date,date-time,password,byte,binary
	Format *string `json:"format"`
}

type Property struct {
	// Description description of current property
	Description string `json:"description"`
	// Schema embedded type info
	*Schema `json:",omitempty"`
}

type Object struct {
	// Type
	// +readonly will always be "object"
	Type string `json:"type"`
	// Properties
	Properties map[string]Property `json:"properties,omitempty"`
	// PatternProperties key should be a regular expression
	PatternProperties map[string]Property `json:"patternProperties,omitempty"`
	// Required required field in properties
	Required []string `json:"required"`
	// AdditionalProperties if non-nil value must be false
	AdditionalProperties *bool `json:"additionalProperties,omitempty"`
}
