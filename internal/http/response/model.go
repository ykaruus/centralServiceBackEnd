package response

type ErrorBody struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Fields  []string `json:"fields"`
}

type Response struct {
	Success   bool       `json:"success"`
	Data      any        `json:"data,omitempty"`
	Error     *ErrorBody `json:"error,omitempty"`
	Timestamp int64      `json:"timestamp"`
	TraceID   string     `json:"trace_id"`
}
