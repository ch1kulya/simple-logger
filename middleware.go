package logger

import (
	"fmt"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	size        int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

func getMethodColor(method string) string {
	switch method {
	case http.MethodGet:
		return Green
	case http.MethodPost:
		return Cyan
	case http.MethodPut:
		return Yellow
	case http.MethodPatch:
		return Magenta
	case http.MethodDelete:
		return Red
	case http.MethodOptions, http.MethodHead:
		return Gray
	default:
		return White
	}
}

func getStatusColor(status int) string {
	switch {
	case status >= 500:
		return Red
	case status >= 400:
		return Yellow
	case status >= 300:
		return Cyan
	default:
		return Green
	}
}

func formatSize(bytes int) string {
	switch {
	case bytes >= 1024*1024:
		return fmt.Sprintf("%.1fM", float64(bytes)/(1024*1024))
	case bytes >= 1024:
		return fmt.Sprintf("%.1fK", float64(bytes)/1024)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &responseWriterWrapper{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		methodColor := getMethodColor(r.Method)
		statusColor := getStatusColor(wrapper.status)
		duration := time.Since(start)

		fmt.Printf("%s %s%-7s%s %-7s %s%-3d%s %s %s%s%s\n",
			getTimestamp(),
			methodColor, r.Method, Reset,
			formatSize(wrapper.size),
			statusColor, wrapper.status, Reset,
			r.URL.Path,
			Dim, duration.Round(time.Microsecond), Reset,
		)
	})
}
