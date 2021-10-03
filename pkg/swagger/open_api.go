package swagger

type Reference struct {
	Ref string `json:"$ref"`
}

type OpenApi struct {
	// OpenApiVersion semantic version number of open api specification
	OpenApiVersion string `json:"openapi"`
	// Info info object
	Info Info `json:"info"`
	// Servers servers object
	Servers []Server `json:"servers,omitempty"`
	// Paths key should follow the pattern /{path}
	Paths map[string]PathItem `json:"paths"`
	// Tags
	Tags []Tag `json:"tags"`
	// Components a collection of reusable objects
	Components *Components `json:"components,omitempty"`
	// Security declaration of which security mechanisms can be used across the API
	Security []map[string][]string `json:"security,omitempty"`
	// ExternalDocs add external docs
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

// Components
// a collection of reuseable objects
// examples,headers,links,callbacks fields declared in OAS is not supported under early development
type Components struct {
	Schemas         map[string]SchemaOrRef    `json:"schemas,omitempty"`
	Responses       map[string]ResponseOrRef  `json:"responses,omitempty"`
	Parameters      map[string]ParameterOrRef `json:"parameters,omitempty"`
	RequestBodies   map[string]RequestBody    `json:"requestBodies,omitempty"`
	SecuritySchemes map[string]SecuritySchema `json:"securitySchemes,omitempty"`
}

type ExternalDocs struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}

type Contact struct {
	// Name name of contact
	Name string `json:"name"`
	// Url contact detail url
	Url string `json:"url"`
	// Email contact email
	Email string `json:"email"`
}

type License struct {
	// Name license name
	Name string `json:"name"`
	// Url license detail url
	Url string `json:"url"`
}

type Server struct {
	// Url a url to target host
	Url string `json:"url"`
	// Description describe of this server
	Description string `json:"description"`
	// Variables server variables
	Variables map[string]ServerVariable `json:"variables"`
}

type ServerVariable struct {
	Enum        []interface{} `json:"enum,omitempty"`
	Default     *string       `json:"default,omitempty"`
	Description *string       `json:"description,omitempty"`
}

// Info provides metadata about certain object
type Info struct {
	// Title
	Title string `json:"title"`
	// Description short description of api
	Description string `json:"description"`
	// TermsOfService url of terms of service
	TermsOfService string `json:"termsOfService"`
	// License license object
	License License `json:"license"`
	// Contact contact information for the api
	Contact Contact `json:"contact"`
	// Version the version of this info object
	Version string `json:"version"`
}

type Tag struct {
	// Name tag name
	Name string `json:"name"`
	// Description short description of tag
	Description string `json:"description"`
	// ExternalDocs additional external documentation
	ExternalDocs ExternalDocs `json:"externalDocs"`
}

type SecuritySchema struct {
	// Type
	// value can only be "apiKey" "http" "oath2" "openIdConnect"
	Type string `json:"type"`
	// Description a short description for this schema
	Description string `json:"description,omitempty"`
	// Name applied to "apiKey", the name of the header,query or cookie
	Name string `json:"name,omitempty"`
	// In applied to "apiKey", value should be "query" "header" or "cookie"
	In string `json:"in,omitempty"`
	// Schema applied to "http", authorization header as defined in RFC7235
	Schema string `json:"schema,omitempty"`
	// BearerFormat applied to "http" generally for documentation purpose
	BearerFormat string `json:"bearerFormat,omitempty"`
	// Flows applied to "oauth2"
	Flows *OAuthFlows `json:"flows,omitempty"`
	// OpenIdConnectUrl
	OpenIdConnectUrl string `json:"openIdConnectUrl,omitempty"`
}

type OAuthFlowsOrRef struct {
	*Reference  `json:",omitempty"`
	*OAuthFlows `json:",omitempty"`
}

type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

type OAuthFlow struct {
	// AuthorizationUrl applies to "oauth2" ("implicit","authorizationCode")
	AuthorizationUrl string `json:"authorizationUrl,omitempty"`
	// TokenUrl applies to "oauth2" ("password","clientCredentials","authorizationCode")
	TokenUrl string `json:"tokenUrl,omitempty"`
	// RefreshUrl used to obtaining refresh tokens
	RefreshUrl string `json:"refreshUrl,omitempty"`
	// Scopes the available scopes for oauth2 security
	Scopes map[string]string `json:"scopes,omitempty"`
}
