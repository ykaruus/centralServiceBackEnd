package infra

import "context"

type Generator interface {
	New() string
	CtxWithValue(context.Context, string) context.Context
	GetTraceIDFromCtx(context.Context) string
	GetTraceIDKey() any
}
