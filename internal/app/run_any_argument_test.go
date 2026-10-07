// Tests that `goremi <argument>` opens nothing and says the command is unknown (task 8.20).
package app

import (
	"io"
	"path/filepath"
	"strings"
	"testing"
)

// anyArgs are the arguments the command no longer knows: the removed commands and a mistake.
var anyArgs = []string{"theme", "help", "--help", "-h", "foo"}

// START: runArg

// runArg runs the command with one argument and returns the models it ran, what it wrote to out and its error.
func runArg(t *testing.T, arg string) (int, string, error) {
	t.Helper()
	run, ran := fakeRun()
	var out strings.Builder
	err := Run([]string{arg}, filepath.Join(t.TempDir(), "config.toml"), fakeProvider{}, &out, run)
	return len(*ran), out.String(), err
}

// END: runArg

// START: TestRunAnyArgumentOpensNothing

func TestRunAnyArgumentOpensNothing(t *testing.T) {
	for _, arg := range anyArgs {
		if ran, _, err := runArg(t, arg); ran != 0 || err == nil {
			t.Errorf("%q: ran %d models, error %v; want no model and an error", arg, ran, err)
		}
	}
}

// END: TestRunAnyArgumentOpensNothing

// START: TestRunAnyArgumentMessage

func TestRunAnyArgumentMessage(t *testing.T) {
	for _, arg := range anyArgs {
		_, _, err := runArg(t, arg)
		if want := `Unknown command "` + arg + `".`; err == nil || err.Error() != want {
			t.Errorf("%q: error %v, want %q", arg, err, want)
		}
	}
}

// END: TestRunAnyArgumentMessage

// START: TestRunAnyArgumentExitsTwo

func TestRunAnyArgumentExitsTwo(t *testing.T) {
	for _, arg := range anyArgs {
		_, _, err := runArg(t, arg)
		if got := ExitCode(err); got != 2 {
			t.Errorf("%q: ExitCode = %d, want 2", arg, got)
		}
	}
}

// END: TestRunAnyArgumentExitsTwo

// START: TestRunHelpPrintsNothing

func TestRunHelpPrintsNothing(t *testing.T) {
	_, out, _ := runArg(t, "help")
	if out != "" {
		t.Fatalf("goremi help wrote %q, want nothing (the help command is gone)", out)
	}
}

// END: TestRunHelpPrintsNothing

// START: TestRunRemovedCommandsOpenNoPicker

func TestRunRemovedCommandsOpenNoPicker(t *testing.T) {
	for _, arg := range []string{"theme", "help"} {
		run, ran := fakeRun()
		err := Run([]string{arg}, filepath.Join(t.TempDir(), "config.toml"), fakeProvider{}, io.Discard, run)
		if len(*ran) != 0 || err == nil {
			t.Errorf("%q: ran %d models, error %v; want none and an error", arg, len(*ran), err)
		}
	}
}

// END: TestRunRemovedCommandsOpenNoPicker
