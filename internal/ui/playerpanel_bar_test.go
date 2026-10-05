// Tests for the progress row of the player panel.
package ui

import (
	"strings"
	"testing"
	"time"

	"goremi/internal/ui/theme"
)

// checkBar fails when the plain BarLine is not filled ━, then ●, then rest ─.
func checkBar(t *testing.T, width, elapsed, total, filled, rest int) {
	t.Helper()
	got := ansiCode.ReplaceAllString(BarLine(theme.Default(), width, time.Duration(elapsed)*time.Second, time.Duration(total)*time.Second), "")
	if want := strings.Repeat("━", filled) + "●" + strings.Repeat("─", rest); got != want {
		t.Fatalf("BarLine(%d, %ds of %ds) = %q, want %q", width, elapsed, total, got, want)
	}
}

// START: TestBarMiddle

func TestBarMiddle(t *testing.T) { checkBar(t, 40, 124, 248, 19, 18) }

// END: TestBarMiddle

// START: TestBarStart

func TestBarStart(t *testing.T) { checkBar(t, 40, 0, 248, 0, 37) }

// END: TestBarStart

// START: TestBarEnd

func TestBarEnd(t *testing.T) { checkBar(t, 40, 248, 248, 37, 0) }

// END: TestBarEnd

// START: TestBarFollowsWidth

func TestBarFollowsWidth(t *testing.T) { checkBar(t, 20, 5, 20, 4, 13) }

// END: TestBarFollowsWidth

// START: TestBarColours

func TestBarColours(t *testing.T) {
	th := theme.Theme{Progress: "#ff0000", Border: "#0000ff"}
	bar := BarLine(th, 40, 124*time.Second, 248*time.Second)
	progress, border := strings.Index(bar, "38;2;255;0;0"), strings.Index(bar, "38;2;0;0;255")
	if progress < 0 || border < 0 || progress > strings.Index(bar, "━") || progress > strings.Index(bar, "●") || border > strings.Index(bar, "─") {
		t.Fatalf("Progress code at %d, Border code at %d in %q, want Progress before ━ and ●, Border before ─", progress, border, bar)
	}
}

// END: TestBarColours
