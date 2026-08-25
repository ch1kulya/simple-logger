package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMiddlewareAllMethods(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(LevelDebug)
	defer SetLevel(LevelDebug)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case http.MethodPut, http.MethodPatch:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusOK)
		}
		w.Write([]byte("response"))
	})

	wrapped := Middleware(handler)

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
	}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/api/"+strings.ToLower(method), nil)
		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)
	}

	out := buf.String()
	t.Log("\n" + out)

	for _, method := range methods {
		if !strings.Contains(out, method) {
			t.Errorf("expected %s in output", method)
		}
	}
}

func TestResponseWriterWrapper(t *testing.T) {
	rec := httptest.NewRecorder()
	wrapper := &responseWriterWrapper{
		ResponseWriter: rec,
		status:         http.StatusOK,
	}

	wrapper.WriteHeader(http.StatusNotFound)
	wrapper.Write([]byte("not found"))
	wrapper.WriteHeader(http.StatusOK)

	if wrapper.status != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", wrapper.status)
	}
	if wrapper.size != 9 {
		t.Errorf("expected size 9, got %d", wrapper.size)
	}
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected underlying status 404, got %d", rec.Code)
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes    int
		expected string
	}{
		{0, "0B"},
		{512, "512B"},
		{1024, "1.0K"},
		{1536, "1.5K"},
		{1048576, "1.0M"},
		{1572864, "1.5M"},
	}

	for _, tt := range tests {
		got := formatSize(tt.bytes)
		if got != tt.expected {
			t.Errorf("formatSize(%d) = %s, want %s", tt.bytes, got, tt.expected)
		}
	}
}

func TestMiddlewareLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(LevelWarn)
	defer SetLevel(LevelDebug)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := Middleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if buf.Len() > 0 {
		t.Error("middleware logs should be filtered at WARN level")
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		d        time.Duration
		expected string
	}{
		{0, "<1µs"},
		{500 * time.Nanosecond, "<1µs"},
		{1 * time.Microsecond, "1µs"},
		{50 * time.Microsecond, "50µs"},
		{1500 * time.Microsecond, "1.5ms"},
		{50 * time.Millisecond, "50ms"},
	}

	for _, tt := range tests {
		got := formatDuration(tt.d)
		if got != tt.expected {
			t.Errorf("formatDuration(%v) = %s, want %s", tt.d, got, tt.expected)
		}
	}
}

func TestMiddlewareLatency(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(LevelDebug)
	defer SetLevel(LevelDebug)

	delays := []time.Duration{
		5 * time.Millisecond,
		50 * time.Millisecond,
	}

	for _, delay := range delays {
		d := delay
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(d)
			w.WriteHeader(http.StatusOK)
		})
		wrapped := Middleware(handler)
		req := httptest.NewRequest(http.MethodGet, "/latency", nil)
		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)
	}

	out := buf.String()
	t.Log("\n" + out)

	if !strings.Contains(out, "ms") {
		t.Error("expected millisecond durations in output")
	}
}

func TestMiddlewareUserAgent(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(LevelDebug)
	defer SetLevel(LevelDebug)

	SetShowUserAgent(true)
	defer SetShowUserAgent(false)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := Middleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/test-ua", nil)
	req.Header.Set("User-Agent", "TestAgent/1.0")
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	out := buf.String()
	t.Log("\n" + out)

	if !strings.Contains(out, "ua:") || !strings.Contains(out, "TestAgent/1.0") {
		t.Error("expected user-agent in output")
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 {
		t.Errorf("expected at least 2 lines of output, got %d", len(lines))
	}
}
