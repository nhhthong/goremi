// Tests for the player panel around a play: the logo before it, and after it the artist at the top with no image and nothing sent to the terminal.
package app

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: panelAfterPlayHelpers

var afterPlayTrack = provider.Track{ID: "p1", Title: "One", Artist: "Daft Punk", Artwork: "http://x/a.jpg"}

// afterPlayApp is a 100-column app holding the results of a search for afterPlayTrack, with a recording player.
func afterPlayApp() Model {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{afterPlayTrack}}}}
	return search(sized(New(p).WithPlayer(&recordingPlayer{}), 100), "daft")
}

// playAndSettle plays the track, runs every command it starts and feeds the messages back until no command is left; it returns the model and every message the commands produced.
func playAndSettle(m Model, track provider.Track) (Model, []tea.Msg) {
	var all []tea.Msg
	next, cmd := m.Update(PlayMsg{Track: track})
	m = next.(Model)
	for i := 0; i < 8 && cmd != nil; i++ {
		var more []tea.Cmd
		for _, msg := range collect(cmd) {
			all = append(all, msg)
			if _, ok := msg.(tickMsg); ok {
				continue
			}
			if _, ok := msg.(noteTickMsg); ok { // the note loop never ends by itself
				continue
			}
			n, c := m.Update(msg)
			m = n.(Model)
			if c != nil {
				more = append(more, c)
			}
		}
		cmd = tea.Batch(more...)
	}
	return m, all
}

// END: panelAfterPlayHelpers

// START: TestNoLogoNoImageAfterPlay

func TestNoLogoNoImageAfterPlay(t *testing.T) {
	m, _ := playAndSettle(afterPlayApp(), afterPlayTrack)
	view := m.View().Content
	if strings.ContainsRune(view, '\U0010EEEE') {
		t.Fatal("the view holds an image placeholder cell")
	}
	for _, l := range plainLines(m) {
		if strings.Contains(l, art0()) {
			t.Fatalf("the logo is still on screen: %q", l)
		}
	}
}

// END: TestNoLogoNoImageAfterPlay

// START: TestPanelStartsWithArtistAfterPlay

func TestPanelStartsWithArtistAfterPlay(t *testing.T) {
	m, _ := playAndSettle(afterPlayApp(), afterPlayTrack)
	for _, l := range plainLines(m) {
		if strings.Contains(l, "▶ One") {
			if !strings.Contains(l, "Daft Punk") {
				t.Fatalf("first panel line = %q, want the artist line", l)
			}
			return
		}
	}
	t.Fatal("no line with the playing track")
}

// END: TestPanelStartsWithArtistAfterPlay

// START: TestPlaySendsNothingToTerminal

func TestPlaySendsNothingToTerminal(t *testing.T) {
	_, msgs := playAndSettle(afterPlayApp(), afterPlayTrack)
	for _, msg := range msgs {
		if raw, ok := msg.(tea.RawMsg); ok {
			t.Fatalf("the play sent %#v to the terminal", raw.Msg)
		}
	}
}

// END: TestPlaySendsNothingToTerminal
