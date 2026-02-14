package trace

import (
	"centralService/internal/infra"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type UUIDGen struct{}

type traceIDKey struct{}

/*


type Generator interface {
	New() string
	CtxWithValue(context.Context, string) context.Context
	GetTraceIDKey() any
	GetTraceIDFromCtx(context.Context) string
}


*/

func NewGenerator() infra.Generator {
	return UUIDGen{}
}

func (U UUIDGen) New() string {
	return uuid.NewString()
}

func (u UUIDGen) GetTraceIDKey() any {
	return traceIDKey{}
}

func (U UUIDGen) CtxWithValue(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func (U UUIDGen) GetTraceIDFromCtx(ctx context.Context) string {

	gen := UUIDGen{}

	if v := ctx.Value(traceIDKey{}); v != nil {
		if s, ok := v.(string); ok {
			{
				return s
			}
		}
	}

	return gen.New()
}

func LogWithTraceID(layer string, ctx context.Context) *slog.Logger {
	traceID := NewGenerator().GetTraceIDFromCtx(ctx)
	log := slog.With("layer", layer, "traceID", traceID)

	return log
}
