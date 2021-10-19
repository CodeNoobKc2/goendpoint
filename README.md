# goendpoint
This is a go http endpoint library which aims to ease the developer experience for gopher developing with http endpoint.
This project is currently unstable and under development
# Content
- [GoEndpoint](#goendpoint)
    - [RoadMap](#Roadmap)
    - [Examples](#Examples)
      - [Auto Binding&Auto Wirting](##AutoBinding&AutoWriting)
# RoadMap
*phase one would only focus on text/plain application/xml and application/json type since these are the most commonly used.*
- [x] auto-binding request object with tags (ph1)
- [x] auto-writing http response with tags (ph1)
- [ ] support auto validation with validator library (ph1)
- [ ] support auto swagger doc generation without extra annotation (ph1)
- [ ] support auto client generator without command line util (ph1)
- [ ] benchmarking and enhance performance (ph1)
- [ ] support other kinds of content-type such as text/html or multipart
# Examples
## AutoBinding&AutoWriting
```go
// the runnable example can be found at pkg example/auto_binding_writing
// Payload serve as request body or response body
type Payload struct {
	Hello string `json:"hello"`
}

// RequestObject accept input through header and body
type RequestObject struct {
	// RequestId would be passed by header
	RequestId *string `header:"X-Request-Id"`
	// Payload would be passed by body
	Payload   `body:"json"`
}

// ResponseObject return application/text->errorInfo or application/json->echoMessage
type ResponseObject struct {
	// Payload if presented return application/json->echoMessage
	*Payload `body:"json"`
	// Error if presented return application/text->errorInfo
	Error    error `body:"text"`
}

// Echo takes in filled requestObject and return corresponding response object
func Echo(ctx context.Context, req RequestObject) (resp ResponseObject) {
	if req.RequestId == nil {
		resp.Error = errors.New("invalid request")
		return
	}

	resp.Payload = &Payload{Hello: fmt.Sprintf("%v : %v", *req.RequestId, req.Hello)}
	return
}
```