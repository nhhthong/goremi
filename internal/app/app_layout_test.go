// Tests for the layout: panel beside the list from 80 columns, stacked below.
package app

import (
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
	"goremi/internal/ui"
)

// START: helpers

// artLines are the lines of the Goremi art.
func artLines() []string { return strings.Split(ui.Logo, "\n") }

// art0 is the first art line without its outer spaces, a fragment that identifies the art in a view line.
func art0() string { return strings.TrimSpace(artLines()[0]) }

// sized sends a WindowSizeMsg of the given width to the model.
func sized(m Model, width int) Model {
	next, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: 40})
	return next.(Model)
}

// abModel returns a model of the given width whose list holds the tracks A and B.
func abModel(width int) Model {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{{Title: "A"}, {Title: "B"}}}}}
	return search(sized(New(p), width), "daft")
}

// sideBySide reports whether one view line holds both the selected track A and the first art line.
func sideBySide(m Model) bool {
	for _, l := range viewLines(m) {
		if strings.Contains(l, "▶ A") && strings.Contains(l, art0()) {
			return true
		}
	}
	return false
}

// lineIndex returns the index of the first line for which match is true, or -1.
func lineIndex(lines []string, match func(string) bool) int {
	for i, l := range lines {
		if match(l) {
			return i
		}
	}
	return -1
}

// END: helpers

// START: wide layout

func TestWideListLeftPanelRight(t *testing.T) {
	if !sideBySide(abModel(100)) {
		t.Fatalf("want the list and the art on one line in %q", abModel(100).View().Content)
	}
}

func TestWidth80IsSideBySide(t *testing.T) {
	if !sideBySide(abModel(80)) {
		t.Fatalf("want the list and the art on one line at 80 columns")
	}
}

func TestWidePanelWithoutTracks(t *testing.T) {
	if got := sized(New(fakeProvider{}), 100).View().Content; !strings.Contains(got, art0()) {
		t.Fatalf("view %q misses the art", got)
	}
}

func TestNoPanelBeforeSize(t *testing.T) {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{{Title: "A"}, {Title: "B"}}}}}
	if got := search(New(p), "daft").View().Content; strings.Contains(got, art0()) {
		t.Fatalf("view %q shows the panel before any size is known", got)
	}
}

// END: wide layout

// START: stacked layout

func isSearch(l string) bool { return strings.HasPrefix(l, "Search:") }
func isArt(l string) bool    { return strings.Contains(l, art0()) }
func isTrackA(l string) bool { return strings.Contains(l, "▶ A") }

func TestNarrowStacksSearchPanelList(t *testing.T) {
	lines := viewLines(abModel(79))
	s, a, l := lineIndex(lines, isSearch), lineIndex(lines, isArt), lineIndex(lines, isTrackA)
	if s < 0 || !(s < a && a < l) {
		t.Fatalf("search line %d, art line %d, list line %d; want search < art < list in %q", s, a, l, lines)
	}
}

func TestNarrowIsNotSideBySide(t *testing.T) {
	if sideBySide(abModel(79)) {
		t.Fatalf("want the art and the list on separate lines at 79 columns")
	}
}

func TestResizeSwitchesLayout(t *testing.T) {
	m := sized(abModel(100), 79)
	lines := viewLines(m)
	s, a, l := lineIndex(lines, isSearch), lineIndex(lines, isArt), lineIndex(lines, isTrackA)
	if sideBySide(m) || !(s < a && a < l) {
		t.Fatalf("want a stacked view after the resize, got %q", lines)
	}
}

func TestNarrowNoTracksOnlyArt(t *testing.T) {
	lines := viewLines(sized(New(fakeProvider{}), 60))
	s := lineIndex(lines, isSearch)
	if got := lines[s+1:]; !reflect.DeepEqual(got, artLines()) {
		t.Fatalf("lines after Search: = %q, want only the art %q", got, artLines())
	}
}

func TestNarrowZeroResultsOnlyArt(t *testing.T) {
	got := strings.Join(viewLines(search(sized(New(fakeProvider{}), 60), "daft")), "\n")
	if strings.Contains(got, "load more...") || !strings.Contains(got, art0()) {
		t.Fatalf("view %q: want the art and no load more line", got)
	}
}

// END: stacked layout
