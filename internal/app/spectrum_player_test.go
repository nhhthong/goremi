// Tests that the spectrum switch of the config reaches mpv: the player starts it with the band filter only when the spectrum shows.
package app

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/player"
)

// optionsOfFirstPlay plays on an mpvPlayer with this switch and returns the options its start function got.
func optionsOfFirstPlay(spectrum bool) player.Options {
	var got player.Options
	pl := newMpvPlayer()
	pl.spectrum = spectrum
	pl.start = func(_ context.Context, opts player.Options) (*player.Player, error) {
		got = opts
		return nil, errors.New("stop here")
	}
	pl.Play("x")
	return got
}

// START: TestMpvPlayerStartsWithoutFilterWhenOff

func TestMpvPlayerStartsWithoutFilterWhenOff(t *testing.T) {
	if got := optionsOfFirstPlay(false); got.Spectrum {
		t.Fatalf("start options = %+v, want Spectrum false", got)
	}
}

// END: TestMpvPlayerStartsWithoutFilterWhenOff

// START: TestMpvPlayerStartsWithFilterWhenOn

func TestMpvPlayerStartsWithFilterWhenOn(t *testing.T) {
	if got := optionsOfFirstPlay(true); !got.Spectrum {
		t.Fatalf("start options = %+v, want Spectrum true", got)
	}
}

// END: TestMpvPlayerStartsWithFilterWhenOn

// runSpectrum runs Run with a config file holding content (none when empty) and returns the spectrum switch of the player and of the model it opened.
func runSpectrum(t *testing.T, content string) (playerOn, modelOn bool) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	if content != "" {
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	run := func(m tea.Model) (tea.Model, error) {
		model := m.(Model)
		playerOn, modelOn = model.player.(*mpvPlayer).spectrum, model.Spectrum()
		return m, nil
	}
	if err := Run(nil, path, fakeProvider{}, io.Discard, run); err != nil {
		t.Fatal(err)
	}
	return playerOn, modelOn
}

// START: TestRunGivesShowSpectrumToPlayer

func TestRunGivesShowSpectrumToPlayer(t *testing.T) {
	if p, m := runSpectrum(t, "[ui]\nshow_spectrum = false\n"); p || m {
		t.Fatalf("show_spectrum = false: player %v, model %v, want both false", p, m)
	}
	if p, m := runSpectrum(t, "[ui]\nshow_spectrum = true\n"); !p || !m {
		t.Fatalf("show_spectrum = true: player %v, model %v, want both true", p, m)
	}
}

// END: TestRunGivesShowSpectrumToPlayer

// START: TestRunDefaultsSpectrumOn

func TestRunDefaultsSpectrumOn(t *testing.T) {
	if p, m := runSpectrum(t, ""); !p || !m {
		t.Fatalf("no config: player %v, model %v, want both true", p, m)
	}
	if p, m := runSpectrum(t, "[ui]\ntheme = \"nord\"\n"); !p || !m {
		t.Fatalf("config without the key: player %v, model %v, want both true", p, m)
	}
}

// END: TestRunDefaultsSpectrumOn
