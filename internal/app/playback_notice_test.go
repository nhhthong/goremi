// Tests for a play request that fails, and for the artist line of the panel.
package app

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// scriptedResolve is a provider whose Resolve answers from errs in turn (nil means success) and whose Details answers with artist.
type scriptedResolve struct {
	fakeProvider
	errs   []error
	artist string
}

func (s *scriptedResolve) Resolve(provider.Track) (string, error) {
	err := s.errs[0]
	s.errs = s.errs[1:]
	return "http://stream/1", err
}

func (s *scriptedResolve) Details(t provider.Track) (provider.Track, error) {
	t.Artist = s.artist
	return t, nil
}

// playOne sends PlayMsg for the track, feeds the resolved URL back, runs the Play command and feeds its outcome back; it returns the model and the command that follows (the details).
func playOne(m Model, track provider.Track) (Model, tea.Cmd) {
	next, cmd := m.Update(PlayMsg{Track: track})
	next, cmd = next.(Model).Update(cmd())
	if cmd != nil {
		next, cmd = next.(Model).Update(cmd())
	}
	return next.(Model), cmd
}

// noticeLine is the view line under Search:.
func noticeLine(m Model) string { return lineAfterSearch(strings.Split(m.View().Content, "\n")) }

// START: TestResolveErrorDoesNotPlay

func TestResolveErrorDoesNotPlay(t *testing.T) {
	pl := &recordingPlayer{}
	playOne(New(&scriptedResolve{errs: []error{provider.ErrNetwork}}).WithPlayer(pl), provider.Track{ID: "a1", Title: "One"})
	if len(pl.urls) != 0 {
		t.Fatalf("Play calls = %v, want none", pl.urls)
	}
}

// END: TestResolveErrorDoesNotPlay

// START: TestResolveErrorNotice

func TestResolveErrorNotice(t *testing.T) {
	m, _ := playOne(New(&scriptedResolve{errs: []error{provider.ErrNetwork}}).WithPlayer(&recordingPlayer{}), provider.Track{ID: "a1", Title: "One"})
	if got, want := noticeLine(m), `Cannot play "One".`; got != want {
		t.Fatalf("line under Search: = %q, want %q", got, want)
	}
}

// END: TestResolveErrorNotice

// START: TestResolveTimeoutNotice

func TestResolveTimeoutNotice(t *testing.T) {
	m, _ := playOne(New(&scriptedResolve{errs: []error{provider.ErrTimeout}}).WithPlayer(&recordingPlayer{}), provider.Track{ID: "a1", Title: "One"})
	if got, want := noticeLine(m), `Cannot play "One".`; got != want {
		t.Fatalf("line under Search: = %q, want %q", got, want)
	}
}

// END: TestResolveTimeoutNotice

// START: TestPlayRecoversAfterResolveError

func TestPlayRecoversAfterResolveError(t *testing.T) {
	pl := &recordingPlayer{}
	m := New(&scriptedResolve{errs: []error{provider.ErrNetwork, nil}}).WithPlayer(pl)
	m, _ = playOne(m, provider.Track{ID: "a1", Title: "One"})
	m, _ = playOne(m, provider.Track{ID: "a2", Title: "Two"})
	if len(pl.urls) != 1 || pl.urls[0] != "http://stream/1" {
		t.Fatalf("Play calls = %v, want one with http://stream/1", pl.urls)
	}
	if got := noticeLine(m); got != "" {
		t.Fatalf("line under Search: = %q, want it gone", got)
	}
}

// END: TestPlayRecoversAfterResolveError

// START: TestArtistLineLoading

func TestArtistLineLoading(t *testing.T) {
	m := sized(New(&scriptedResolve{errs: []error{nil}, artist: "Daft Punk"}).WithPlayer(&recordingPlayer{}), 100)
	next, _ := m.Update(PlayMsg{Track: provider.Track{ID: "a1", Title: "One"}})
	if view := plain(next.(Model).View().Content); !strings.Contains(view, "Loading…") {
		t.Fatalf("view has no Loading… while the details load:\n%s", view)
	}
}

// END: TestArtistLineLoading

// START: TestArtistLineShowsArtist

func TestArtistLineShowsArtist(t *testing.T) {
	m := sized(New(&scriptedResolve{errs: []error{nil}, artist: "Daft Punk"}).WithPlayer(&recordingPlayer{}), 100)
	m, cmd := playOne(m, provider.Track{ID: "a1", Title: "One"})
	next, _ := m.Update(cmd())
	view := plain(next.(Model).View().Content)
	if !strings.Contains(view, "Daft Punk") || strings.Contains(view, "Loading…") {
		t.Fatalf("view should show Daft Punk and no Loading…:\n%s", view)
	}
}

// END: TestArtistLineShowsArtist
