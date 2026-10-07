// Tests of the fixed search bar and the scrolling list (ui tasks 3.8.1 to 3.8.5).
package app

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: helpers

// pageAt is a model of the given size holding one full page of tracks T0 to T9 (so a `load more...` line follows); the list has the focus and, when played, T0 plays.
func pageAt(width, height int, played bool) Model {
	var tracks []provider.Track
	for i := 0; i < provider.PageSize; i++ {
		tracks = append(tracks, provider.Track{ID: fmt.Sprintf("id%d", i), Title: fmt.Sprintf("T%d", i)})
	}
	m := search(sizeTo(New(&scriptedProvider{steps: []step{{tracks: tracks}}}).WithPlayer(&recordingPlayer{}), width, height), "daft")
	if played {
		next, _ := m.Update(PlayMsg{Track: tracks[0]})
		m = next.(Model)
	}
	return m
}

// downs presses Down n times.
func downs(m Model, n int) Model {
	for i := 0; i < n; i++ {
		next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
		m = next.(Model)
	}
	return m
}

// ups presses Up n times.
func ups(m Model, n int) Model {
	for i := 0; i < n; i++ {
		next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
		m = next.(Model)
	}
	return m
}

// endsWithBar checks that the last three lines are a rule, the prompt and a rule.
func endsWithBar(t *testing.T, lines []string) {
	t.Helper()
	n := len(lines)
	if n < 3 || !strings.HasPrefix(lines[n-3], "─") || !strings.HasPrefix(lines[n-2], "❯") || !strings.HasPrefix(lines[n-1], "─") {
		t.Fatalf("the last three lines are %q, want the search bar", lines[max(n-3, 0):])
	}
}

// END: helpers

// START: TestViewKeepsHeightWithLongList

func TestViewKeepsHeightWithLongList(t *testing.T) {
	for _, played := range []bool{false, true} {
		lines := plainLines(pageAt(100, 25, played))
		if len(lines) != 25 {
			t.Errorf("played %v: %d lines, want 25", played, len(lines))
		}
		endsWithBar(t, lines)
	}
}

// END: TestViewKeepsHeightWithLongList

// START: TestListShowsOnlyTheRowsLeft

func TestListShowsOnlyTheRowsLeft(t *testing.T) {
	got := plain(pageAt(100, 25, false).View().Content)
	for _, want := range []string{"▶ T0", "  T1", "  T2"} {
		if !strings.Contains(got, want) {
			t.Errorf("view lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "T3") {
		t.Fatalf("view shows T3, want 3 rows (25 lines − 3 bar − 1 hint − 12 mascot − 6 banner):\n%s", got)
	}
}

// END: TestListShowsOnlyTheRowsLeft

// START: TestScrollDownKeepsSelectionVisible

func TestScrollDownKeepsSelectionVisible(t *testing.T) {
	got := plain(downs(pageAt(100, 25, false), 10).View().Content)
	for _, want := range []string{"  T8", "  T9", "▶ load more..."} {
		if !strings.Contains(got, want) {
			t.Errorf("view lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "T7") || strings.Contains(got, "T0") {
		t.Fatalf("view still shows an old row:\n%s", got)
	}
}

// END: TestScrollDownKeepsSelectionVisible

// START: TestScrollUpOnlyAtTheTop

func TestScrollUpOnlyAtTheTop(t *testing.T) {
	m := ups(downs(pageAt(100, 25, false), 10), 2) // selection on T8, still inside the window
	if got := plain(m.View().Content); !strings.Contains(got, "▶ T8") || !strings.Contains(got, "load more...") || strings.Contains(got, "T7") {
		t.Fatalf("after two Up the window moved:\n%s", got)
	}
	m = ups(m, 1) // selection on T7, above the window
	if got := plain(m.View().Content); !strings.Contains(got, "▶ T7") || !strings.Contains(got, "  T9") || strings.Contains(got, "load more...") || strings.Contains(got, "T6") {
		t.Fatalf("after three Up the window did not scroll up by one line:\n%s", got)
	}
}

// END: TestScrollUpOnlyAtTheTop

// START: TestContentIsCutWhenNothingIsLeft

func TestContentIsCutWhenNothingIsLeft(t *testing.T) {
	lines := plainLines(pageAt(60, 12, true))
	if len(lines) != 12 {
		t.Fatalf("%d lines, want 12", len(lines))
	}
	endsWithBar(t, lines)
}

// END: TestContentIsCutWhenNothingIsLeft
