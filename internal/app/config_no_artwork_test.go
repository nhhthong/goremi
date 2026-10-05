// Tests that the config no longer knows show_artwork: a leftover key is ignored and a saved config does not write it.
package app

import (
	"os"
	"strings"
	"testing"
)

// START: TestLoadConfigIgnoresShowArtwork

func TestLoadConfigIgnoresShowArtwork(t *testing.T) {
	got := LoadConfig(writeConfig(t, "[ui]\ntheme = \"nord\"\nshow_spectrum = false\nshow_artwork = true\n"))
	if got.Theme != "nord" || got.ShowSpectrum || !got.Mouse {
		t.Fatalf("LoadConfig = %+v, want Theme nord, ShowSpectrum false, Mouse true", got)
	}
}

// END: TestLoadConfigIgnoresShowArtwork

// START: TestSaveThemeDropsShowArtwork

func TestSaveThemeDropsShowArtwork(t *testing.T) {
	path := writeConfig(t, "[ui]\nshow_artwork = false\nmouse = false\n")
	if err := SaveTheme(path, "nord"); err != nil {
		t.Fatal(err)
	}
	text, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(text), "show_artwork") {
		t.Fatalf("the saved config still holds show_artwork:\n%s", text)
	}
	if got := LoadConfig(path); got.Mouse || got.Theme != "nord" {
		t.Fatalf("LoadConfig = %+v, want Mouse false and Theme nord", got)
	}
}

// END: TestSaveThemeDropsShowArtwork

// START: TestLoadConfigStaleKeyOfOtherTypeIsIgnored

func TestLoadConfigStaleKeyOfOtherTypeIsIgnored(t *testing.T) {
	if got := LoadConfig(writeConfig(t, "[ui]\ntheme = \"nord\"\nshow_artwork = \"yes\"\n")); got.Theme != "nord" {
		t.Fatalf("Theme = %q, want nord", got.Theme)
	}
}

// END: TestLoadConfigStaleKeyOfOtherTypeIsIgnored
