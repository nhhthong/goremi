// Tests for the play/pause glyph following the pause state mpv reports.
package app

import (
	"strings"
	"testing"
)

// glyphAfterTick plays a track on a fake player whose pause state is mpvPaused, optionally presses k, runs one tick and returns the controls row.
func glyphAfterTick(t *testing.T, pressKey bool, mpvPaused bool) string {
	t.Helper()
	pl := &recordingPlayer{pausedNow: mpvPaused}
	m := sized(New(&scriptedResolve{errs: []error{nil}, artist: "Daft Punk"}).WithPlayer(pl), 100)
	m, cmd := playOne(m, panelTrack)
	next, _ := m.Update(cmd())
	m = next.(Model)
	if pressKey {
		m = pressK(m)
	}
	next, read := m.Update(tickMsg{})
	next, _ = next.Update(read())
	return controlsLine(next.(Model))
}

// START: TestGlyphFollowsMpvPause

func TestGlyphFollowsMpvPause(t *testing.T) {
	if row := glyphAfterTick(t, false, true); !strings.Contains(row, "▷") || strings.Contains(row, "❚❚") {
		t.Fatalf("controls row %q while mpv is paused, want ▷ and no ❚❚", row)
	}
}

// END: TestGlyphFollowsMpvPause

// START: TestGlyphCorrectsAppFlag

func TestGlyphCorrectsAppFlag(t *testing.T) {
	if row := glyphAfterTick(t, true, false); !strings.Contains(row, "❚❚") || strings.Contains(row, "▷") {
		t.Fatalf("controls row %q after k while mpv plays, want ❚❚ and no ▷", row)
	}
}

// END: TestGlyphCorrectsAppFlag

// START: TestGlyphStaysWhenAgreed

func TestGlyphStaysWhenAgreed(t *testing.T) {
	if row := glyphAfterTick(t, true, true); !strings.Contains(row, "▷") || strings.Contains(row, "❚❚") {
		t.Fatalf("controls row %q after k while mpv is paused, want ▷ and no ❚❚", row)
	}
}

// END: TestGlyphStaysWhenAgreed
