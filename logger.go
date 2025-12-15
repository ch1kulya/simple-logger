package logger

import (
	"fmt"
	"os"
	"time"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
)

func getTimestamp() string {
	return fmt.Sprintf("%s[%s]%s", Gray, time.Now().Format("15:04:05"), Reset)
}

func Info(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	fmt.Printf("%s %sINFO%s  %s\n", getTimestamp(), Green, Reset, msg)
}

func Error(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintf(os.Stderr, "%s %sERR%s   %s\n", getTimestamp(), Red, Reset, msg)
}

func Warn(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	fmt.Printf("%s %sWARN%s  %s\n", getTimestamp(), Yellow, Reset, msg)
}

func Debug(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	fmt.Printf("%s %sDEBUG%s %s\n", getTimestamp(), Blue, Reset, msg)
}
