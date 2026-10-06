// Helpers for the spectrum tests of the app: a player that reports band levels, and an app of 60 columns (panel above list) with one result.
package app

import (
	"strings"

	"goremi/internal/player"
	"goremi/internal/provider"
	"goremi/internal/ui"
)

// bandsPlayer is a recording player that also reports the same level for every band.
type bandsPlayer struct {
	*recordingPlayer
	db float64
	on bool
}

func (b *bandsPlayer) Bands() [ui.SpectrumBars]float64 {
	var out [ui.SpectrumBars]float64
	for i := range out {
		out[i] = b.db
	}
	return out
}

func (b *bandsPlayer) Spectrum() bool { return b.on }

// newBandsPlayer is a bands player whose levels are db, with the spectrum on.
func newBandsPlayer(db float64) *bandsPlayer {
	return &bandsPlayer{recordingPlayer: &recordingPlayer{events: make(chan player.Event, 1)}, db: db, on: true}
}

var specTrack = provider.Track{ID: "s1", Title: "One", Artist: "Daft Punk", Artwork: "http://x/a.jpg"}

// specApp is a 60-column app (search, panel, list from top to bottom) with the result specTrack, this player and the spectrum on or off.
func specApp(pl Player, spectrum bool) Model {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{specTrack}}}}
	return search(sized(New(p).WithPlayer(pl).WithSpectrum(spectrum), 60), "daft")
}

// panelLines are the lines of the view after the Search: line, ANSI stripped.
func panelLines(m Model) []string {
	lines := plainLines(m)
	for i, l := range lines {
		if strings.HasPrefix(l, "Search:") {
			return lines[i+1:]
		}
	}
	return nil
}

// artistIndex is the index of the artist line among the panel lines, or -1.
func artistIndex(m Model) int {
	for i, l := range panelLines(m) {
		if strings.Contains(l, "Daft Punk") && !strings.Contains(l, "▶") {
			return i
		}
	}
	return -1
}
