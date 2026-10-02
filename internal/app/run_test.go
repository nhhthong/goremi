// Tests for the goremi command: `goremi theme` picks a theme first, `goremi` opens the app.
package app

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/ui/theme"
)

// fakeRun returns a run function that feeds keys to a ThemeModel, only records an app Model, and the list of the models it was given.
func fakeRun(keys ...tea.KeyPressMsg) (func(tea.Model) (tea.Model, error), *[]tea.Model) {
	var ran []tea.Model
	return func(m tea.Model) (tea.Model, error) {
		ran = append(ran, m)
		if tm, ok := m.(ThemeModel); ok {
			final, _ := press(tm, keys...)
			return final, nil
		}
		return m, nil
	}, &ran
}

// START: TestRunThemeOpensAppAfterEnter

func TestRunThemeOpensAppAfterEnter(t *testing.T) {
	path := writeConfig(t, "[ui]\ntheme = \"default\"\n")
	run, ran := fakeRun(thDown, thEnter)
	if err := Run([]string{"theme"}, path, fakeProvider{}, io.Discard, run); err != nil {
		t.Fatal(err)
	}
	if len(*ran) != 2 {
		t.Fatalf("ran %d models, want 2 (picker, app)", len(*ran))
	}
	if _, ok := (*ran)[0].(ThemeModel); !ok {
		t.Errorf("first model is %T, want the picker", (*ran)[0])
	}
	if m, ok := (*ran)[1].(Model); !ok || m.Theme() != theme.Catppuccin() {
		t.Errorf("second model is %T, want the app with the catppuccin theme", (*ran)[1])
	}
	if got := LoadConfig(path).Theme; got != "catppuccin" {
		t.Errorf("saved Theme = %q, want catppuccin", got)
	}
}

// END: TestRunThemeOpensAppAfterEnter

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

// START: TestRunThemeEscExits

func TestRunThemeEscExits(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	run, ran := fakeRun(thDown, thEsc)
	if err := Run([]string{"theme"}, path, fakeProvider{}, io.Discard, run); err != nil {
		t.Fatal(err)
	}
	if _, ok := (*ran)[0].(ThemeModel); !ok || len(*ran) != 1 {
		t.Errorf("ran %v, want only the picker", *ran)
	}
	if _, err := os.Stat(path); err == nil {
		t.Errorf("the config file %s exists, want nothing saved", path)
	}
}

// END: TestRunThemeEscExits
