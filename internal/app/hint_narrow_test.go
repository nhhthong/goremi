// Tests that the hint line is cut to the terminal width and never wraps.
package app

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

// wholeHint is the hint the first view line shows when it fits.
const wholeHint = "Ctrl+C: quit · run goremi theme to choose a theme"

// START: TestHintIsCutToNarrowWidth

func TestHintIsCutToNarrowWidth(t *testing.T) {
	lines := strings.Split(plain(sized(New(fakeProvider{}), 30).View().Content), "\n")
	if w := lipgloss.Width(lines[0]); w > 30 || !strings.HasSuffix(lines[0], "…") {
		t.Fatalf("first line = %q (%d columns), want at most 30 and ending with …", lines[0], w)
	}
	if !strings.HasPrefix(lines[1], "Search:") {
		t.Fatalf("second line = %q, want it to start with Search:", lines[1])
	}
}

// END: TestHintIsCutToNarrowWidth

// START: TestHintShowsWholeAtItsWidth

func TestHintShowsWholeAtItsWidth(t *testing.T) {
	if got := strings.Split(plain(sized(New(fakeProvider{}), lipgloss.Width(wholeHint)).View().Content), "\n")[0]; got != wholeHint {
		t.Fatalf("first line = %q, want %q", got, wholeHint)
	}
}

// END: TestHintShowsWholeAtItsWidth

// START: TestHintShowsWholeInWideTerminal

func TestHintShowsWholeInWideTerminal(t *testing.T) {
	if got := strings.Split(plain(sized(New(fakeProvider{}), 80).View().Content), "\n")[0]; got != wholeHint {
		t.Fatalf("first line = %q, want %q", got, wholeHint)
	}
}

// END: TestHintShowsWholeInWideTerminal

// START: TestHintShowsWholeWithUnknownWidth

func TestHintShowsWholeWithUnknownWidth(t *testing.T) {
	if got := strings.Split(plain(New(fakeProvider{}).View().Content), "\n")[0]; got != wholeHint {
		t.Fatalf("first line = %q, want %q", got, wholeHint)
	}
}

// END: TestHintShowsWholeWithUnknownWidth
