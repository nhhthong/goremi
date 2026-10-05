// Tests that the stderr of a failing yt-dlp reaches the log.
package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// START: fakeYtDlp

// fakeYtDlp puts a yt-dlp on PATH that prints stderr to its standard error and stdout to its standard output, then exits with code.
func fakeYtDlp(t *testing.T, stderr, stdout string, code int) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the fake yt-dlp is a shell script")
	}
	dir := t.TempDir()
	script := fmt.Sprintf("#!/bin/sh\nprintf '%%s\\n' %q >&2\nprintf '%%s' %q\nexit %d\n", stderr, stdout, code)
	if err := os.WriteFile(filepath.Join(dir, "yt-dlp"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+"/usr/bin:/bin")
}

// logged returns a provider that records its log lines, and the record.
func logged() (*YouTubeProvider, *[]string) {
	var lines []string
	return &YouTubeProvider{Log: func(format string, a ...any) { lines = append(lines, fmt.Sprintf(format, a...)) }}, &lines
}

// loggedOnce checks that exactly one log call holds want.
func loggedOnce(t *testing.T, lines []string, want string) {
	t.Helper()
	if len(lines) != 1 || !strings.Contains(lines[0], want) {
		t.Fatalf("log = %q, want one call holding %q", lines, want)
	}
}

// END: fakeYtDlp

// START: failure logging tests

func TestSearchFailureLogsStderr(t *testing.T) {
	fakeYtDlp(t, "ERROR: boom", "", 1)
	p, lines := logged()
	if _, err := p.Search("daft punk", 1); err == nil {
		t.Fatal("Search returned no error")
	}
	loggedOnce(t, *lines, "ERROR: boom")
}

func TestResolveFailureLogsStderr(t *testing.T) {
	fakeYtDlp(t, "ERROR: boom", "", 1)
	p, lines := logged()
	if _, err := p.Resolve(Track{ID: "a1"}); err == nil {
		t.Fatal("Resolve returned no error")
	}
	loggedOnce(t, *lines, "ERROR: boom")
}

func TestDetailsFailureLogsStderr(t *testing.T) {
	fakeYtDlp(t, "ERROR: boom", "", 1)
	p, lines := logged()
	if _, err := p.Details(Track{ID: "a1"}); err == nil {
		t.Fatal("Details returned no error")
	}
	loggedOnce(t, *lines, "ERROR: boom")
}

func TestSuccessLogsNothing(t *testing.T) {
	fakeYtDlp(t, "WARNING: slow", "https://example.test/stream\n", 0)
	p, lines := logged()
	if _, err := p.Resolve(Track{ID: "a1"}); err != nil {
		t.Fatal(err)
	}
	if len(*lines) != 0 {
		t.Fatalf("log = %q, want nothing", *lines)
	}
}

// END: failure logging tests
