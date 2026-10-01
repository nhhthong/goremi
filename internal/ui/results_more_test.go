// Tests for the "load more..." line below the track list.
package ui

import (
	"strings"
	"testing"

	"goremi/internal/provider"
)

// START: threeTracks

func threeTracks() []provider.Track {
	return []provider.Track{{Title: "Get Lucky"}, {Title: "One More Time"}, {Title: "Instant Crush"}}
}

// END: threeTracks

// START: TestMoreLineAfterTracks

func TestMoreLineAfterTracks(t *testing.T) {
	tracks := threeTracks()
	out := RenderResultsMore(tracks, 0)
	got := lines(out)
	if len(got) != 4 {
		t.Fatalf("want 4 lines, got %q", got)
	}
	if !strings.HasPrefix(out, RenderResults(tracks, 0)+"\n") {
		t.Errorf("track lines changed: %q", out)
	}
	if !strings.Contains(got[3], "load more...") {
		t.Errorf("last line = %q, want load more...", got[3])
	}
}

// END: TestMoreLineAfterTracks

// START: TestMoreLineSelected

func TestMoreLineSelected(t *testing.T) {
	out := RenderResultsMore(threeTracks(), 3)
	got := lines(out)
	if !strings.HasPrefix(got[3], "▶ ") {
		t.Errorf("last line = %q, want prefix ▶ ", got[3])
	}
	if n := strings.Count(out, "▶"); n != 1 {
		t.Errorf("want exactly one ▶, got %d in %q", n, out)
	}
}

// END: TestMoreLineSelected

// START: TestMoreLineUnselected

func TestMoreLineUnselected(t *testing.T) {
	got := lines(RenderResultsMore(threeTracks(), 1))
	if !strings.HasPrefix(got[3], "  ") || !strings.Contains(got[3], "load more...") {
		t.Errorf("last line = %q, want two-space prefix and load more...", got[3])
	}
}

// END: TestMoreLineUnselected
