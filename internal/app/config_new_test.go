// Tests for starting the app from the config file.
package app

import (
	"path/filepath"
	"testing"

	"goremi/internal/ui/theme"
)

// START: TestFromConfigUsesTheme

func TestFromConfigUsesTheme(t *testing.T) {
	m := NewFromConfig(writeConfig(t, "[ui]\ntheme = \"dracula\"\n"), fakeProvider{})
	if m.Theme() != theme.Dracula() {
		t.Fatalf("Theme() = %v, want the dracula theme", m.Theme())
	}
}

// END: TestFromConfigUsesTheme

// START: TestFromConfigDefaultsDefault

func TestFromConfigDefaultsDefault(t *testing.T) {
	m := NewFromConfig(filepath.Join(t.TempDir(), "none.toml"), fakeProvider{})
	if m.Theme() != theme.Default() {
		t.Fatalf("Theme() = %v, want the default theme", m.Theme())
	}
}

// END: TestFromConfigDefaultsDefault
