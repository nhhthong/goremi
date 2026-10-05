// Tests for keys and Enter that act on the playing track: after a new search, and playing it again.
package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// listProvider is a scripted search provider that resolves every track to one URL.
type listProvider struct{ *scriptedProvider }

func (listProvider) Resolve(provider.Track) (string, error) { return "http://stream/1", nil }

var (
	trackD = provider.Track{ID: "d", Title: "D"}
	trackE = provider.Track{ID: "e", Title: "E"}
)

// playedAfterSearch is a list-focused model where A plays and a second search replaced the list with D and E.
func playedAfterSearch() Model {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{trackA, trackB}}, {tracks: []provider.Track{trackD, trackE}}}}
	m := search(sized(New(p).WithPlayer(&recordingPlayer{}), 100), "daft")
	next, _ := m.Update(PlayMsg{Track: trackA})
	return search(next.(Model), "punk")
}

// START: TestKeyNAfterNewSearch

func TestKeyNAfterNewSearch(t *testing.T) {
	if cmd := pressKey(playedAfterSearch(), 'n'); cmd != nil {
		t.Fatalf("n returned a command after the playing track left the list: %v", cmd())
	}
}

// END: TestKeyNAfterNewSearch

// START: TestKeyPAfterNewSearch

func TestKeyPAfterNewSearch(t *testing.T) {
	if cmd := pressKey(playedAfterSearch(), 'p'); cmd != nil {
		t.Fatalf("p returned a command after the playing track left the list: %v", cmd())
	}
}

// END: TestKeyPAfterNewSearch

// playingA is a list-focused model of A and B that plays A through the given fake player.
func playingA(pl *recordingPlayer) Model {
	p := listProvider{&scriptedProvider{steps: []step{{tracks: []provider.Track{trackA, trackB}}}}}
	m := search(sized(New(p).WithPlayer(pl), 100), "daft")
	m, _ = playOne(m, trackA)
	return m
}

// pressEnter presses Enter on the model and returns the command.
func pressEnter(m Model) tea.Cmd {
	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	return cmd
}

// START: TestEnterOnPlayingTrack

func TestEnterOnPlayingTrack(t *testing.T) {
	cmd := pressEnter(playingA(&recordingPlayer{}))
	if cmd == nil {
		t.Fatal("Enter on the playing track returned no command")
	}
	if msg := cmd(); msg != (PlayMsg{Track: trackA}) {
		t.Fatalf("Enter gave %v, want PlayMsg for A", msg)
	}
}

// END: TestEnterOnPlayingTrack

// START: TestEnterRestartsFromStart

func TestEnterRestartsFromStart(t *testing.T) {
	pl := &recordingPlayer{}
	m := playingA(pl)
	first := len(pl.urls)
	next, cmd := m.Update(pressEnter(m)())
	next, cmd = next.Update(cmd())
	cmd()
	if len(pl.urls) != first+1 || pl.urls[first] != pl.urls[first-1] || pl.toggles != 0 || len(pl.seeks) != 0 {
		t.Fatalf("Play calls %v, toggles %d, seeks %v, want the same URL played again and no pause or seek", pl.urls, pl.toggles, pl.seeks)
	}
	_ = next
}

// END: TestEnterRestartsFromStart
