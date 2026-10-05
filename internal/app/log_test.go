// Tests for the path of the log file.
package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// START: TestLogPathUnderCacheDir

func TestLogPathUnderCacheDir(t *testing.T) {
	tmp := t.TempDir()
	cache := tmp
	switch runtime.GOOS {
	case "darwin":
		t.Setenv("HOME", tmp)
		cache = filepath.Join(tmp, "Library", "Caches")
	case "windows":
		t.Setenv("LocalAppData", tmp)
	default:
		t.Setenv("XDG_CACHE_HOME", tmp)
	}
	got, err := LogPath()
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(cache, "goremi", "goremi.log"); got != want {
		t.Fatalf("LogPath = %q, want %q", got, want)
	}
}

// END: TestLogPathUnderCacheDir

// START: log helpers

// closeLog closes the log file at the end of the test; Windows cannot remove a temp directory that holds an open file.
func closeLog(t *testing.T, l *Logger) {
	t.Helper()
	if c, ok := l.w.(io.Closer); ok {
		t.Cleanup(func() { _ = c.Close() })
	}
}

// logLines opens a log in a temp directory, runs write on it, and returns the lines of the file.
func logLines(t *testing.T, write func(*Logger)) []string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "goremi.log")
	l, err := OpenLog(path)
	if err != nil {
		t.Fatal(err)
	}
	closeLog(t, l)
	write(l)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
}

// END: log helpers

// START: TestLogLineStartsWithPID

func TestLogLineStartsWithPID(t *testing.T) {
	lines := logLines(t, func(l *Logger) { l.Printf("hello") })
	if want := fmt.Sprintf("%d hello", os.Getpid()); len(lines) != 1 || lines[0] != want {
		t.Fatalf("lines = %q, want [%q]", lines, want)
	}
}

// END: TestLogLineStartsWithPID

// START: TestLogEveryLineHasPID

func TestLogEveryLineHasPID(t *testing.T) {
	lines := logLines(t, func(l *Logger) { l.Printf("one\ntwo") })
	pid := os.Getpid()
	want := []string{fmt.Sprintf("%d one", pid), fmt.Sprintf("%d two", pid)}
	if len(lines) != 2 || lines[0] != want[0] || lines[1] != want[1] {
		t.Fatalf("lines = %q, want %q", lines, want)
	}
}

// END: TestLogEveryLineHasPID

// START: TestLogAppends

func TestLogAppends(t *testing.T) {
	path := filepath.Join(t.TempDir(), "goremi.log")
	for _, msg := range []string{"first", "second"} {
		l, err := OpenLog(path)
		if err != nil {
			t.Fatal(err)
		}
		closeLog(t, l)
		l.Printf("%s", msg)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	pid := os.Getpid()
	if want := fmt.Sprintf("%d first\n%d second\n", pid, pid); string(data) != want {
		t.Fatalf("file = %q, want %q", data, want)
	}
}

// END: TestLogAppends

// START: TestLogUnwritableDoesNotPanic

func TestLogUnwritableDoesNotPanic(t *testing.T) {
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	l, err := OpenLog(filepath.Join(file, "goremi.log"))
	if err == nil {
		t.Fatal("OpenLog under a regular file returned no error")
	}
	if l == nil {
		t.Fatal("OpenLog returned a nil logger")
	}
	l.Printf("still fine")
}

// END: TestLogUnwritableDoesNotPanic

// START: TestLogNilLoggerDoesNotPanic

func TestLogNilLoggerDoesNotPanic(t *testing.T) {
	var l *Logger
	l.Printf("nothing")
}

// END: TestLogNilLoggerDoesNotPanic
