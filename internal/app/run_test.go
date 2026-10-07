// Tests for the goremi command: `goremi` opens the app.
package app

import (
	"io"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// fakeRun returns a run function that only records the models it was given, and the list of them.
func fakeRun() (func(tea.Model) (tea.Model, error), *[]tea.Model) {
	var ran []tea.Model
	return func(m tea.Model) (tea.Model, error) {
		ran = append(ran, m)
		return m, nil
	}, &ran
}

// START: TestRunWithoutArgsSkipsPicker

func TestRunWithoutArgsSkipsPicker(t *testing.T) {
	run, ran := fakeRun()
	if err := Run(nil, filepath.Join(t.TempDir(), "config.toml"), fakeProvider{}, io.Discard, run); err != nil {
		t.Fatal(err)
	}
	if _, ok := (*ran)[0].(Model); !ok || len(*ran) != 1 {
		t.Fatalf("ran %v, want only the app", *ran)
	}
}

// END: TestRunWithoutArgsSkipsPicker
