// Tests for `goremi help`.
package app

import (
	"path/filepath"
	"strings"
	"testing"
)

// runHelp runs `goremi help` and returns the number of models it ran, the text it wrote and its error.
func runHelp(t *testing.T) (int, string, error) {
	t.Helper()
	run, ran := fakeRun()
	var out strings.Builder
	err := Run([]string{"help"}, filepath.Join(t.TempDir(), "config.toml"), fakeProvider{}, &out, run)
	return len(*ran), out.String(), err
}

// START: TestRunHelpOpensNothing

func TestRunHelpOpensNothing(t *testing.T) {
	ran, _, err := runHelp(t)
	if ran != 0 || err != nil {
		t.Fatalf("ran %d models and returned %v, want no model and no error", ran, err)
	}
}

// END: TestRunHelpOpensNothing

// START: TestRunHelpPrintsText

func TestRunHelpPrintsText(t *testing.T) {
	if _, out, _ := runHelp(t); strings.TrimSpace(out) == "" {
		t.Fatal("goremi help wrote nothing")
	}
}

// END: TestRunHelpPrintsText
