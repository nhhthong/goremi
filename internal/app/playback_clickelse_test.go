// Tests for clicks that must do nothing: off the controls, between them, or with another button.
package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// clickAt left-clicks (or with the given button) a cell of a model that plays B in the list A, B, C and returns the fake player and the command.
func clickAt(t *testing.T, button tea.MouseButton, cell func(m Model) (int, int)) (*recordingPlayer, tea.Cmd) {
	t.Helper()
	pl := &recordingPlayer{}
	m := playingModelWith(pl, trackB)
	x, y := cell(m)
	_, cmd := m.Update(tea.MouseClickMsg{X: x, Y: y, Button: button})
	return pl, cmd
}

// quiet fails when the click moved the player or gave a command.
func quiet(t *testing.T, pl *recordingPlayer, cmd tea.Cmd) {
	t.Helper()
	if pl.toggles != 0 || len(pl.seeks) != 0 || cmd != nil {
		t.Fatalf("toggles %d, seeks %v, command %v, want nothing", pl.toggles, pl.seeks, cmd != nil)
	}
}

// START: TestClickElsewhereDoesNothing

func TestClickElsewhereDoesNothing(t *testing.T) {
	pl, cmd := clickAt(t, tea.MouseLeft, func(Model) (int, int) { return 0, 0 })
	quiet(t, pl, cmd)
}

// END: TestClickElsewhereDoesNothing

// START: TestRightClickDoesNothing

func TestRightClickDoesNothing(t *testing.T) {
	pl, cmd := clickAt(t, tea.MouseRight, func(m Model) (int, int) { return controlCell(t, m, 2) })
	quiet(t, pl, cmd)
}

// END: TestRightClickDoesNothing

// START: TestClickBetweenControlsDoesNothing

func TestClickBetweenControlsDoesNothing(t *testing.T) {
	pl, cmd := clickAt(t, tea.MouseLeft, func(m Model) (int, int) {
		x, y := controlCell(t, m, 0)
		return x + 2, y // the blank cell between ◀◀ and ◀
	})
	quiet(t, pl, cmd)
}

// END: TestClickBetweenControlsDoesNothing
