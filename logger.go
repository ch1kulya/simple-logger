package logger

import (
	"fmt"
	"os"
	"time"
)

const (
	Reset = "\033[0m"
	Dim   = "\033[2m"

	Red     = "\033[91m"
	Green   = "\033[92m"
	Yellow  = "\033[93m"
	Blue    = "\033[94m"
	Magenta = "\033[95m"
	Cyan    = "\033[96m"
	White   = "\033[97m"
	Gray    = "\033[90m"
)

func getTimestamp() string {
	return fmt.Sprintf("%s%s%s", Dim, time.Now().Format("15:04:05"), Reset)
}

func Info(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	fmt.Printf("%s %s%-7s%s %s\n", getTimestamp(), Magenta, "INFO", Reset, msg)
}

func Error(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintf(os.Stderr, "%s %s%-7s%s %s\n", getTimestamp(), Red, "ERR", Reset, msg)
}

func Warn(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	fmt.Printf("%s %s%-7s%s %s\n", getTimestamp(), Yellow, "WARN", Reset, msg)
}

func Debug(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	fmt.Printf("%s %s%-7s%s %s\n", getTimestamp(), Blue, "DEBUG", Reset, msg)
}
