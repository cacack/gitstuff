package verbosity

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestSetLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		expected Level
	}{
		{"Normal level", Normal, Normal},
		{"Info level", InfoLevel, InfoLevel},
		{"Debug level", DebugLevel, DebugLevel},
		{"Trace level", TraceLevel, TraceLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.level)
			if GetLevel() != tt.expected {
				t.Errorf("SetLevel(%v) = %v, want %v", tt.level, GetLevel(), tt.expected)
			}
		})
	}
}

func TestSetFromCount(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		expected Level
	}{
		{"Negative count", -1, Normal},
		{"Zero count", 0, Normal},
		{"Count 1", 1, InfoLevel},
		{"Count 2", 2, DebugLevel},
		{"Count 3", 3, TraceLevel},
		{"Count too high", 5, TraceLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetFromCount(tt.count)
			if GetLevel() != tt.expected {
				t.Errorf("SetFromCount(%d) = %v, want %v", tt.count, GetLevel(), tt.expected)
			}
		})
	}
}

func TestIsEnabled(t *testing.T) {
	tests := []struct {
		name         string
		currentLevel Level
		checkLevel   Level
		expected     bool
	}{
		{"Normal checks Normal", Normal, Normal, true},
		{"Normal checks Info", Normal, InfoLevel, false},
		{"Info checks Normal", InfoLevel, Normal, true},
		{"Info checks Info", InfoLevel, InfoLevel, true},
		{"Info checks Debug", InfoLevel, DebugLevel, false},
		{"Debug checks all", DebugLevel, TraceLevel, false},
		{"Trace checks all", TraceLevel, DebugLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.currentLevel)
			result := IsEnabled(tt.checkLevel)
			if result != tt.expected {
				t.Errorf("IsEnabled(%v) with level %v = %v, want %v",
					tt.checkLevel, tt.currentLevel, result, tt.expected)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	tests := []struct {
		name           string
		currentLevel   Level
		printLevel     Level
		message        string
		shouldPrint    bool
		expectedPrefix string
	}{
		{"Normal prints at Normal", Normal, Normal, "test", true, ""},
		{"Normal doesn't print at Info", Normal, InfoLevel, "test", false, ""},
		{"Info prints Info with prefix", InfoLevel, InfoLevel, "test", true, "ℹ️  "},
		{"Debug prints Debug with prefix", DebugLevel, DebugLevel, "test", true, "🐛 [DEBUG] "},
		{"Trace prints Trace with prefix", TraceLevel, TraceLevel, "test", true, "🔍 [TRACE] "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.currentLevel)

			// Capture output
			oldStdout := os.Stdout
			oldStderr := os.Stderr

			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Stderr = w

			Print(tt.printLevel, tt.message)

			w.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr

			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			output := buf.String()

			if tt.shouldPrint {
				if tt.expectedPrefix == "" {
					// Normal level goes to stdout
					if output != tt.message+"\n" {
						t.Errorf("Expected output '%s\\n', got '%s'", tt.message, output)
					}
				} else {
					// Other levels go to stderr with prefix
					expected := tt.expectedPrefix + tt.message + "\n"
					if output != expected {
						t.Errorf("Expected output '%s', got '%s'", expected, output)
					}
				}
			} else {
				if output != "" {
					t.Errorf("Expected no output, got '%s'", output)
				}
			}
		})
	}
}

func TestConvenienceFunctions(t *testing.T) {
	SetLevel(TraceLevel) // Enable all levels

	// Capture output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Info("info test")
	Debug("debug test")
	Trace("trace test")

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	expectedOutputs := []string{
		"ℹ️  info test\n",
		"🐛 [DEBUG] debug test\n",
		"🔍 [TRACE] trace test\n",
	}

	for _, expected := range expectedOutputs {
		if !bytes.Contains(buf.Bytes(), []byte(expected)) {
			t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
		}
	}
}

func TestPrintWithTiming(t *testing.T) {
	SetLevel(DebugLevel)

	startTime := time.Now().Add(-100 * time.Millisecond)

	// Capture output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	PrintWithTiming(DebugLevel, startTime, "operation completed")

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if !bytes.Contains(buf.Bytes(), []byte("🐛 [DEBUG] operation completed (took")) {
		t.Errorf("Expected timing output, got '%s'", output)
	}
}

func TestTimingConvenienceFunctions(t *testing.T) {
	SetLevel(TraceLevel)

	startTime := time.Now().Add(-50 * time.Millisecond)

	// Capture output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	DebugTiming(startTime, "debug operation")
	TraceTiming(startTime, "trace operation")

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	expectedOutputs := []string{
		"🐛 [DEBUG] debug operation (took",
		"🔍 [TRACE] trace operation (took",
	}

	for _, expected := range expectedOutputs {
		if !bytes.Contains(buf.Bytes(), []byte(expected)) {
			t.Errorf("Expected output to contain '%s', got '%s'", expected, output)
		}
	}
}
