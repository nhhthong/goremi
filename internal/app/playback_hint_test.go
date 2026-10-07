// Tests for the keys in the search input and for the hint line under the interface.
package app

import (
	"strings"
	"testing"

	"goremi/internal/provider"
)

// inputWhilePlaying is a model that plays a track and has the search input focused.
func inputWhilePlaying(pl *recordingPlayer) Model {
	m, _ := playOne(New(&scriptedResolve{errs: []error{nil}}).WithPlayer(pl), provider.Track{ID: "a1", Title: "One"})
	return m.WithFocus(FocusInput)
}

// START: TestInputKeysAreText

func TestInputKeysAreText(t *testing.T) {
	pl := &recordingPlayer{}
	m, _ := typed(inputWhilePlaying(pl), "jklnp")
	if m.Query() != "jklnp" || pl.toggles != 0 || len(pl.seeks) != 0 {
		t.Fatalf("Query() = %q, toggles %d, seeks %v, want jklnp and no control", m.Query(), pl.toggles, pl.seeks)
	}
}

// END: TestInputKeysAreText

// START: TestInputSpaceIsText

func TestInputSpaceIsText(t *testing.T) {
	pl := &recordingPlayer{}
	m, _ := typed(inputWhilePlaying(pl), "a b")
	if m.Query() != "a b" || pl.toggles != 0 {
		t.Fatalf("Query() = %q, toggles %d, want \"a b\" and no pause", m.Query(), pl.toggles)
	}
}

// END: TestInputSpaceIsText

// lastLine is the last line of the plain view.
func lastLine(m Model) string {
	lines := strings.Split(plain(m.View().Content), "\n")
	return lines[len(lines)-1]
}

// listWhilePlaying is a model of the given width that plays a track with the list focused.
func listWhilePlaying(width int) Model {
	m, _ := playOne(sized(New(&scriptedResolve{errs: []error{nil}}).WithPlayer(&recordingPlayer{}), width), provider.Track{ID: "a1", Title: "One"})
	return m.WithFocus(FocusList)
}

const hintText = "j -10s  k pause  l +10s  p prev  n next"

// START: TestHintAbsentInInput

func TestHintAbsentInInput(t *testing.T) {
	if view := plain(inputWhilePlaying(&recordingPlayer{}).View().Content); strings.Contains(view, "j -10s") {
		t.Fatalf("hint line shows while the input has focus:\n%s", view)
	}
}

// END: TestHintAbsentInInput

// START: TestHintAbsentBeforePlay

func TestHintAbsentBeforePlay(t *testing.T) {
	m := sized(New(fakeProvider{}), 100).WithFocus(FocusList)
	if view := plain(m.View().Content); strings.Contains(view, "j -10s") {
		t.Fatalf("hint line shows before any track plays:\n%s", view)
	}
}

// END: TestHintAbsentBeforePlay
