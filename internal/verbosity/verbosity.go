package verbosity

import (
	"fmt"
	"os"
	"time"
)

type Level int

const (
	Normal Level = iota
	InfoLevel
	DebugLevel
	TraceLevel
)

var currentLevel Level = Normal

func SetLevel(level Level) {
	currentLevel = level
}

func GetLevel() Level {
	return currentLevel
}

func SetFromCount(count int) {
	if count < 0 {
		count = 0
	}
	if count > 3 {
		count = 3
	}
	currentLevel = Level(count)
}

func IsEnabled(level Level) bool {
	return currentLevel >= level
}

func Print(level Level, format string, args ...interface{}) {
	if !IsEnabled(level) {
		return
	}

	var prefix string
	switch level {
	case Normal:
		prefix = ""
	case InfoLevel:
		prefix = "ℹ️  "
	case DebugLevel:
		prefix = "🐛 [DEBUG] "
	case TraceLevel:
		prefix = "🔍 [TRACE] "
	}

	message := fmt.Sprintf(format, args...)
	if prefix != "" {
		fmt.Fprintf(os.Stderr, "%s%s\n", prefix, message)
	} else {
		fmt.Println(message)
	}
}

func Info(format string, args ...interface{}) {
	Print(InfoLevel, format, args...)
}

func Debug(format string, args ...interface{}) {
	Print(DebugLevel, format, args...)
}

func Trace(format string, args ...interface{}) {
	Print(TraceLevel, format, args...)
}

func Printf(format string, args ...interface{}) {
	Print(Normal, format, args...)
}

func PrintWithTiming(level Level, startTime time.Time, format string, args ...interface{}) {
	if !IsEnabled(level) {
		return
	}
	elapsed := time.Since(startTime)
	message := fmt.Sprintf(format, args...)
	Print(level, "%s (took %v)", message, elapsed)
}

func DebugTiming(startTime time.Time, format string, args ...interface{}) {
	PrintWithTiming(DebugLevel, startTime, format, args...)
}

func TraceTiming(startTime time.Time, format string, args ...interface{}) {
	PrintWithTiming(TraceLevel, startTime, format, args...)
}
