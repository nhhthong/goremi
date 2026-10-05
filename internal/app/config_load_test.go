// Tests for reading config.toml.
package app

import (
	"os"
	"path/filepath"
	"testing"
)

// writeConfig writes content to a config.toml in a temp dir and returns its path.
func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// START: TestLoadConfigTheme

func TestLoadConfigTheme(t *testing.T) {
	got := LoadConfig(writeConfig(t, "[ui]\ntheme = \"dracula\"\n"))
	if got.Theme != "dracula" {
		t.Fatalf("Theme = %q, want dracula", got.Theme)
	}
}

// END: TestLoadConfigTheme

// START: TestLoadConfigAllKeys

func TestLoadConfigAllKeys(t *testing.T) {
	got := LoadConfig(writeConfig(t, "[ui]\ntheme = \"nord\"\nshow_artwork = true\nshow_spectrum = true\n"))
	want := Config{Theme: "nord", ShowArtwork: true, ShowSpectrum: true, Mouse: true}
	if got != want {
		t.Fatalf("LoadConfig = %+v, want %+v", got, want)
	}
}

// END: TestLoadConfigAllKeys
