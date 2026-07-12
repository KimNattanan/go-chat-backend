package requestid

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestWithContextRoundTrip(t *testing.T) {
	ctx := WithContext(context.Background(), "abc-123")
	if got := FromContext(ctx); got != "abc-123" {
		t.Fatalf("FromContext = %q, want abc-123", got)
	}
}

func TestFromIncomingMetadata(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(MetadataKey, "rid-1"))
	if got := FromIncomingMetadata(ctx); got != "rid-1" {
		t.Fatalf("FromIncomingMetadata = %q, want rid-1", got)
	}
}

func TestAppendToOutgoingContext(t *testing.T) {
	ctx := WithContext(context.Background(), "rid-2")
	ctx = AppendToOutgoingContext(ctx)
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("expected outgoing metadata")
	}
	if got := md.Get(MetadataKey); len(got) != 1 || got[0] != "rid-2" {
		t.Fatalf("metadata = %v, want [rid-2]", got)
	}
}
