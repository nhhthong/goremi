// Tests for the application model: focus, keys and view.
package app

import (
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: fakeProvider

// fakeProvider is a provider that returns nothing.
type fakeProvider struct{}

func (fakeProvider) Search(string, int) ([]provider.Track, error)     { return nil, nil }
func (fakeProvider) Details(t provider.Track) (provider.Track, error) { return t, nil }
func (fakeProvider) Resolve(provider.Track) (string, error)           { return "", nil }

// END: fakeProvider

// START: TestStartsOnSearchInput

func TestStartsOnSearchInput(t *testing.T) {
	if got := New(fakeProvider{}).Focus(); got != FocusInput {
		t.Fatalf("Focus() = %v, want FocusInput", got)
	}
}

// END: TestStartsOnSearchInput

// START: helpers

func key(r rune) tea.KeyPressMsg { return tea.KeyPressMsg{Code: r, Text: string(r)} }

// typed feeds the text to the model one key at a time and returns the model and the last command.
func typed(m Model, text string) (Model, tea.Cmd) {
	var cmd tea.Cmd
	for _, r := range text {
		var next tea.Model
		next, cmd = m.Update(key(r))
		m = next.(Model)
	}
	return m, cmd
}

// quits reports whether the command yields a quit message.
func quits(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

// END: helpers

// START: TestTypingAddsCharacters

func TestTypingAddsCharacters(t *testing.T) {
	m, _ := typed(New(fakeProvider{}), "daft")
	if m.Query() != "daft" {
		t.Fatalf("Query() = %q, want daft", m.Query())
	}
}

// END: TestTypingAddsCharacters

// START: TestQTypesInInput

func TestQTypesInInput(t *testing.T) {
	m, cmd := typed(New(fakeProvider{}), "q")
	if m.Query() != "q" || quits(cmd) {
		t.Fatalf("Query() = %q, quit = %v; want q and no quit", m.Query(), quits(cmd))
	}
}

// END: TestQTypesInInput

// START: TestSpaceTypesInInput

func TestSpaceTypesInInput(t *testing.T) {
	m, _ := typed(New(fakeProvider{}), "daft punk")
	if m.Query() != "daft punk" {
		t.Fatalf("Query() = %q, want daft punk", m.Query())
	}
}

// END: TestSpaceTypesInInput

// START: TestTypingMultiByte

func TestTypingMultiByte(t *testing.T) {
	m, _ := typed(New(fakeProvider{}), "bài")
	if m.Query() != "bài" {
		t.Fatalf("Query() = %q, want bài", m.Query())
	}
}

// END: TestTypingMultiByte

// START: TestEscInInputQuits

func TestEscInInputQuits(t *testing.T) {
	_, cmd := New(fakeProvider{}).Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if !quits(cmd) {
		t.Fatal("want a quit command")
	}
}

// END: TestEscInInputQuits

// START: TestCtrlCQuitsFromInput

func TestCtrlCQuitsFromInput(t *testing.T) {
	next, cmd := New(fakeProvider{}).Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	if !quits(cmd) || next.(Model).Query() != "" {
		t.Fatalf("quit = %v, Query() = %q; want quit and an empty query", quits(cmd), next.(Model).Query())
	}
}

// END: TestCtrlCQuitsFromInput

// START: TestCtrlCQuitsFromList

func TestCtrlCQuitsFromList(t *testing.T) {
	_, cmd := New(fakeProvider{}).WithFocus(FocusList).Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	if !quits(cmd) {
		t.Fatal("want a quit command")
	}
}

// END: TestCtrlCQuitsFromList

// START: TestViewHintOnFirstLine

func TestViewHintOnFirstLine(t *testing.T) {
	first := strings.SplitN(New(fakeProvider{}).View().Content, "\n", 2)[0]
	if !strings.Contains(first, "Ctrl+C") {
		t.Fatalf("first line = %q, want it to contain Ctrl+C", first)
	}
}

// END: TestViewHintOnFirstLine

// START: TestNewIsATeaModel

func TestNewIsATeaModel(t *testing.T) {
	var m tea.Model = New(fakeProvider{})
	if m.Init() != nil || m.View().Content == "" {
		t.Fatalf("Init() nil = %v, View empty = %v", m.Init() == nil, m.View().Content == "")
	}
}

// END: TestNewIsATeaModel

// START: recordingProvider

// recordingProvider remembers the Search calls it receives.
type recordingProvider struct {
	fakeProvider
	calls []searchCall
}

type searchCall struct {
	query string
	page  int
}

func (r *recordingProvider) Search(q string, page int) ([]provider.Track, error) {
	r.calls = append(r.calls, searchCall{q, page})
	return nil, nil
}

// END: recordingProvider

// START: TestBackspaceRemovesLast

func TestBackspaceRemovesLast(t *testing.T) {
	m, _ := typed(New(fakeProvider{}), "daft")
	next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
	if got := next.(Model).Query(); got != "daf" {
		t.Fatalf("Query() = %q, want daf", got)
	}
}

// END: TestBackspaceRemovesLast

// START: TestBackspaceOnEmpty

func TestBackspaceOnEmpty(t *testing.T) {
	next, _ := New(fakeProvider{}).Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
	if got := next.(Model).Query(); got != "" {
		t.Fatalf("Query() = %q, want empty", got)
	}
}

// END: TestBackspaceOnEmpty

// START: TestBackspaceMultiByte

func TestBackspaceMultiByte(t *testing.T) {
	m, _ := typed(New(fakeProvider{}), "bà")
	next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
	if got := next.(Model).Query(); got != "b" {
		t.Fatalf("Query() = %q, want b", got)
	}
}

// END: TestBackspaceMultiByte

// START: TestEnterRunsSearch

func TestEnterRunsSearch(t *testing.T) {
	p := &recordingProvider{}
	m, _ := typed(New(p), "daft punk")
	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("want a command")
	}
	cmd()
	if want := []searchCall{{"daft punk", 1}}; !reflect.DeepEqual(p.calls, want) {
		t.Fatalf("calls = %v, want %v", p.calls, want)
	}
}

// END: TestEnterRunsSearch

// START: scriptedProvider

// scriptedProvider answers each Search call with the next scripted step.
type scriptedProvider struct {
	fakeProvider
	steps []step
	next  int
}

type step struct {
	tracks []provider.Track
	err    error
}

func (s *scriptedProvider) Search(string, int) ([]provider.Track, error) {
	st := s.steps[s.next]
	s.next++
	return st.tracks, st.err
}

// END: scriptedProvider

// START: search helper

// search puts the focus on the input, types the text, presses Enter and feeds the search result back.
func search(m Model, text string) Model {
	m, _ = typed(m.WithFocus(FocusInput), text)
	next, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	next, _ = next.Update(cmd())
	return next.(Model)
}

func titles(ts []provider.Track) []string {
	var out []string
	for _, t := range ts {
		out = append(out, t.Title)
	}
	return out
}

// END: search helper

// START: TestEnterMovesFocusToList

func TestEnterMovesFocusToList(t *testing.T) {
	m, _ := typed(New(fakeProvider{}), "daft")
	next, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	next, _ = next.Update(cmd()) // the focus moves when the results arrive (18.3.1)
	if got := next.(Model).Focus(); got != FocusList {
		t.Fatalf("Focus() = %v, want FocusList", got)
	}
}

// END: TestEnterMovesFocusToList

// START: TestResultsShownInList

func TestResultsShownInList(t *testing.T) {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{{Title: "A"}, {Title: "B"}, {Title: "C"}}}}}
	m := search(New(p), "daft")
	if got, want := titles(m.Tracks()), []string{"A", "B", "C"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("tracks = %v, want %v", got, want)
	}
}

// END: TestResultsShownInList

// START: TestNewSearchReplacesList

func TestNewSearchReplacesList(t *testing.T) {
	p := &scriptedProvider{steps: []step{
		{tracks: []provider.Track{{Title: "A"}, {Title: "B"}}},
		{tracks: []provider.Track{{Title: "C"}}},
	}}
	m := search(search(New(p), "daft"), "x")
	if got, want := titles(m.Tracks()), []string{"C"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("tracks = %v, want %v", got, want)
	}
}

// END: TestNewSearchReplacesList

// START: TestEmptySearchClearsList

func TestEmptySearchClearsList(t *testing.T) {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{{Title: "A"}}}, {}}}
	m := search(search(New(p), "daft"), "x")
	if len(m.Tracks()) != 0 {
		t.Fatalf("tracks = %v, want none", titles(m.Tracks()))
	}
}

// END: TestEmptySearchClearsList

// START: TestSearchErrorMessage

func TestSearchErrorMessage(t *testing.T) {
	p := &scriptedProvider{steps: []step{{err: provider.ErrNetwork}}}
	view := search(New(p), "daft").View().Content
	for _, want := range []string{"Unable to search.", "Check your internet connection."} {
		if !strings.Contains(view, want) {
			t.Errorf("view %q misses %q", view, want)
		}
	}
}

// END: TestSearchErrorMessage

// START: TestSearchErrorKeepsRunning

func TestSearchErrorKeepsRunning(t *testing.T) {
	p := &scriptedProvider{steps: []step{{err: provider.ErrNetwork}}}
	m, _ := typed(New(p), "daft")
	next, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	next, cmd = next.Update(cmd())
	if quits(cmd) || next.(Model).Query() != "daft" {
		t.Fatalf("quit = %v, Query() = %q; want the app to keep running with the query", quits(cmd), next.(Model).Query())
	}
}

// END: TestSearchErrorKeepsRunning

// START: TestSearchRecovers

func TestSearchRecovers(t *testing.T) {
	p := &scriptedProvider{steps: []step{{err: provider.ErrNetwork}, {tracks: []provider.Track{{Title: "A"}}}}}
	m := search(search(New(p), "daft"), "x")
	if strings.Contains(m.View().Content, "Unable to search.") {
		t.Errorf("error text still shown: %q", m.View().Content)
	}
	if got := titles(m.Tracks()); !reflect.DeepEqual(got, []string{"A"}) {
		t.Errorf("tracks = %v, want A", got)
	}
}

// END: TestSearchRecovers

// START: TestTabInputToList

func TestTabInputToList(t *testing.T) {
	next, _ := New(fakeProvider{}).Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if got := next.(Model).Focus(); got != FocusList {
		t.Fatalf("Focus() = %v, want FocusList", got)
	}
}

// END: TestTabInputToList

// START: TestTabListToInputKeepsQuery

func TestTabListToInputKeepsQuery(t *testing.T) {
	m, _ := typed(New(fakeProvider{}), "daft")
	next, _ := m.WithFocus(FocusList).Update(tea.KeyPressMsg{Code: tea.KeyTab})
	got := next.(Model)
	if got.Focus() != FocusInput || got.Query() != "daft" {
		t.Fatalf("Focus() = %v, Query() = %q; want FocusInput and daft", got.Focus(), got.Query())
	}
}

// END: TestTabListToInputKeepsQuery

// START: TestTabIsNotTyped

func TestTabIsNotTyped(t *testing.T) {
	m, _ := typed(New(fakeProvider{}), "daft")
	next, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if next.(Model).Query() != "daft" || quits(cmd) {
		t.Fatalf("Query() = %q, quit = %v; want daft and no quit", next.(Model).Query(), quits(cmd))
	}
}

// END: TestTabIsNotTyped

// START: TestQInListQuits

func TestQInListQuits(t *testing.T) {
	_, cmd := New(fakeProvider{}).WithFocus(FocusList).Update(key('q'))
	if !quits(cmd) {
		t.Fatal("want a quit command")
	}
}

// END: TestQInListQuits

// START: TestEscInListReturnsToInput

func TestEscInListReturnsToInput(t *testing.T) {
	m, _ := typed(New(fakeProvider{}), "daft")
	next, _ := m.WithFocus(FocusList).Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	got := next.(Model)
	if got.Focus() != FocusInput || got.Query() != "daft" {
		t.Fatalf("Focus() = %v, Query() = %q; want FocusInput and daft", got.Focus(), got.Query())
	}
}

// END: TestEscInListReturnsToInput

// START: TestEscInListDoesNotQuit

func TestEscInListDoesNotQuit(t *testing.T) {
	_, cmd := New(fakeProvider{}).WithFocus(FocusList).Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if quits(cmd) {
		t.Fatal("Esc in the list must not quit")
	}
}

// END: TestEscInListDoesNotQuit

// START: TestViewInputAboveList

func TestViewInputAboveList(t *testing.T) {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{{Title: "A"}, {Title: "B"}}}}}
	lines := strings.Split(search(New(p), "daft").View().Content, "\n")
	input, list := -1, -1
	for i, l := range lines {
		if strings.Contains(l, "Search:") && input < 0 {
			input = i
		}
		if strings.Contains(l, "A") && strings.Contains(l, "▶") && list < 0 {
			list = i
		}
	}
	if input < 0 || list < 0 || input >= list {
		t.Fatalf("search line %d, first track line %d; want search above list in %q", input, list, lines)
	}
}

// END: TestViewInputAboveList

// START: TestViewShowsAllTracks

func TestViewShowsAllTracks(t *testing.T) {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{{Title: "A"}, {Title: "B"}, {Title: "C"}}}}}
	c := search(New(p), "daft").View().Content
	last := -1
	for _, s := range []string{"A", "B", "C", "load more..."} {
		i := strings.Index(c[strings.Index(c, "Search:"):], s)
		if i < 0 || i <= last {
			t.Fatalf("%q missing or out of order in %q", s, c)
		}
		last = i
	}
}

// END: TestViewShowsAllTracks
