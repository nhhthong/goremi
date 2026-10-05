// Tests for the keys that play the next and the previous track of the results list.
package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

var (
	trackA = provider.Track{ID: "a", Title: "A"}
	trackB = provider.Track{ID: "b", Title: "B"}
	trackC = provider.Track{ID: "c", Title: "C"}
)

// playingModel is a model whose list holds A, B and C (list focused) and which plays the given track.
func playingModel(playing provider.Track) Model { return playingModelWith(&recordingPlayer{}, playing) }

// playingModelWith is playingModel with the given fake player.
func playingModelWith(pl *recordingPlayer, playing provider.Track) Model {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{trackA, trackB, trackC}}}}
	m := search(sized(New(p).WithPlayer(pl), 100), "daft")
	next, _ := m.Update(PlayMsg{Track: playing})
	return next.(Model)
}

// pressKey presses one rune key on the model and returns the command.
func pressKey(m Model, r rune) tea.Cmd {
	_, cmd := m.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	return cmd
}

// START: TestKeyNPlaysNext

func TestKeyNPlaysNext(t *testing.T) {
	cmd := pressKey(playingModel(trackA), 'n')
	if cmd == nil {
		t.Fatal("n returned no command")
	}
	if msg := cmd(); msg != (PlayMsg{Track: trackB}) {
		t.Fatalf("n gave %v, want PlayMsg for B", msg)
	}
}

// END: TestKeyNPlaysNext

// START: TestKeyNAtLastDoesNothing

func TestKeyNAtLastDoesNothing(t *testing.T) {
	if cmd := pressKey(playingModel(trackC), 'n'); cmd != nil {
		t.Fatalf("n on the last track returned a command: %v", cmd())
	}
}

// END: TestKeyNAtLastDoesNothing

// START: TestKeyPPlaysPrevious

func TestKeyPPlaysPrevious(t *testing.T) {
	cmd := pressKey(playingModel(trackB), 'p')
	if cmd == nil {
		t.Fatal("p returned no command")
	}
	if msg := cmd(); msg != (PlayMsg{Track: trackA}) {
		t.Fatalf("p gave %v, want PlayMsg for A", msg)
	}
}

// END: TestKeyPPlaysPrevious

// START: TestKeyPAtFirstDoesNothing

func TestKeyPAtFirstDoesNothing(t *testing.T) {
	if cmd := pressKey(playingModel(trackA), 'p'); cmd != nil {
		t.Fatalf("p on the first track returned a command: %v", cmd())
	}
}

// END: TestKeyPAtFirstDoesNothing
