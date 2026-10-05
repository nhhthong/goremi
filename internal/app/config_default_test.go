// Tests for config defaults, fallbacks and saving the theme.
package app

import (
	"os"
	"path/filepath"
	"testing"
)

// START: TestShowKeysDefaultTrue

func TestShowKeysDefaultTrue(t *testing.T) {
	got := LoadConfig(writeConfig(t, "[ui]\ntheme = \"default\"\n"))
	if !got.ShowSpectrum {
		t.Fatalf("LoadConfig = %+v, want the show key true", got)
	}
}

// END: TestShowKeysDefaultTrue

// START: TestShowKeyFalseKept

func TestShowKeyFalseKept(t *testing.T) {
	got := LoadConfig(writeConfig(t, "[ui]\nshow_spectrum = false\n"))
	if got.ShowSpectrum {
		t.Fatalf("LoadConfig = %+v, want ShowSpectrum false", got)
	}
}

// END: TestShowKeyFalseKept

// START: TestUnlistedThemeFallsBack

func TestUnlistedThemeFallsBack(t *testing.T) {
	if got := LoadConfig(writeConfig(t, "[ui]\ntheme = \"neon\"\n")); got.Theme != "default" {
		t.Fatalf("Theme = %q, want default", got.Theme)
	}
}

// END: TestUnlistedThemeFallsBack

// START: TestMissingConfigDefaults

func TestMissingConfigDefaults(t *testing.T) {
	got := LoadConfig(filepath.Join(t.TempDir(), "none.toml"))
	if want := (Config{Theme: "default", ShowSpectrum: true, Mouse: true}); got != want {
		t.Fatalf("LoadConfig = %+v, want %+v", got, want)
	}
}

// END: TestMissingConfigDefaults

// START: TestBrokenConfigFallsBack

func TestBrokenConfigFallsBack(t *testing.T) {
	path := writeConfig(t, "[ui]\ntheme = default\n")
	if got, want := LoadConfig(path), (Config{Theme: "default", ShowSpectrum: true, Mouse: true}); got != want {
		t.Fatalf("broken file: LoadConfig = %+v, want %+v", got, want)
	}
	if err := os.WriteFile(path, []byte("[ui]\ntheme = \"nord\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := LoadConfig(path); got.Theme != "nord" {
		t.Fatalf("after the fix: Theme = %q, want nord", got.Theme)
	}
}

// END: TestBrokenConfigFallsBack

// START: TestSaveThemeCreatesDir

func TestSaveThemeCreatesDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "config.toml")
	if err := SaveTheme(path, "nord"); err != nil {
		t.Fatal(err)
	}
	if got := LoadConfig(path); got.Theme != "nord" {
		t.Fatalf("Theme = %q, want nord", got.Theme)
	}
}

// END: TestSaveThemeCreatesDir

// START: TestSaveThemeKeepsOtherKeys

func TestSaveThemeKeepsOtherKeys(t *testing.T) {
	path := writeConfig(t, "[ui]\ntheme = \"default\"\nshow_spectrum = false\n")
	if err := SaveTheme(path, "dracula"); err != nil {
		t.Fatal(err)
	}
	want := Config{Theme: "dracula", ShowSpectrum: false, Mouse: true}
	if got := LoadConfig(path); got != want {
		t.Fatalf("LoadConfig = %+v, want %+v", got, want)
	}
}

// END: TestSaveThemeKeepsOtherKeys
