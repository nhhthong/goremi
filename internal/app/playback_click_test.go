// Tests for the left click on the controls of the player panel.
package app

import (
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
)

// START: helpers

// controlCell is the cell of the first rune of control n (0 previous … 4 next) in the drawn view of m.
func controlCell(t *testing.T, m Model, n int) (x, y int) {
	t.Helper()
	lines := plainLines(m)
	y = lineWith(lines, "◀◀")
	if y < 0 {
		t.Fatalf("no controls row in:\n%s", strings.Join(lines, "\n"))
	}
	i := strings.Index(lines[y], controlsRow)
	x = utf8.RuneCountInString(lines[y][:i])
	for _, f := range strings.Split(controlsRow, " ")[:n] {
		x += utf8.RuneCountInString(f) + 1
	}
	return x, y
}

// clickOn plays B in the list A, B, C and left-clicks control n; it returns the fake player and the command.
func clickOn(t *testing.T, n int) (*recordingPlayer, tea.Cmd) {
	t.Helper()
	pl := &recordingPlayer{}
	m := playingModelWith(pl, trackB)
	x, y := controlCell(t, m, n)
	_, cmd := m.Update(tea.MouseClickMsg{X: x, Y: y, Button: tea.MouseLeft})
	return pl, cmd
}

// END: helpers

// START: TestClickPrevious

func TestClickPrevious(t *testing.T) {
	_, cmd := clickOn(t, 0)
	if cmd == nil || cmd() != (PlayMsg{Track: trackA}) {
		t.Fatal("a click on ◀◀ did not give a PlayMsg for A")
	}
}

// END: TestClickPrevious

// START: TestClickSeekBack

func TestClickSeekBack(t *testing.T) {
	if pl, _ := clickOn(t, 1); !reflect.DeepEqual(pl.seeks, []float64{-10}) {
		t.Fatalf("Seek calls after a click on ◀ = %v, want [-10]", pl.seeks)
	}
}

// END: TestClickSeekBack

// START: TestClickPause

func TestClickPause(t *testing.T) {
	if pl, _ := clickOn(t, 2); pl.toggles != 1 {
		t.Fatalf("TogglePause calls after a click on ❚❚ = %d, want 1", pl.toggles)
	}
}

// END: TestClickPause

// START: TestClickSeekForward

func TestClickSeekForward(t *testing.T) {
	if pl, _ := clickOn(t, 3); !reflect.DeepEqual(pl.seeks, []float64{10}) {
		t.Fatalf("Seek calls after a click on ▶ = %v, want [10]", pl.seeks)
	}
}

// END: TestClickSeekForward

// START: TestClickNext

func TestClickNext(t *testing.T) {
	_, cmd := clickOn(t, 4)
	if cmd == nil || cmd() != (PlayMsg{Track: trackC}) {
		t.Fatal("a click on ▶▶ did not give a PlayMsg for C")
	}
}

// END: TestClickNext
