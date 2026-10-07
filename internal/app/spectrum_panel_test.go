// Tests for the top of the player panel: the logo before the first play, then the spectrum (8 rows above the artist line) or nothing when it is off.
package app

import (
	"os"
	"path/filepath"
	"testing"
	"unicode/utf8"
)

// START: TestSpectrumAboveArtistAfterPlay

func TestSpectrumAboveArtistAfterPlay(t *testing.T) {
	m, _ := playAndSettle(specApp(&recordingPlayer{}, true), specTrack)
	if got := artistIndex(m); got != 9 {
		t.Fatalf("the artist line is the panel line %d, want 9 (below the 8 rows of the spectrum and the baseline): %q", got, panelLines(m))
	}
}

// END: TestSpectrumAboveArtistAfterPlay

// START: TestSpectrumRowsAre40Wide

func TestSpectrumRowsAre40Wide(t *testing.T) {
	m, _ := playAndSettle(specApp(&recordingPlayer{}, true), specTrack)
	for i, l := range panelLines(m)[:8] {
		if n := utf8.RuneCountInString(l); n != 40 {
			t.Fatalf("spectrum row %d is %d columns (%q), want 40", i, n, l)
		}
	}
}

// END: TestSpectrumRowsAre40Wide

// START: TestSpectrumRowsDrawTheBars

func TestSpectrumRowsDrawTheBars(t *testing.T) {
	m, _ := playAndSettle(specApp(&recordingPlayer{}, true), specTrack)
	m.bars[0] = 64
	top := []rune(panelLines(m)[0])
	if len(top) < 1 || top[0] != '▄' {
		t.Fatalf("top spectrum row = %q, want ▄ in column 0 for a bar of 64", string(top))
	}
}

// END: TestSpectrumRowsDrawTheBars

// START: TestNoSpectrumRowsWhenOff

func TestNoSpectrumRowsWhenOff(t *testing.T) {
	m, _ := playAndSettle(specApp(&recordingPlayer{}, false), specTrack)
	if got := artistIndex(m); got != 0 {
		t.Fatalf("the artist line is the panel line %d, want 0 (nothing above it): %q", got, panelLines(m))
	}
}

// END: TestNoSpectrumRowsWhenOff

// START: TestConfigShowSpectrumFalse

func TestConfigShowSpectrumFalse(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[ui]\nshow_spectrum = false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if NewFromConfig(path, fakeProvider{}).Spectrum() {
		t.Fatal("Spectrum() is true with show_spectrum = false")
	}
	if !NewFromConfig(filepath.Join(t.TempDir(), "none.toml"), fakeProvider{}).Spectrum() {
		t.Fatal("Spectrum() is false with no config, want true (the default)")
	}
}

// END: TestConfigShowSpectrumFalse
