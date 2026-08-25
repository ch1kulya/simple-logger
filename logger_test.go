package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogFunctions(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(LevelDebug)
	defer SetLevel(LevelDebug)

	Debug("debug message %d", 1)
	Info("info message %s", "test")
	Warn("warn message")
	Error("error message")

	out := buf.String()
	t.Log("\n" + out)

	if !strings.Contains(out, "DEBUG") {
		t.Error("expected DEBUG in output")
	}
	if !strings.Contains(out, "INFO") {
		t.Error("expected INFO in output")
	}
	if !strings.Contains(out, "WARN") {
		t.Error("expected WARN in output")
	}
	if !strings.Contains(out, "ERR") {
		t.Error("expected ERR in output")
	}
}

func TestLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(LevelWarn)
	defer SetLevel(LevelDebug)

	Debug("should not appear")
	Info("should not appear")
	Warn("should appear")
	Error("should appear")

	out := buf.String()
	t.Log("\n" + out)

	if strings.Contains(out, "DEBUG") {
		t.Error("DEBUG should be filtered")
	}
	if strings.Contains(out, "INFO") {
		t.Error("INFO should be filtered")
	}
	if !strings.Contains(out, "WARN") {
		t.Error("WARN should appear")
	}
	if !strings.Contains(out, "ERR") {
		t.Error("ERR should appear")
	}
}

func TestFatal(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	var exitCode int
	originalExit := exitFunc
	exitFunc = func(code int) { exitCode = code }
	defer func() { exitFunc = originalExit }()

	Fatal("fatal error %s", "test")

	out := buf.String()
	t.Log("\n" + out)

	if !strings.Contains(out, "FATAL") {
		t.Error("expected FATAL in output")
	}
	if !strings.Contains(out, "fatal error test") {
		t.Error("expected message in output")
	}
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
}

func TestWithMetadata(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(LevelDebug)
	defer SetLevel(LevelDebug)

	With(map[string]any{"user_id": 123, "role": "admin"}).Info("user action")
	With(map[string]any{"error_code": 404}).Error("not found")
	With(map[string]any{"key": "val"}).Warn("warning")
	With(map[string]any{"debug_info": "trace"}).Debug("debugging")

	out := buf.String()
	t.Log("\n" + out)

	if !strings.Contains(out, "user_id:") || !strings.Contains(out, "123") {
		t.Error("expected user_id in output")
	}
	if !strings.Contains(out, "role:") || !strings.Contains(out, "admin") {
		t.Error("expected role in output")
	}
	if !strings.Contains(out, "error_code:") || !strings.Contains(out, "404") {
		t.Error("expected error_code in output")
	}
}

func TestWithChaining(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(LevelDebug)
	defer SetLevel(LevelDebug)

	entry := With(map[string]any{"service": "auth"})
	entry.With(map[string]any{"ip": "127.0.0.1"}).Info("login attempt")

	out := buf.String()
	t.Log("\n" + out)

	if !strings.Contains(out, "service:") || !strings.Contains(out, "auth") {
		t.Error("expected service in output")
	}
	if !strings.Contains(out, "ip:") || !strings.Contains(out, "127.0.0.1") {
		t.Error("expected ip in output")
	}
}

func TestWithFatal(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	var exitCode int
	originalExit := exitFunc
	exitFunc = func(code int) { exitCode = code }
	defer func() { exitFunc = originalExit }()

	With(map[string]any{"fatal_reason": "db_down"}).Fatal("system crash")

	out := buf.String()
	t.Log("\n" + out)

	if !strings.Contains(out, "FATAL") {
		t.Error("expected FATAL in output")
	}
	if !strings.Contains(out, "fatal_reason:") || !strings.Contains(out, "db_down") {
		t.Error("expected fatal_reason in output")
	}
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
}
