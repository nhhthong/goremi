// Tests for the focus change after a search and the Enter key in the results list.
package app

import (
	"errors"
	"reflect"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: helpers

// pagedProvider answers Search from a page map and records every call.
type pagedProvider struct {
	fakeProvider
	pages map[int][]provider.Track
	calls []searchCall
}

func (p *pagedProvider) Search(q string, page int) ([]provider.Track, error) {
	p.calls = append(p.calls, searchCall{q, page})
	return p.pages[page], nil
}

// listOf returns a model with the list A, B, C for the query "daft", the focus on the list and the selected line set.
func listOf(selected int) (Model, *pagedProvider) {
	p := &pagedProvider{pages: map[int][]provider.Track{
		1: {{Title: "A"}, {Title: "B"}, {Title: "C"}},
		2: {{Title: "D"}, {Title: "E"}},
	}}
	m := search(New(p), "daft").WithFocus(FocusList)
	p.calls = nil
	for i := 0; i < selected; i++ {
		next, _ := m.Update(down)
		m = next.(Model)
	}
	return m, p
}

// pressEnterOnList sends Enter to the model and returns the message its command yields, or nil.
func pressEnterOnList(m Model) (Model, tea.Msg) {
	next, cmd := m.Update(enter)
	if cmd == nil {
		return next.(Model), nil
	}
	return next.(Model), cmd()
}

// END: helpers

// START: focus after a search

func TestFocusToListOnResults(t *testing.T) {
	m := search(New(fakeProvider{}), "daft")
	if m.Focus() != FocusList {
		t.Fatalf("Focus() = %v, want FocusList", m.Focus())
	}
}

func TestFocusStaysOnInputAfterError(t *testing.T) {
	m := search(New(&scriptedProvider{steps: []step{{err: errors.New("boom")}}}), "daft")
	if m.Focus() != FocusInput {
		t.Fatalf("Focus() = %v, want FocusInput", m.Focus())
	}
}

func TestFocusWaitsForResults(t *testing.T) {
	m, _ := typed(New(fakeProvider{}), "daft")
	next, _ := m.Update(enter)
	if got := next.(Model).Focus(); got != FocusInput {
		t.Fatalf("Focus() = %v, want FocusInput", got)
	}
}

// END: focus after a search

// START: load more

func TestLoadMoreCallsPage2(t *testing.T) {
	m, p := listOf(3)
	pressEnterOnList(m)
	if want := []searchCall{{"daft", 2}}; !reflect.DeepEqual(p.calls, want) {
		t.Fatalf("calls = %v, want %v", p.calls, want)
	}
}

func TestLoadMoreAppendsTracks(t *testing.T) {
	m, _ := listOf(3)
	next, msg := pressEnterOnList(m)
	after, _ := next.Update(msg)
	if got, want := titles(after.(Model).Tracks()), []string{"A", "B", "C", "D", "E"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("tracks = %v, want %v", got, want)
	}
}

func TestEnterOnTrackDoesNotLoadMore(t *testing.T) {
	m, p := listOf(1)
	pressEnterOnList(m)
	if len(p.calls) != 0 {
		t.Fatalf("Search called %v, want no call", p.calls)
	}
}

// END: load more

// START: play

func TestEnterPlaysSelectedTrack(t *testing.T) {
	m, _ := listOf(1)
	_, msg := pressEnterOnList(m)
	if got, ok := msg.(PlayMsg); !ok || got.Track.Title != "B" {
		t.Fatalf("message = %#v, want PlayMsg for B", msg)
	}
}

func TestEnterPlaysFirstTrack(t *testing.T) {
	m, _ := listOf(0)
	_, msg := pressEnterOnList(m)
	if got, ok := msg.(PlayMsg); !ok || got.Track.Title != "A" {
		t.Fatalf("message = %#v, want PlayMsg for A", msg)
	}
}

func TestEnterOnLoadMoreDoesNotPlay(t *testing.T) {
	m, _ := listOf(3)
	_, msg := pressEnterOnList(m)
	if _, ok := msg.(PlayMsg); ok {
		t.Fatalf("message = %#v, want no PlayMsg", msg)
	}
}

func TestEnterOnEmptyListDoesNotPlay(t *testing.T) {
	_, msg := pressEnterOnList(New(fakeProvider{}).WithFocus(FocusList))
	if _, ok := msg.(PlayMsg); ok {
		t.Fatalf("message = %#v, want no PlayMsg", msg)
	}
}

// END: play
