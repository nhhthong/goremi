// Tests for the list and artwork pane layout.
package ui

import (
	"strings"
	"testing"
)

// START: checkLeftOf

func checkLeftOf(t *testing.T, line, left, right string) {
	t.Helper()
	l, r := strings.Index(line, left), strings.Index(line, right)
	if l < 0 || r < 0 || l >= r {
		t.Errorf("line %q: want %q before %q", line, left, right)
	}
}

// END: checkLeftOf

// START: TestJoinPanesEqual

func TestJoinPanesEqual(t *testing.T) {
	list := []string{"Get Lucky", "One More Time", "Instant Crush"}
	pane := []string{"ART1", "ART2", "ART3"}
	got := lines(JoinPanes(strings.Join(list, "\n"), strings.Join(pane, "\n")))
	if len(got) != 3 {
		t.Fatalf("want 3 lines, got %q", got)
	}
	for i := range list {
		checkLeftOf(t, got[i], list[i], pane[i])
	}
}

// END: TestJoinPanesEqual

// START: TestJoinPanesPaneTaller

func TestJoinPanesPaneTaller(t *testing.T) {
	pane := []string{"A", "B", "C"}
	got := lines(JoinPanes("Get Lucky", strings.Join(pane, "\n")))
	if len(got) != 3 {
		t.Fatalf("want 3 lines, got %q", got)
	}
	for i, p := range pane {
		if !strings.Contains(got[i], p) {
			t.Errorf("line %d = %q, want it to contain %q", i, got[i], p)
		}
	}
	checkLeftOf(t, got[0], "Get Lucky", "A")
}

// END: TestJoinPanesPaneTaller

// START: TestJoinPanesListTaller

func TestJoinPanesListTaller(t *testing.T) {
	list := []string{"Get Lucky", "One More Time", "Instant Crush"}
	got := lines(JoinPanes(strings.Join(list, "\n"), "ART"))
	if len(got) != 3 {
		t.Fatalf("want 3 lines, got %q", got)
	}
	for i, l := range list {
		if !strings.Contains(got[i], l) {
			t.Errorf("line %d = %q, want it to contain %q", i, got[i], l)
		}
	}
	checkLeftOf(t, got[0], "Get Lucky", "ART")
}

// END: TestJoinPanesListTaller
