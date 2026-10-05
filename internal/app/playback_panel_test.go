// Tests for the rows of the player panel while a track plays: controls, their order and the play/pause glyph.
package app

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/provider"
)

const controlsRow = "◀◀ ◀ ❚❚ ▶ ▶▶"

// panelTrack is the track the panel tests play: One by Daft Punk, 248 s.
var panelTrack = provider.Track{ID: "a1", Title: "One", Duration: 248 * time.Second}

// playedModel is a model of the given width that plays panelTrack, with its details loaded.
func playedModel(width int) Model {
	m := sized(New(&scriptedResolve{errs: []error{nil}, artist: "Daft Punk"}).WithPlayer(&recordingPlayer{}), width)
	m, cmd := playOne(m, panelTrack)
	next, _ := m.Update(cmd())
	return next.(Model)
}

// plainLines is the plain view split in lines.
func plainLines(m Model) []string { return strings.Split(plain(m.View().Content), "\n") }

// lineWith returns the index of the first line holding text, or -1.
func lineWith(lines []string, text string) int {
	for i, l := range lines {
		if strings.Contains(l, text) {
			return i
		}
	}
	return -1
}

// START: TestPanelShowsControls

func TestPanelShowsControls(t *testing.T) {
	if lineWith(plainLines(playedModel(100)), controlsRow) < 0 {
		t.Fatalf("no controls row in:\n%s", plain(playedModel(100).View().Content))
	}
}

// END: TestPanelShowsControls

// START: TestPanelShowsControlsNarrow

func TestPanelShowsControlsNarrow(t *testing.T) {
	if lineWith(plainLines(playedModel(70)), controlsRow) < 0 {
		t.Fatalf("no controls row at width 70 in:\n%s", plain(playedModel(70).View().Content))
	}
}

// END: TestPanelShowsControlsNarrow

// START: TestNoControlsBeforePlay

func TestNoControlsBeforePlay(t *testing.T) {
	if lineWith(plainLines(sized(New(fakeProvider{}), 100)), "◀◀") >= 0 {
		t.Fatal("controls row shows before any track plays")
	}
}

// END: TestNoControlsBeforePlay

// checkRowOrder fails unless the artist, title, bar, clock and controls lines come in that order.
func checkRowOrder(t *testing.T, m Model) {
	t.Helper()
	lines := plainLines(m)
	rows := []int{lineWith(lines, "Daft Punk"), lineWith(lines, "One"), lineWith(lines, "●"), lineWith(lines, "0:00 / 4:08"), lineWith(lines, "◀◀")}
	for i, r := range rows {
		if r < 0 || (i > 0 && r <= rows[i-1]) {
			t.Fatalf("rows artist, title, bar, clock, controls at lines %v in:\n%s", rows, plain(m.View().Content))
		}
	}
}

// START: TestPanelRowOrder

func TestPanelRowOrder(t *testing.T) { checkRowOrder(t, playedModel(100)) }

// END: TestPanelRowOrder

// START: TestPanelRowOrderNarrow

func TestPanelRowOrderNarrow(t *testing.T) { checkRowOrder(t, playedModel(70)) }

// END: TestPanelRowOrderNarrow

// START: TestPanelBarWidth

func TestPanelBarWidth(t *testing.T) {
	lines := plainLines(playedModel(100))
	i := lineWith(lines, "●")
	if i < 0 {
		t.Fatalf("no bar row in:\n%s", strings.Join(lines, "\n"))
	}
	cells := 0
	for _, r := range lines[i] {
		if strings.ContainsRune("━●─", r) {
			cells++
		}
	}
	if cells != 38 {
		t.Fatalf("bar row %q has %d cells, want 38", lines[i], cells)
	}
}

// END: TestPanelBarWidth

// pressK presses k with the list focused and returns the model.
func pressK(m Model) Model {
	next, _ := m.WithFocus(FocusList).Update(tea.KeyPressMsg{Code: 'k', Text: "k"})
	return next.(Model)
}

// controlsLine is the controls row of the view.
func controlsLine(m Model) string {
	lines := plainLines(m)
	return lines[lineWith(lines, "◀◀")]
}

// START: TestPlayGlyphWhilePlaying

func TestPlayGlyphWhilePlaying(t *testing.T) {
	if row := controlsLine(playedModel(100)); !strings.Contains(row, "❚❚") {
		t.Fatalf("controls row %q has no ❚❚ while playing", row)
	}
}

// END: TestPlayGlyphWhilePlaying

// START: TestPlayGlyphWhilePaused

func TestPlayGlyphWhilePaused(t *testing.T) {
	row := controlsLine(pressK(playedModel(100)))
	if !strings.Contains(row, "▷") || strings.Contains(row, "❚❚") || !strings.Contains(row, "▶") {
		t.Fatalf("controls row %q while paused, want ▷, no ❚❚ and still ▶", row)
	}
}

// END: TestPlayGlyphWhilePaused

// START: TestPlayGlyphAfterResume

func TestPlayGlyphAfterResume(t *testing.T) {
	if row := controlsLine(pressK(pressK(playedModel(100)))); !strings.Contains(row, "❚❚") {
		t.Fatalf("controls row %q after resuming, want ❚❚", row)
	}
}

// END: TestPlayGlyphAfterResume
