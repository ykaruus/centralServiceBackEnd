package response

import (
	"centralService/internal/infra/trace"
	"time"

	"github.com/gin-gonic/gin"
)

func OK(ctx *gin.Context, data any) *Response {
	return &Response{
		Success:   true,
		Data:      data,
		Timestamp: time.Now().Unix(),
		TraceID:   trace.NewGenerator().GetTraceIDFromCtx(ctx),
	}
}

func Fail(ctx *gin.Context, errors *ErrorBody) *Response {
	return &Response{
		Success:   false,
		Error:     errors,
		Timestamp: time.Now().Unix(),
		TraceID:   trace.NewGenerator().GetTraceIDFromCtx(ctx),
	}
}

func BadRequest(ctx *gin.Context) *Response {
	return &Response{
		Success: false,
		Error: &ErrorBody{
			Code:    ErrBadRequest.Code,
			Message: ErrBadRequest.Message,
		},
		Timestamp: time.Now().Unix(),
		TraceID:   trace.NewGenerator().GetTraceIDFromCtx(ctx),
	}
}
