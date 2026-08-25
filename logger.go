package logger

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var (
	level         = LevelDebug
	showUserAgent = false
	mu            sync.RWMutex
	output        io.Writer = os.Stdout
	errOut        io.Writer = os.Stderr
	exitFunc                = os.Exit
)

var (
	dimStyle     = lipgloss.NewStyle().Faint(true)
	debugStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	fatalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	greenStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	cyanStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	yellowStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	magentaStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	redStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	grayStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	whiteStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
)

func SetLevel(l Level) {
	mu.Lock()
	defer mu.Unlock()
	level = l
}

func SetOutput(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	output = w
	errOut = w
}

func SetShowUserAgent(show bool) {
	mu.Lock()
	defer mu.Unlock()
	showUserAgent = show
}

func getLevel() Level {
	mu.RLock()
	defer mu.RUnlock()
	return level
}

func getShowUserAgent() bool {
	mu.RLock()
	defer mu.RUnlock()
	return showUserAgent
}

func getTimestamp() string {
	return dimStyle.Render(time.Now().Format("02 Jan 15:04:05"))
}

type field struct {
	key   string
	value any
}

type Entry struct {
	fields []field
}

func (e *Entry) With(meta map[string]any) *Entry {
	newFields := make([]field, len(e.fields))
	copy(newFields, e.fields)
	keys := make([]string, 0, len(meta))
	for k := range meta {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		newFields = append(newFields, field{key: k, value: meta[k]})
	}
	return &Entry{fields: newFields}
}

func (e *Entry) Debug(format string, v ...any) {
	logMessage(output, LevelDebug, debugStyle, "DEBUG", format, e.fields, v...)
}

func (e *Entry) Info(format string, v ...any) {
	logMessage(output, LevelInfo, infoStyle, "INFO", format, e.fields, v...)
}

func (e *Entry) Warn(format string, v ...any) {
	logMessage(output, LevelWarn, warnStyle, "WARN", format, e.fields, v...)
}

func (e *Entry) Error(format string, v ...any) {
	logMessage(errOut, LevelError, errorStyle, "ERR", format, e.fields, v...)
}

func (e *Entry) Fatal(format string, v ...any) {
	logMessage(errOut, LevelFatal, fatalStyle, "FATAL", format, e.fields, v...)
	exitFunc(1)
}

func With(meta map[string]any) *Entry {
	return (&Entry{}).With(meta)
}

func renderMetadata(fields []field) string {
	if len(fields) == 0 {
		return ""
	}
	maxKeyLen := 0
	for _, f := range fields {
		if l := lipgloss.Width(f.key); l > maxKeyLen {
			maxKeyLen = l
		}
	}

	indent := strings.Repeat(" ", lipgloss.Width(getTimestamp())+1)
	var lines []string
	for _, f := range fields {
		k := dimStyle.Width(maxKeyLen + 1).Render(f.key + ":")
		valStr := dimStyle.Render(fmt.Sprintf("%v", f.value))
		lines = append(lines, indent+k+" "+valStr)
	}
	return strings.Join(lines, "\n")
}

func logMessage(w io.Writer, lvl Level, badge lipgloss.Style, label, format string, fields []field, v ...any) {
	if getLevel() > lvl {
		return
	}
	msg := fmt.Sprintf(format, v...)
	line := fmt.Sprintf("%s %s %s", getTimestamp(), badge.Width(7).Render(label), msg)
	if meta := renderMetadata(fields); meta != "" {
		line += "\n" + meta
	}
	fmt.Fprintln(w, line)
}

func Info(format string, v ...any) {
	logMessage(output, LevelInfo, infoStyle, "INFO", format, nil, v...)
}

func Error(format string, v ...any) {
	logMessage(errOut, LevelError, errorStyle, "ERR", format, nil, v...)
}

func Warn(format string, v ...any) {
	logMessage(output, LevelWarn, warnStyle, "WARN", format, nil, v...)
}

func Debug(format string, v ...any) {
	logMessage(output, LevelDebug, debugStyle, "DEBUG", format, nil, v...)
}

func Fatal(format string, v ...any) {
	logMessage(errOut, LevelFatal, fatalStyle, "FATAL", format, nil, v...)
	exitFunc(1)
}

func SetForceColor(force bool) {
	if force {
		lipgloss.Writer.Profile = colorprofile.TrueColor
	} else {
		lipgloss.Writer.Profile = colorprofile.Ascii
	}
}
