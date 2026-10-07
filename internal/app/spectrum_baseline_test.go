// Test that the baseline row sits under the spectrum, above the artist line, in the Border colour (task 6.24).
package app

import (
	"strings"
	"testing"
)

// START: TestBaselineUnderSpectrum

func TestBaselineUnderSpectrum(t *testing.T) {
	m, _ := playAndSettle(specApp(&recordingPlayer{}, true), specTrack)
	lines := panelLines(m)
	if lines[8] != strings.Repeat("─", 40) || !strings.Contains(lines[9], "Daft Punk") {
		t.Fatalf("panel lines 8 and 9 = %q and %q, want the baseline then the artist", lines[8], lines[9])
	}
	raw := viewLines(m)
	if !strings.Contains(raw[panelHead(m)+8], "38;2;107;114;128") {
		t.Fatalf("baseline %q lacks the Border colour 107;114;128", raw[panelHead(m)+8])
	}
}

// END: TestBaselineUnderSpectrum
