// Tests of the notes in the view: they run with a playing track, stop with a pause, and are drawn in the mascot (ui task 3.11.5).
package app

import (
	"strings"
	"testing"

	"goremi/internal/provider"
)

// noting plays a track in a 100x40 model and returns it after the details arrived, a random source that always spawns being set.
func noting(t *testing.T) Model {
	t.Helper()
	m := playedList(100, 40)
	m.intn = func(int) int { return 0 }
	next, _ := m.Update(detailsMsg{track: provider.Track{ID: trackB.ID, Title: trackB.Title, Artist: "Daft Punk"}})
	return next.(Model)
}

// START: TestNoteTickRunsWhilePlaying

func TestNoteTickRunsWhilePlaying(t *testing.T) {
	m := noting(t)
	if !m.noting {
		t.Fatal("no note tick runs after the details of a play")
	}
	next, cmd := m.Update(noteTickMsg{})
	got := next.(Model)
	if len(got.notes) != 1 || cmd == nil {
		t.Fatalf("after a tick: %d notes, command %v; want one note and the next tick", len(got.notes), cmd)
	}
}

// END: TestNoteTickRunsWhilePlaying

// START: TestNotesAreDrawnInTheMascot

func TestNotesAreDrawnInTheMascot(t *testing.T) {
	next, _ := noting(t).Update(noteTickMsg{})
	if got := plain(next.(Model).View().Content); !strings.Contains(got, "⢸⢳") || !strings.Contains(got, "⠸⠟") {
		t.Fatalf("view holds no note glyph:\n%s", got)
	}
}

// END: TestNotesAreDrawnInTheMascot

// START: TestPauseStopsTheNotes

func TestPauseStopsTheNotes(t *testing.T) {
	m := noting(t)
	next, _ := m.Update(noteTickMsg{})
	m = next.(Model)
	m.paused = true
	next, cmd := m.Update(noteTickMsg{})
	got := next.(Model)
	if len(got.notes) != 0 || got.noting || cmd != nil {
		t.Fatalf("paused: %d notes, noting %v, command %v; want none, false, nil", len(got.notes), got.noting, cmd)
	}
}

// END: TestPauseStopsTheNotes

// START: TestNoNotesBeforeFirstPlay

func TestNoNotesBeforeFirstPlay(t *testing.T) {
	m := listAt(100, 40)
	m.intn = func(int) int { return 0 }
	next, cmd := m.Update(noteTickMsg{})
	if got := next.(Model); len(got.notes) != 0 || got.noting || cmd != nil {
		t.Fatalf("before a play: %d notes, noting %v, command %v; want none, false, nil", len(got.notes), got.noting, cmd)
	}
}

// END: TestNoNotesBeforeFirstPlay

// START: TestNotesComeBackWhenPlayResumes

func TestNotesComeBackWhenPlayResumes(t *testing.T) {
	m := noting(t)
	m.paused = true
	next, _ := m.Update(noteTickMsg{}) // the pause ends the tick loop
	m = next.(Model)
	next, cmd := m.Update(progressMsg{paused: false, total: 200e9}) // the next progress reading finds the play running again
	got := next.(Model)
	if !got.noting || cmd == nil {
		t.Fatalf("noting %v after the play resumed, want the tick loop started again", got.noting)
	}
}

// END: TestNotesComeBackWhenPlayResumes
