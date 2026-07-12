package requestid

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

const (
	// MetadataKey is the gRPC metadata / HTTP header key (lowercase for gRPC).
	MetadataKey = "x-request-id"
	// EchoContextKey is the key used with echo.Context.Set/Get.
	EchoContextKey = "requestID"
)

type ctxKey struct{}

// New generates a request ID.
func New() string {
	return uuid.NewString()
}

// FromContext returns the request ID stored on ctx, or "".
func FromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	rid, _ := ctx.Value(ctxKey{}).(string)
	return rid
}

// WithContext stores rid on ctx.
func WithContext(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, ctxKey{}, rid)
}

// FromIncomingMetadata reads x-request-id from gRPC incoming metadata.
func FromIncomingMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	vals := md.Get(MetadataKey)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

// AppendToOutgoingContext adds x-request-id to outgoing gRPC metadata when present on ctx.
func AppendToOutgoingContext(ctx context.Context) context.Context {
	rid := FromContext(ctx)
	if rid == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, MetadataKey, rid)
}
