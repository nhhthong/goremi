// Tests for smoothing the bar heights, drawing the bars as rows, and the random heights while no audio plays.
package ui

import (
	"math/rand"
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

// START: TestSpectrumRowsFullAndCentred

func TestSpectrumRowsFullAndCentred(t *testing.T) {
	want := strings.Repeat(" ", 4) + strings.Repeat("█", 32) + strings.Repeat(" ", 4)
	for i, r := range SpectrumRows(heightsOf(64)) {
		if r != want {
			t.Fatalf("row %d = %q, want %q", i, r, want)
		}
	}
}

// END: TestSpectrumRowsFullAndCentred

// START: TestSpectrumRowsTopCell

func TestSpectrumRowsTopCell(t *testing.T) {
	var h [SpectrumBars]int
	h[0] = 13
	rows := SpectrumRows(h)
	cell := func(row int) string { return string([]rune(rows[row])[4]) }
	if cell(7) != "█" || cell(6) != "▅" {
		t.Fatalf("bottom cell %q and the one above %q, want █ and ▅", cell(7), cell(6))
	}
	for r := 0; r < 6; r++ {
		if cell(r) != " " {
			t.Fatalf("row %d cell = %q, want blank", r, cell(r))
		}
	}
}

// END: TestSpectrumRowsTopCell

// START: TestSpectrumRowsBoundaryHeights

func TestSpectrumRowsBoundaryHeights(t *testing.T) {
	var h [SpectrumBars]int
	h[0], h[1] = 1, 8
	rows := SpectrumRows(h)
	for col, want := range map[int]string{4: "▁", 5: "█"} {
		if got := string([]rune(rows[7])[col]); got != want {
			t.Fatalf("bottom cell of column %d = %q, want %q", col, got, want)
		}
		if got := string([]rune(rows[6])[col]); got != " " {
			t.Fatalf("cell above column %d = %q, want blank", col, got)
		}
	}
}

// END: TestSpectrumRowsBoundaryHeights

// START: TestIdleHeightsMax

func TestIdleHeightsMax(t *testing.T) {
	asked := 0
	got := IdleHeights(func(n int) int { asked = n; return n - 1 })
	if asked != 4 {
		t.Fatalf("bound asked = %d, want 4 (0 to 3)", asked)
	}
	if got != heightsOf(3) {
		t.Fatalf("heights = %v, want all 3", got)
	}
}

// END: TestIdleHeightsMax

// START: TestIdleHeightsMin

func TestIdleHeightsMin(t *testing.T) {
	if got := IdleHeights(func(int) int { return 0 }); got != heightsOf(0) {
		t.Fatalf("heights = %v, want all 0", got)
	}
}

// END: TestIdleHeightsMin

// START: TestIdleHeightsRandomWithinRange

func TestIdleHeightsRandomWithinRange(t *testing.T) {
	varied := false
	for i := 0; i < 1000; i++ {
		h := IdleHeights(rand.Intn)
		for _, v := range h {
			if v < 0 || v > 3 {
				t.Fatalf("height %d outside 0 to 3", v)
			}
		}
		if h != heightsOf(h[0]) {
			varied = true
		}
	}
	if !varied {
		t.Fatal("the 32 heights were the same in all 1000 frames")
	}
}

// END: TestIdleHeightsRandomWithinRange
