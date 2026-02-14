package middlewares

import (
	"centralService/internal/infra"

	"github.com/gin-gonic/gin"
)

func TraceIDMiddleware(uuuid infra.Generator) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		traceID := ctx.Request.Header.Get("X-Trace-ID")

		if traceID == "" {
			traceID = uuuid.New()
		}
		ctx.Set(uuuid.GetTraceIDKey(), traceID)

		ctx.Request.Header.Set("X-Trace-ID", traceID)
		ctx.Next()
	}
}
