// Tests for `goremi` with an argument it does not know.
package app

import (
	"io"
	"path/filepath"
	"strings"
	"testing"
)

// runUnknown runs the command with args and returns the number of models it ran and its error.
func runUnknown(t *testing.T, args ...string) (int, error) {
	t.Helper()
	run, ran := fakeRun()
	err := Run(args, filepath.Join(t.TempDir(), "config.toml"), fakeProvider{}, io.Discard, run)
	return len(*ran), err
}

// START: TestRunUnknownArgDoesNotOpenApp

func TestRunUnknownArgDoesNotOpenApp(t *testing.T) {
	ran, err := runUnknown(t, "foo")
	if ran != 0 || err == nil {
		t.Fatalf("ran %d models and returned %v, want no model and an error", ran, err)
	}
}

// END: TestRunUnknownArgDoesNotOpenApp

// START: TestRunUnknownArgHintsHelp

func TestRunUnknownArgHintsHelp(t *testing.T) {
	_, err := runUnknown(t, "foo")
	if err == nil || !strings.Contains(err.Error(), "goremi help") {
		t.Fatalf("error %v, want a message containing %q", err, "goremi help")
	}
}

// END: TestRunUnknownArgHintsHelp

// START: TestRunOldThemeFlagIsUnknown

func TestRunOldThemeFlagIsUnknown(t *testing.T) {
	ran, err := runUnknown(t, "--theme", "dark")
	if ran != 0 || err == nil || !strings.Contains(err.Error(), "goremi help") {
		t.Fatalf("ran %d models and returned %v, want no model and an error containing %q", ran, err, "goremi help")
	}
}

// END: TestRunOldThemeFlagIsUnknown
