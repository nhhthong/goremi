// Tests for the results model: ↑ and ↓ move the selection.
package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// START: helpers

func abc() []provider.Track { return []provider.Track{{Title: "A"}, {Title: "B"}, {Title: "C"}} }

// pressed sends one key to a list with the given tracks and selected line and returns the selected line.
func pressed(tracks []provider.Track, selected int, code rune) int {
	r, _ := NewResults(&fakeProvider{}, "q", tracks).Select(selected).Update(tea.KeyPressMsg{Code: code})
	return r.Selected()
}

// END: helpers

// START: TestDownMovesSelection

func TestDownMovesSelection(t *testing.T) {
	if got := pressed(abc(), 0, tea.KeyDown); got != 1 {
		t.Fatalf("Selected() = %d, want 1", got)
	}
}

// END: TestDownMovesSelection

// START: TestUpMovesSelection

func TestUpMovesSelection(t *testing.T) {
	if got := pressed(abc(), 1, tea.KeyUp); got != 0 {
		t.Fatalf("Selected() = %d, want 0", got)
	}
}

// END: TestUpMovesSelection

// START: TestUpStopsAtFirst

func TestUpStopsAtFirst(t *testing.T) {
	if got := pressed(abc(), 0, tea.KeyUp); got != 0 {
		t.Fatalf("Selected() = %d, want 0", got)
	}
}

// END: TestUpStopsAtFirst

// START: TestDownReachesLoadMore

func TestDownReachesLoadMore(t *testing.T) {
	if got := pressed(abc(), 2, tea.KeyDown); got != 3 {
		t.Fatalf("Selected() = %d, want 3", got)
	}
}

// END: TestDownReachesLoadMore

// START: TestDownStopsAtLoadMore

func TestDownStopsAtLoadMore(t *testing.T) {
	if got := pressed(abc(), 3, tea.KeyDown); got != 3 {
		t.Fatalf("Selected() = %d, want 3", got)
	}
}

// END: TestDownStopsAtLoadMore

// START: TestDownOnEmptyList

func TestDownOnEmptyList(t *testing.T) {
	if got := pressed(nil, 0, tea.KeyDown); got != 0 {
		t.Fatalf("Selected() = %d, want 0", got)
	}
}

// END: TestDownOnEmptyList
