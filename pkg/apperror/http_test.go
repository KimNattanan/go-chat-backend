package apperror

import (
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type sampleReq struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

func TestParseHttp_ValidationErrors(t *testing.T) {
	err := validator.New().Struct(sampleReq{})
	if err == nil {
		t.Fatal("expected validation error")
	}

	code, resp := ParseHttp(err)
	if code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", code)
	}
	if resp.Message != "validation failed" {
		t.Fatalf("message = %q, want validation failed", resp.Message)
	}
	if resp.Code != CodeValidationFailed {
		t.Fatalf("code string = %q, want %s", resp.Code, CodeValidationFailed)
	}
	if resp.Errors["Email"] == "" || resp.Errors["Password"] == "" {
		t.Fatalf("errors = %#v, want Email and Password fields", resp.Errors)
	}
}

func TestParseHttp_UnauthorizedHasCode(t *testing.T) {
	code, resp := ParseHttp(Unauthorized("unauthorized", nil))
	if code != http.StatusUnauthorized {
		t.Fatalf("code = %d, want 401", code)
	}
	if resp.Code != CodeUnauthorized {
		t.Fatalf("code string = %q, want %s", resp.Code, CodeUnauthorized)
	}
	if resp.Errors != nil {
		t.Fatalf("errors = %#v, want nil", resp.Errors)
	}
}

func TestParseHttp_InternalDefault(t *testing.T) {
	code, resp := ParseHttp(assertErr("boom"))
	if code != http.StatusInternalServerError {
		t.Fatalf("code = %d, want 500", code)
	}
	if resp.Code != CodeInternal {
		t.Fatalf("code string = %q, want %s", resp.Code, CodeInternal)
	}
}

type assertErr string

func (e assertErr) Error() string { return string(e) }

func TestParseGrpc_ForeignKeyAlreadyExists(t *testing.T) {
	err := ParseGrpc(gorm.ErrForeignKeyViolated)
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected grpc status")
	}
	if st.Code() != codes.AlreadyExists {
		t.Fatalf("code = %v, want AlreadyExists", st.Code())
	}
}

func TestParseGrpc_RedisErrorInternalNoLeak(t *testing.T) {
	err := ParseGrpc(redis.ErrCrossSlot)
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected grpc status")
	}
	if st.Code() != codes.Internal {
		t.Fatalf("code = %v, want Internal", st.Code())
	}
	if st.Message() != "redis error" {
		t.Fatalf("message = %q, want redis error", st.Message())
	}
}

func TestForbidden(t *testing.T) {
	err := Forbidden("forbidden", nil)
	if err.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", err.Code)
	}
	code, resp := ParseHttp(err)
	if code != http.StatusForbidden || resp.Code != CodeForbidden {
		t.Fatalf("got status=%d code=%s", code, resp.Code)
	}
}
