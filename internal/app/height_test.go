// Tests that the model keeps the terminal height from the WindowSizeMsg (task 3.5.2).
package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// START: sizeTo

// sizeTo sends a WindowSizeMsg of the given size to the model.
func sizeTo(m Model, width, height int) Model {
	next, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: height})
	return next.(Model)
}

// END: sizeTo

// START: TestHeightStoredFromWindowSize

func TestHeightStoredFromWindowSize(t *testing.T) {
	if got := sizeTo(New(fakeProvider{}), 100, 30).Height(); got != 30 {
		t.Fatalf("Height() = %d, want 30", got)
	}
}

// END: TestHeightStoredFromWindowSize

// START: TestHeightIsZeroBeforeSize

func TestHeightIsZeroBeforeSize(t *testing.T) {
	if got := New(fakeProvider{}).Height(); got != 0 {
		t.Fatalf("Height() = %d before any WindowSizeMsg, want 0", got)
	}
}

// END: TestHeightIsZeroBeforeSize

// START: TestHeightFollowsResize

func TestHeightFollowsResize(t *testing.T) {
	m := sizeTo(sizeTo(New(fakeProvider{}), 100, 30), 100, 20)
	if got := m.Height(); got != 20 {
		t.Fatalf("Height() = %d after a resize to 20, want 20", got)
	}
}

// END: TestHeightFollowsResize
