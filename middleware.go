package logger

import (
	"fmt"
	"net/http"
	"time"

	"charm.land/lipgloss/v2"
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

func getMethodStyle(method string) lipgloss.Style {
	switch method {
	case http.MethodGet:
		return greenStyle
	case http.MethodPost:
		return cyanStyle
	case http.MethodPut:
		return yellowStyle
	case http.MethodPatch:
		return magentaStyle
	case http.MethodDelete:
		return redStyle
	case http.MethodOptions, http.MethodHead:
		return grayStyle
	default:
		return whiteStyle
	}
}

func getStatusStyle(status int) lipgloss.Style {
	switch {
	case status >= 500:
		return redStyle
	case status >= 400:
		return yellowStyle
	case status >= 300:
		return cyanStyle
	default:
		return greenStyle
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

func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return "<1µs"
	}
	return d.Round(time.Microsecond).String()
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &responseWriterWrapper{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		if getLevel() > LevelInfo {
			return
		}

		methodStyle := getMethodStyle(r.Method)
		statusStyle := getStatusStyle(wrapper.status)
		duration := time.Since(start)

		line := fmt.Sprintf("%s %s %-7s %s %s %s",
			getTimestamp(),
			methodStyle.Width(7).Render(r.Method),
			formatSize(wrapper.size),
			statusStyle.Width(3).Render(fmt.Sprintf("%d", wrapper.status)),
			r.URL.Path,
			dimStyle.Render(formatDuration(duration)),
		)

		var fields []field
		if getShowUserAgent() && r.UserAgent() != "" {
			fields = append(fields, field{key: "ua", value: r.UserAgent()})
		}

		if meta := renderMetadata(fields); meta != "" {
			line += "    " + meta
		}

		fmt.Fprintln(output, line)
	})
}
