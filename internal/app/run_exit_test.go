// Tests for the exit code of the goremi command and for the help flags.
package app

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// START: TestExitCodeUnknownArg

func TestExitCodeUnknownArg(t *testing.T) {
	_, err := runUnknown(t, "foo")
	if got := ExitCode(err); got != 2 {
		t.Fatalf("ExitCode = %d, want 2", got)
	}
}

// END: TestExitCodeUnknownArg

// START: TestExitCodeOtherError

func TestExitCodeOtherError(t *testing.T) {
	failing := func(tea.Model) (tea.Model, error) { return nil, errors.New("boom") }
	err := Run(nil, filepath.Join(t.TempDir(), "config.toml"), fakeProvider{}, &strings.Builder{}, failing)
	if got := ExitCode(err); got != 1 {
		t.Fatalf("ExitCode = %d, want 1 (error %v)", got, err)
	}
}

// END: TestExitCodeOtherError

// START: TestExitCodeNil

func TestExitCodeNil(t *testing.T) {
	if got := ExitCode(nil); got != 0 {
		t.Fatalf("ExitCode(nil) = %d, want 0", got)
	}
}

// END: TestExitCodeNil

// helpWith runs the command with one argument and returns the models run, the text written and the error.
func helpWith(t *testing.T, arg string) (int, string, error) {
	t.Helper()
	run, ran := fakeRun()
	var out strings.Builder
	err := Run([]string{arg}, filepath.Join(t.TempDir(), "config.toml"), fakeProvider{}, &out, run)
	return len(*ran), out.String(), err
}

// checkLikeHelp fails unless the flag does what `help` does.
func checkLikeHelp(t *testing.T, flag string) {
	t.Helper()
	_, want, _ := helpWith(t, "help")
	ran, got, err := helpWith(t, flag)
	if ran != 0 || err != nil || got != want || got == "" {
		t.Fatalf("%s: ran %d models, error %v, text %q; want no model, no error and the text of help %q", flag, ran, err, got, want)
	}
}
