package grpcserver

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/pkg/requestid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryServerRequestID ensures every RPC has a request ID on context.
// It reuses incoming x-request-id metadata when present, otherwise generates one,
// and mirrors it onto the response headers.
func UnaryServerRequestID() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		rid := requestid.FromIncomingMetadata(ctx)
		if rid == "" {
			rid = requestid.New()
		}
		ctx = requestid.WithContext(ctx, rid)
		_ = grpc.SetHeader(ctx, metadata.Pairs(requestid.MetadataKey, rid))
		return handler(ctx, req)
	}
}
