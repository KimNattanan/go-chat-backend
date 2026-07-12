package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/rs/zerolog"
)

func captureLogger(t *testing.T) (*Logger, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	z := zerolog.New(&buf).Level(zerolog.DebugLevel)
	return &Logger{logger: &z}, &buf
}

func TestError_ErrWithLabel(t *testing.T) {
	l, buf := captureLogger(t)
	err := errors.New("boom")

	l.Error(err, "rest - v1 - login")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal log: %v\nraw: %s", err, buf.String())
	}
	if entry["message"] != "rest - v1 - login" {
		t.Fatalf("message = %v, want context label", entry["message"])
	}
	if entry["error"] != "boom" {
		t.Fatalf("error = %v, want boom", entry["error"])
	}
	if entry["level"] != "error" {
		t.Fatalf("level = %v, want error", entry["level"])
	}
}

func TestError_ErrOnly(t *testing.T) {
	l, buf := captureLogger(t)

	l.Error(errors.New("only-err"))

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal log: %v\nraw: %s", err, buf.String())
	}
	if entry["error"] != "only-err" {
		t.Fatalf("error = %v, want only-err", entry["error"])
	}
	if entry["message"] != "only-err" {
		t.Fatalf("message = %v, want only-err", entry["message"])
	}
}

func TestError_StringMessage(t *testing.T) {
	l, buf := captureLogger(t)

	l.Error("plain message")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal log: %v\nraw: %s", err, buf.String())
	}
	if entry["message"] != "plain message" {
		t.Fatalf("message = %v, want plain message", entry["message"])
	}
}

func TestWarn_ErrWithLabel(t *testing.T) {
	l, buf := captureLogger(t)

	l.Warn(errors.New("bad request"), "rest - v1 - login")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal log: %v\nraw: %s", err, buf.String())
	}
	if entry["level"] != "warn" {
		t.Fatalf("level = %v, want warn", entry["level"])
	}
	if entry["message"] != "rest - v1 - login" {
		t.Fatalf("message = %v", entry["message"])
	}
}

func TestWith_AddsStructuredField(t *testing.T) {
	l, buf := captureLogger(t)

	l.With("x-request-id", "rid-1").Info("hello")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal log: %v\nraw: %s", err, buf.String())
	}
	if entry["x-request-id"] != "rid-1" {
		t.Fatalf("x-request-id = %v, want rid-1", entry["x-request-id"])
	}
}

