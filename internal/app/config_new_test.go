// Tests for starting the app from the config file.
package app

import (
	"path/filepath"
	"testing"

	"goremi/internal/ui/theme"
)

// START: TestFromConfigUsesTheme

func TestFromConfigUsesTheme(t *testing.T) {
	m := NewFromConfig(writeConfig(t, "[ui]\ntheme = \"cyberpunk\"\n"), fakeProvider{})
	if m.Theme() != theme.Cyberpunk() {
		t.Fatalf("Theme() = %v, want the cyberpunk theme", m.Theme())
	}
}

// END: TestFromConfigUsesTheme

// START: TestFromConfigDefaultsDark

func TestFromConfigDefaultsDark(t *testing.T) {
	m := NewFromConfig(filepath.Join(t.TempDir(), "none.toml"), fakeProvider{})
	if m.Theme() != theme.Dark() {
		t.Fatalf("Theme() = %v, want the dark theme", m.Theme())
	}
}

// END: TestFromConfigDefaultsDark
