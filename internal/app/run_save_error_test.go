// Tests for `goremi theme` when the theme cannot be saved.
package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// blockedPath returns a config path below a regular file, so no directory can be made for it.
func blockedPath(t *testing.T) string {
	t.Helper()
	file := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	return filepath.Join(file, "config.toml")
}

// START: TestRunSaveFailureMessage

func TestRunSaveFailureMessage(t *testing.T) {
	run, ran := fakeRun(thDown, thEnter)
	err := Run([]string{"theme"}, blockedPath(t), fakeProvider{}, &strings.Builder{}, run)
	if err == nil || !strings.HasPrefix(err.Error(), "Cannot save the theme:") {
		t.Fatalf("error %v, want a message starting with %q", err, "Cannot save the theme:")
	}
	if _, ok := (*ran)[0].(ThemeModel); !ok || len(*ran) != 1 {
		t.Fatalf("ran %v, want only the picker", *ran)
	}
}

// END: TestRunSaveFailureMessage

// START: TestRunSaveFailureExitsAndRecovers

func TestRunSaveFailureExitsAndRecovers(t *testing.T) {
	run, _ := fakeRun(thDown, thEnter)
	err := Run([]string{"theme"}, blockedPath(t), fakeProvider{}, &strings.Builder{}, run)
	if got := ExitCode(err); got != 1 {
		t.Fatalf("ExitCode = %d, want 1 (error %v)", got, err)
	}
	run, ran := fakeRun(thDown, thEnter)
	good := filepath.Join(t.TempDir(), "config.toml")
	if err := Run([]string{"theme"}, good, fakeProvider{}, &strings.Builder{}, run); err != nil || len(*ran) != 2 {
		t.Fatalf("after the file system works: error %v, ran %d models, want no error and 2", err, len(*ran))
	}
}

// END: TestRunSaveFailureExitsAndRecovers
