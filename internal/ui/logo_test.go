// Tests for the embedded logo art.
package ui

import (
	"strings"
	"testing"
)

// START: TestLogoRows

func TestLogoRows(t *testing.T) {
	rows := strings.Split(Logo, "\n")
	if len(rows) != 6 || rows[3] != "" {
		t.Fatalf("Logo has %d rows and row 4 is %q, want 6 rows and an empty row 4", len(rows), rows[3])
	}
}

// END: TestLogoRows

// trimmedRows returns rows from to of the art with the spaces at both ends removed.
func trimmedRows(from, to int) []string {
	var out []string
	for _, r := range strings.Split(Logo, "\n")[from:to] {
		out = append(out, strings.TrimSpace(r))
	}
	return out
}

// START: TestLogoLetters

func TestLogoLetters(t *testing.T) {
	want := []string{"▄▀▀▀ ▄▀▀▄ █▀▀▄ █▀▀▀ █▄ ▄█ ▀█▀", "█ ▀█ █  █ █▀█  █▀▀  █ ▀ █  █", "▀▀▀  ▀▀  ▀  ▀ ▀▀▀▀ ▀   ▀ ▀▀▀"}
	for i, got := range trimmedRows(0, 3) {
		if got != want[i] {
			t.Errorf("row %d = %q, want %q", i+1, got, want[i])
		}
	}
}

// END: TestLogoLetters

// START: TestLogoEqualizer

func TestLogoEqualizer(t *testing.T) {
	want := []string{"▂   ▅ ▁ █ ▃   ▆   ▄   ▇ ▁   ▃", "▃ ▆ █ ▇ █ █ █ █ ▆ █ █ █ ▅ █ █ ▇ █ ▄ █ ▃"}
	for i, got := range trimmedRows(4, 6) {
		if got != want[i] {
			t.Errorf("row %d = %q, want %q", i+5, got, want[i])
		}
	}
}

// END: TestLogoEqualizer
