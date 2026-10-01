// Tests for results list rendering: one line per track and the selection marker.
package ui

import (
	"strings"
	"testing"

	"goremi/internal/provider"
)

// START: lines

func lines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(strings.TrimRight(s, "\n"), "\n")
}

// END: lines

// START: TestResultsLines

func TestResultsLines(t *testing.T) {
	titles := []string{"Get Lucky", "One More Time", "Instant Crush"}
	var tracks []provider.Track
	for _, ti := range titles {
		tracks = append(tracks, provider.Track{Title: ti})
	}
	got := lines(RenderResults(tracks, 0))
	if len(got) != len(titles) {
		t.Fatalf("want %d lines, got %d: %q", len(titles), len(got), got)
	}
	for i, ti := range titles {
		if !strings.Contains(got[i], ti) {
			t.Errorf("line %d = %q, want it to contain %q", i, got[i], ti)
		}
	}
}

// END: TestResultsLines

// START: TestResultsEmpty

func TestResultsEmpty(t *testing.T) {
	if got := lines(RenderResults(nil, 0)); len(got) != 0 {
		t.Fatalf("want 0 lines, got %q", got)
	}
}

// END: TestResultsEmpty

// START: TestResultsTitleOnly

func TestResultsTitleOnly(t *testing.T) {
	out := RenderResults([]provider.Track{{Title: "Get Lucky", Artist: "Daft Punk"}}, 0)
	if !strings.Contains(out, "Get Lucky") {
		t.Errorf("output %q misses the Title", out)
	}
	if strings.Contains(out, "Daft Punk") {
		t.Errorf("output %q shows the Artist", out)
	}
}

// END: TestResultsTitleOnly

// START: checkMarker

func checkMarker(t *testing.T, selected int) {
	t.Helper()
	tracks := []provider.Track{{Title: "Get Lucky"}, {Title: "One More Time"}, {Title: "Instant Crush"}}
	out := RenderResults(tracks, selected)
	got := lines(out)
	for i, l := range got {
		want := "  "
		if i == selected {
			want = "▶ "
		}
		if !strings.HasPrefix(l, want) {
			t.Errorf("line %d = %q, want prefix %q", i, l, want)
		}
	}
	if n := strings.Count(out, "▶"); n != 1 {
		t.Errorf("want exactly one ▶, got %d in %q", n, out)
	}
}

// END: checkMarker

// START: TestSelectedMarkerMiddle

func TestSelectedMarkerMiddle(t *testing.T) { checkMarker(t, 1) }

// END: TestSelectedMarkerMiddle
// START: TestSelectedMarkerFirst

func TestSelectedMarkerFirst(t *testing.T) { checkMarker(t, 0) }

// END: TestSelectedMarkerFirst
// START: TestSelectedMarkerLast

func TestSelectedMarkerLast(t *testing.T) { checkMarker(t, 2) }

// END: TestSelectedMarkerLast
