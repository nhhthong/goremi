// Tests for smoothing the bar heights and drawing the bars as rows.
package ui

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// START: TestSmoothRisesAtOnce

func TestSmoothRisesAtOnce(t *testing.T) {
	if got := SmoothHeight(10, 40); got != 40 {
		t.Fatalf("SmoothHeight(10, 40) = %d, want 40", got)
	}
}

// END: TestSmoothRisesAtOnce

// START: TestSmoothFallsTo85Percent

func TestSmoothFallsTo85Percent(t *testing.T) {
	if got := SmoothHeight(60, 0); got != 51 {
		t.Fatalf("SmoothHeight(60, 0) = %d, want 51", got)
	}
	if got := SmoothHeight(64, 0); got != 54 {
		t.Fatalf("SmoothHeight(64, 0) = %d, want 54", got)
	}
}

// END: TestSmoothFallsTo85Percent

// START: TestSmoothKeepsTargetAboveFall

func TestSmoothKeepsTargetAboveFall(t *testing.T) {
	if got := SmoothHeight(60, 55); got != 55 {
		t.Fatalf("SmoothHeight(60, 55) = %d, want 55", got)
	}
}

// END: TestSmoothKeepsTargetAboveFall

// START: TestSmoothReachesZero

func TestSmoothReachesZero(t *testing.T) {
	h := 64
	for i := 0; i < 60; i++ {
		h = SmoothHeight(h, 0)
	}
	if h != 0 {
		t.Fatalf("height after 60 frames of silence = %d, want 0", h)
	}
}

// END: TestSmoothReachesZero

// heightsOf returns 32 heights, all h.
func heightsOf(h int) [SpectrumBars]int {
	var out [SpectrumBars]int
	for i := range out {
		out[i] = h
	}
	return out
}

// START: TestSpectrumRowsAre8By40

func TestSpectrumRowsAre8By40(t *testing.T) {
	rows := SpectrumRows(heightsOf(37))
	if len(rows) != 8 {
		t.Fatalf("%d rows, want 8", len(rows))
	}
	for i, r := range rows {
		if n := utf8.RuneCountInString(r); n != 40 {
			t.Fatalf("row %d is %d columns, want 40", i, n)
		}
	}
}

// END: TestSpectrumRowsAre8By40

// START: TestSpectrumRowsEmptyWhenZero

func TestSpectrumRowsEmptyWhenZero(t *testing.T) {
	for i, r := range SpectrumRows(heightsOf(0)) {
		if r != strings.Repeat(" ", 40) {
			t.Fatalf("row %d = %q, want 40 spaces", i, r)
		}
	}
}

// END: TestSpectrumRowsEmptyWhenZero
