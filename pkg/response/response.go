package response

type ParsedField struct {
	// Field
	Field string
	// In value should be "header" "body"
	In string
	// Key If In is "header" then key must be presented
	Key *string
	// ContentType if In is "body" then content type should be presented
	ContentType *string
	// Code corresponding response statusCode
	Code int
}

// ParsedResponse parsed responseObject, code should be httpStatusCode
type ParsedResponse struct {
	CodeToFields map[int][]ParsedField
}

func NewParsedResponse() *ParsedResponse {
	return &ParsedResponse{CodeToFields: map[int][]ParsedField{}}
}
