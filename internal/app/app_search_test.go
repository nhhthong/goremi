// Tests for the search input rules and the list keys of the application model.
package app

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: helpers

var enter = tea.KeyPressMsg{Code: tea.KeyEnter}

// searchedWith types the text, presses Enter, runs the command and returns the provider's calls.
func searchedWith(text string) []searchCall {
	p := &recordingProvider{}
	m, _ := typed(New(p), text)
	_, cmd := m.Update(enter)
	cmd()
	return p.calls
}

// abcList returns a model whose list holds the tracks A, B and C, with the focus on the list.
func abcList() Model {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{{Title: "A"}, {Title: "B"}, {Title: "C"}}}}}
	return search(New(p), "daft").WithFocus(FocusList)
}

// END: helpers

// START: query trimming

func TestEnterTrimsQuery(t *testing.T) {
	if got, want := searchedWith("  daft punk  "), []searchCall{{"daft punk", 1}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("calls = %v, want %v", got, want)
	}
}

func TestEnterTrimsLeading(t *testing.T) {
	if got, want := searchedWith("  daft"), []searchCall{{"daft", 1}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("calls = %v, want %v", got, want)
	}
}

func TestEnterTrimsTrailing(t *testing.T) {
	if got, want := searchedWith("daft  "), []searchCall{{"daft", 1}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("calls = %v, want %v", got, want)
	}
}

func TestEnterKeepsInnerSpaces(t *testing.T) {
	if got, want := searchedWith("daft  punk"), []searchCall{{"daft  punk", 1}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("calls = %v, want %v", got, want)
	}
}

// END: query trimming

// START: overlapping searches

func TestSecondEnterWhileSearching(t *testing.T) {
	p := &recordingProvider{}
	m, _ := typed(New(p), "daft")
	next, cmd1 := m.Update(enter)
	_, cmd2 := next.(Model).WithFocus(FocusInput).Update(enter)
	cmd1()
	if cmd2 != nil {
		cmd2()
	}
	if len(p.calls) != 1 {
		t.Fatalf("Search called %d times, want 1", len(p.calls))
	}
}

func TestEnterAllowedAfterResult(t *testing.T) {
	p := &recordingProvider{}
	m, _ := typed(New(p), "daft")
	next, cmd1 := m.Update(enter)
	next, _ = next.Update(cmd1())
	_, cmd2 := next.(Model).WithFocus(FocusInput).Update(enter)
	if cmd2 == nil {
		t.Fatal("want a command for the second search")
	}
	cmd2()
	if len(p.calls) != 2 {
		t.Fatalf("Search called %d times, want 2", len(p.calls))
	}
}

// END: overlapping searches

// START: TestFailedSearchKeepsList

func TestFailedSearchKeepsList(t *testing.T) {
	p := &scriptedProvider{steps: []step{
		{tracks: []provider.Track{{Title: "A"}, {Title: "B"}}},
		{err: errors.New("offline")},
	}}
	m := search(search(New(p), "daft"), "punk")
	if got, want := titles(m.Tracks()), []string{"A", "B"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("tracks = %v, want %v", got, want)
	}
}

// END: TestFailedSearchKeepsList

// START: list selection keys

var (
	down = tea.KeyPressMsg{Code: tea.KeyDown}
	up   = tea.KeyPressMsg{Code: tea.KeyUp}
)

func TestDownInListMovesSelection(t *testing.T) {
	next, _ := abcList().Update(down)
	if got := next.(Model).Selected(); got != 1 {
		t.Fatalf("Selected() = %d, want 1", got)
	}
}

func TestUpInListMovesSelection(t *testing.T) {
	next, _ := abcList().Update(down)
	next, _ = next.Update(up)
	if got := next.(Model).Selected(); got != 0 {
		t.Fatalf("Selected() = %d, want 0", got)
	}
}

func TestDownInInputIgnored(t *testing.T) {
	// abcList searched for "daft", so the input already holds that query.
	next, _ := abcList().WithFocus(FocusInput).Update(down)
	got := next.(Model)
	if got.Selected() != 0 || got.Query() != "daft" {
		t.Fatalf("Selected() = %d, Query() = %q; want 0 and daft", got.Selected(), got.Query())
	}
}

func TestViewFollowsSelection(t *testing.T) {
	next, _ := abcList().Update(down)
	c := next.(Model).View().Content
	if !strings.Contains(c, "▶ B") || strings.Contains(c, "▶ A") {
		t.Fatalf("want the marker on B only in %q", c)
	}
}

// END: list selection keys

// START: hint style

// faint is the start of the style code of a faint line; the colour joins it, as in "\x1b[2;38;2;…m".
const faint = "\x1b[2;"

func TestHintIsFaint(t *testing.T) {
	first := strings.Split(New(fakeProvider{}).View().Content, "\n")[0]
	if !strings.Contains(first, "Ctrl+C: quit") || !strings.Contains(first, faint) {
		t.Fatalf("first line %q: want the hint with the faint code", first)
	}
}

func TestSearchLineNotFaint(t *testing.T) {
	second := strings.Split(New(fakeProvider{}).View().Content, "\n")[1]
	if strings.Contains(second, faint) {
		t.Fatalf("second line %q must not be faint", second)
	}
}

// END: hint style

// START: empty query

// enterOn types the text, presses Enter, runs the returned command if any and returns the model and the provider.
func enterOn(text string) (Model, *recordingProvider) {
	p := &recordingProvider{}
	m, _ := typed(New(p), text)
	next, cmd := m.Update(enter)
	if cmd != nil {
		cmd()
	}
	return next.(Model), p
}

func TestEnterOnEmptyDoesNotSearch(t *testing.T) {
	if _, p := enterOn(""); len(p.calls) != 0 {
		t.Fatalf("Search called %v, want no call", p.calls)
	}
}

func TestEnterOnSpacesDoesNotSearch(t *testing.T) {
	if _, p := enterOn("   "); len(p.calls) != 0 {
		t.Fatalf("Search called %v, want no call", p.calls)
	}
}

func TestEnterOnSpacesKeepsFocus(t *testing.T) {
	if m, _ := enterOn("   "); m.Focus() != FocusInput {
		t.Fatalf("Focus() = %v, want FocusInput", m.Focus())
	}
}

func TestEnterAfterEmptyStillSearches(t *testing.T) {
	m, p := enterOn("")
	m, _ = typed(m, "daft")
	_, cmd := m.Update(enter)
	if cmd == nil {
		t.Fatal("want a command for the search")
	}
	cmd()
	if want := []searchCall{{"daft", 1}}; !reflect.DeepEqual(p.calls, want) {
		t.Fatalf("calls = %v, want %v", p.calls, want)
	}
}

// END: empty query
