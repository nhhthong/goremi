// Tests for the results model: Enter on the "load more..." line loads the next page.
package ui

import (
	"errors"
	"reflect"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: fakeProvider

// fakeProvider answers Search from a page map and records every call.
type fakeProvider struct {
	pages map[int][]provider.Track
	err   error
	calls [][2]any
}

func (f *fakeProvider) Search(query string, page int) ([]provider.Track, error) {
	f.calls = append(f.calls, [2]any{query, page})
	return f.pages[page], f.err
}

func (f *fakeProvider) Details(t provider.Track) (provider.Track, error) { return t, nil }
func (f *fakeProvider) Resolve(t provider.Track) (string, error)         { return "", nil }

// END: fakeProvider

// START: helpers

func titles(ts []provider.Track) []string {
	var out []string
	for _, t := range ts {
		out = append(out, t.Title)
	}
	return out
}

// pressEnter sends Enter and, when a command comes back, runs it and feeds the result back.
func pressEnter(r Results) (Results, tea.Cmd) {
	r, cmd := r.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd != nil {
		r, _ = r.Update(cmd())
	}
	return r, cmd
}

// END: helpers

// START: TestLoadMoreAppends

func TestLoadMoreAppends(t *testing.T) {
	f := &fakeProvider{pages: map[int][]provider.Track{2: {{Title: "K"}, {Title: "L"}}}}
	r := NewResults(f, "daft punk", tenTracks()).Select(10)
	r, cmd := pressEnter(r)
	if cmd == nil {
		t.Fatal("want a command")
	}
	if want := [][2]any{{"daft punk", 2}}; !reflect.DeepEqual(f.calls, want) {
		t.Fatalf("calls = %v, want %v", f.calls, want)
	}
	if got, want := titles(r.Tracks()), append(titles(tenTracks()), "K", "L"); !reflect.DeepEqual(got, want) {
		t.Fatalf("tracks = %v, want %v", got, want)
	}
}

// END: TestLoadMoreAppends

// START: TestEnterOnTrackDoesNotLoad

func TestEnterOnTrackDoesNotLoad(t *testing.T) {
	f := &fakeProvider{}
	r := NewResults(f, "q", []provider.Track{{Title: "A"}, {Title: "B"}, {Title: "C"}}).Select(1)
	r, cmd := pressEnter(r)
	if cmd != nil || len(f.calls) != 0 {
		t.Fatalf("want no command and no Search, got cmd=%v calls=%v", cmd != nil, f.calls)
	}
	if got := titles(r.Tracks()); !reflect.DeepEqual(got, []string{"A", "B", "C"}) {
		t.Fatalf("tracks = %v", got)
	}
}

// END: TestEnterOnTrackDoesNotLoad

// START: TestLoadMoreSecondPage

func TestLoadMoreSecondPage(t *testing.T) {
	f := &fakeProvider{pages: map[int][]provider.Track{2: numbered("K", 10), 3: {{Title: "E"}}}}
	r := NewResults(f, "q", tenTracks()).Select(10)
	r, _ = pressEnter(r)
	r, _ = pressEnter(r.Select(len(r.Tracks())))
	if want := [][2]any{{"q", 2}, {"q", 3}}; !reflect.DeepEqual(f.calls, want) {
		t.Fatalf("calls = %v, want %v", f.calls, want)
	}
	if got := titles(r.Tracks()); len(got) != 21 || got[20] != "E" {
		t.Fatalf("tracks = %v, want 21 tracks ending with E", got)
	}
}

// END: TestLoadMoreSecondPage

// START: TestLoadMoreError

func TestLoadMoreError(t *testing.T) {
	f := &fakeProvider{err: errors.New("boom")}
	r := NewResults(f, "q", tenTracks()).Select(10)
	r, _ = pressEnter(r)
	if got := titles(r.Tracks()); !reflect.DeepEqual(got, titles(tenTracks())) {
		t.Fatalf("tracks = %v, want unchanged", got)
	}
}

// END: TestLoadMoreError
