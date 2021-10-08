package testdata

type RequestBody struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

type InHeader struct {
	HeaderStr         string  `header:"Str"`
	HeaderOptionalStr *string `header:"Optional-Str"`
	HeaderInt         int     `header:"Int"`
}

type InQuery struct {
	QueryStr         string   `query:"str"`
	QueryInt         int      `query:"int"`
	QueryOptionalStr *string  `query:"optionalStr"`
	QueryOptionalInt *int     `query:"optionalInt"`
	QueryStrSlice    []string `query:"strSlice"`
}

type InPath struct {
	PathStr string `path:"str"`
	PathInt int    `path:"int"`
}

type BindObject struct {
	InHeader
	InQuery
	InPath
	RequestBody `body:"json"`
}
