// Tests of the click and the drag on the progress bar of the player panel (tasks 5.32 to 5.36.1).
package app

import (
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

// barCells is the width of the bar in cells: the 40-column panel minus 2.
const barCells = 38

// START: barModel

// barModel plays a track of the given length at width 100 and returns the model and its fake player; the bar starts at 0 s.
func barModel(length time.Duration) (Model, *recordingPlayer) {
	pl := &recordingPlayer{}
	track := provider.Track{ID: "a1", Title: "One", Duration: length}
	m := sized(New(&scriptedResolve{errs: []error{nil}, artist: "Daft Punk"}).WithPlayer(pl), 100)
	m, cmd := playOne(m, track)
	next, _ := m.Update(cmd())
	return next.(Model), pl
}

// barCell is the column of the first cell of the bar and its row in the drawn view of m.
func barCell(t *testing.T, m Model) (x, y int) {
	t.Helper()
	lines := plainLines(m)
	y = lineWith(lines, "●")
	if y < 0 {
		t.Fatalf("no bar in:\n%s", strings.Join(lines, "\n"))
	}
	return utf8.RuneCountInString(lines[y][:strings.Index(lines[y], "●")]), y
}

// press, move and release send the three mouse messages of a drag at the given cell.
func press(m Model, x, y int) Model {
	next, _ := m.Update(tea.MouseClickMsg{X: x, Y: y, Button: tea.MouseLeft})
	return next.(Model)
}

func move(m Model, x, y int) Model {
	next, _ := m.Update(tea.MouseMotionMsg{X: x, Y: y, Button: tea.MouseLeft})
	return next.(Model)
}

func release(m Model, x, y int) Model {
	next, _ := m.Update(tea.MouseReleaseMsg{X: x, Y: y, Button: tea.MouseLeft})
	return next.(Model)
}

// clickCell presses and releases on cell i of the bar.
func clickCell(t *testing.T, m Model, i int) Model {
	t.Helper()
	x, y := barCell(t, m)
	return release(press(m, x+i, y), x+i, y)
}

// END: barModel

// START: TestMpvPlayerSeekToBeforePlayDoesNothing

func TestMpvPlayerSeekToBeforePlayDoesNothing(t *testing.T) {
	if err := newMpvPlayer().SeekTo(5 * time.Second); err != nil {
		t.Fatalf("SeekTo before the first Play returned %v, want nil (nothing happens)", err)
	}
}

// END: TestMpvPlayerSeekToBeforePlayDoesNothing

// START: TestClickMiddleCellSeeks

func TestClickMiddleCellSeeks(t *testing.T) {
	m, pl := barModel(248 * time.Second)
	clickCell(t, m, 19)
	if want := []time.Duration{127 * time.Second}; !reflect.DeepEqual(pl.seekTos, want) {
		t.Fatalf("SeekTo calls = %v, want %v (248 s × 19 / 37, rounded down)", pl.seekTos, want)
	}
}

// END: TestClickMiddleCellSeeks

// START: TestClickEndCellsSeek

func TestClickEndCellsSeek(t *testing.T) {
	for cell, want := range map[int]time.Duration{0: 0, barCells - 1: 248 * time.Second} {
		m, pl := barModel(248 * time.Second)
		clickCell(t, m, cell)
		if !reflect.DeepEqual(pl.seekTos, []time.Duration{want}) {
			t.Errorf("cell %d: SeekTo calls = %v, want [%v]", cell, pl.seekTos, want)
		}
	}
}

// END: TestClickEndCellsSeek

// START: TestPressLeftOfBarSeeksNothing

func TestPressLeftOfBarSeeksNothing(t *testing.T) {
	m, pl := barModel(248 * time.Second)
	x, y := barCell(t, m)
	release(press(m, x-1, y), x-1, y)
	if len(pl.seekTos) != 0 {
		t.Fatalf("SeekTo calls = %v after a click one cell left of the bar, want none", pl.seekTos)
	}
}

// END: TestPressLeftOfBarSeeksNothing

// START: TestDragShowsPositionWithoutSeeking

func TestDragShowsPositionWithoutSeeking(t *testing.T) {
	m, pl := barModel(248 * time.Second)
	x, y := barCell(t, m)
	m = move(press(m, x+10, y), x+30, y)
	if len(pl.seekTos) != 0 {
		t.Fatalf("SeekTo calls = %v during the drag, want none", pl.seekTos)
	}
	if got := plain(m.View().Content); !strings.Contains(got, "3:21 / 4:08") {
		t.Fatalf("view during the drag has no clock 3:21 / 4:08 (248 s × 30 / 37 = 201 s):\n%s", got)
	}
}

// END: TestDragShowsPositionWithoutSeeking

// START: TestReleaseSeeksOnce

func TestReleaseSeeksOnce(t *testing.T) {
	m, pl := barModel(248 * time.Second)
	x, y := barCell(t, m)
	release(move(press(m, x+10, y), x+30, y), x+30, y)
	if want := []time.Duration{201 * time.Second}; !reflect.DeepEqual(pl.seekTos, want) {
		t.Fatalf("SeekTo calls = %v, want %v", pl.seekTos, want)
	}
}

// END: TestReleaseSeeksOnce

// START: TestReleaseRightOfBarSeeksToEnd

func TestReleaseRightOfBarSeeksToEnd(t *testing.T) {
	m, pl := barModel(248 * time.Second)
	x, y := barCell(t, m)
	release(move(press(m, x+10, y), x+barCells+5, y), x+barCells+5, y)
	if want := []time.Duration{248 * time.Second}; !reflect.DeepEqual(pl.seekTos, want) {
		t.Fatalf("SeekTo calls = %v, want %v (the nearest end)", pl.seekTos, want)
	}
}

// END: TestReleaseRightOfBarSeeksToEnd

// START: TestReleaseLeftOfBarSeeksToStart

func TestReleaseLeftOfBarSeeksToStart(t *testing.T) {
	m, pl := barModel(248 * time.Second)
	x, y := barCell(t, m)
	release(move(press(m, x+10, y), x-3, y), x-3, y)
	if want := []time.Duration{0}; !reflect.DeepEqual(pl.seekTos, want) {
		t.Fatalf("SeekTo calls = %v, want %v (the nearest end)", pl.seekTos, want)
	}
}

// END: TestReleaseLeftOfBarSeeksToStart

// START: TestMouseOffSeeksNothing

func TestMouseOffSeeksNothing(t *testing.T) {
	m, pl := barModel(248 * time.Second)
	clickCell(t, m.WithMouse(false), 19)
	if len(pl.seekTos) != 0 {
		t.Fatalf("SeekTo calls = %v with mouse = false, want none", pl.seekTos)
	}
}

// END: TestMouseOffSeeksNothing

// START: TestNoLengthSeeksNothing

func TestNoLengthSeeksNothing(t *testing.T) {
	m, pl := barModel(0)
	clickCell(t, m, 19)
	if len(pl.seekTos) != 0 {
		t.Fatalf("SeekTo calls = %v with a total of 0, want none", pl.seekTos)
	}
}

// END: TestNoLengthSeeksNothing

// START: TestReleaseWithoutPressSeeksNothing

func TestReleaseWithoutPressSeeksNothing(t *testing.T) {
	m, pl := barModel(248 * time.Second)
	x, y := barCell(t, m)
	release(m, x+19, y)
	if len(pl.seekTos) != 0 {
		t.Fatalf("SeekTo calls = %v after a release with no press, want none", pl.seekTos)
	}
}

// END: TestReleaseWithoutPressSeeksNothing
