// Tests for cutting a long track title so the list line fits beside the panel.
package app

import (
	"strings"
	"testing"

	"goremi/internal/provider"
)

// titled returns a model of the given width whose list holds one track with a title of n letters t.
func titled(width, n int) Model {
	p := &scriptedProvider{steps: []step{{tracks: []provider.Track{{Title: strings.Repeat("t", n)}}}}}
	m := search(sized(New(p), width), "daft")
	next, _ := m.Update(PlayMsg{Track: m.Tracks()[0]}) // after the first play the panel sits beside the list
	return next.(Model)
}

// listLine is the plain text of the first view line that holds the selected track.
func listLine(t *testing.T, m Model) []rune {
	t.Helper()
	for _, l := range viewLines(m) {
		if p := plain(l); strings.Contains(p, "▶ t") {
			return []rune(p)
		}
	}
	t.Fatalf("no list line in %q", viewLines(m))
	return nil
}

// START: TestSideBySideLongTitleCut

func TestSideBySideLongTitleCut(t *testing.T) {
	line := listLine(t, titled(100, 70))
	if got := string(line[:58]); !strings.HasSuffix(got, "…") || line[58] != ' ' || line[59] != ' ' {
		t.Fatalf("list part %q, then %q; want 58 columns ending with … and the two-space gap", got, string(line[58:60]))
	}
}

// END: TestSideBySideLongTitleCut

// START: TestSideBySideTitleAtLimit

func TestSideBySideTitleAtLimit(t *testing.T) {
	line := listLine(t, titled(100, 56))
	if got := string(line[:58]); strings.Contains(got, "…") || !strings.HasSuffix(got, "t") {
		t.Fatalf("list part %q, want the whole title (2 + 56 = 58 columns)", got)
	}
}

// END: TestSideBySideTitleAtLimit

// START: TestSideBySideTitleOverLimit

func TestSideBySideTitleOverLimit(t *testing.T) {
	line := listLine(t, titled(100, 57))
	if got := string(line[:58]); !strings.HasSuffix(got, "…") {
		t.Fatalf("list part %q, want the title cut with … (2 + 57 = 59 columns)", got)
	}
}

// END: TestSideBySideTitleOverLimit

// START: TestStackedLongTitleCut

func TestStackedLongTitleCut(t *testing.T) {
	line := listLine(t, titled(60, 70))
	if len(line) != 60 || line[59] != '…' {
		t.Fatalf("list line %q (%d columns), want 60 columns ending with …", string(line), len(line))
	}
}

// END: TestStackedLongTitleCut

// START: TestStackedTitleAtLimit

func TestStackedTitleAtLimit(t *testing.T) {
	line := listLine(t, titled(60, 58))
	if len(line) != 60 || line[59] != 't' {
		t.Fatalf("list line %q (%d columns), want the whole title (2 + 58 = 60 columns)", string(line), len(line))
	}
}

// END: TestStackedTitleAtLimit
