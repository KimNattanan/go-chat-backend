package requestid

import (
	"context"

	"google.golang.org/grpc"
)

// UnaryClientInterceptor forwards the request ID from ctx into outgoing gRPC metadata.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = AppendToOutgoingContext(ctx)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
