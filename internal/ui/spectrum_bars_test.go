// Tests of the 20 bars of half blocks and the baseline of the spectrum (tasks 6.22 to 6.24).
package ui

import (
	"strings"
	"testing"
)

// START: TestBarsTwentyAtEvenColumns

func TestBarsTwentyAtEvenColumns(t *testing.T) {
	for r, row := range SpectrumRows(heightsOf(64)) {
		cells := []rune(row)
		if len(cells) != 40 {
			t.Fatalf("row %d is %d columns, want 40", r, len(cells))
		}
		for c, ch := range cells {
			want := ' '
			if c%2 == 0 && c <= 38 {
				want = '▄'
			}
			if ch != want {
				t.Fatalf("row %d column %d = %q, want %q", r, c, ch, want)
			}
		}
	}
}

// END: TestBarsTwentyAtEvenColumns

// START: TestLitCellIsHalfBlock

func TestLitCellIsHalfBlock(t *testing.T) {
	var h [SpectrumBars]int
	h[0] = 12 // two cells
	rows := SpectrumRows(h)
	for i, r := range rows {
		want := strings.Repeat(" ", 40)
		if i >= 6 {
			want = "▄" + strings.Repeat(" ", 39)
		}
		if r != want {
			t.Fatalf("row %d = %q, want %q", i, r, want)
		}
	}
}

// END: TestLitCellIsHalfBlock

// START: TestBarTakesHighestBand

func TestBarTakesHighestBand(t *testing.T) {
	var h [SpectrumBars]int
	h[2] = 64 // bands 1 and 2 belong to bar 1, at column 2
	bottom := []rune(SpectrumRows(h)[7])
	if bottom[2] != '▄' || bottom[0] != ' ' || bottom[4] != ' ' {
		t.Fatalf("bottom row %q, want only column 2 lit", string(bottom))
	}
}

// END: TestBarTakesHighestBand

// START: TestLastBarTakesLastBands

func TestLastBarTakesLastBands(t *testing.T) {
	var h [SpectrumBars]int
	h[31] = 64 // bands 30 and 31 belong to bar 19, at column 38
	bottom := []rune(SpectrumRows(h)[7])
	if bottom[38] != '▄' || bottom[36] != ' ' {
		t.Fatalf("bottom row %q, want only column 38 lit", string(bottom))
	}
}

// END: TestLastBarTakesLastBands

// START: TestBaselineIsForty

func TestBaselineIsForty(t *testing.T) {
	if got := SpectrumBaseline(); got != strings.Repeat("─", 40) {
		t.Fatalf("SpectrumBaseline() = %q, want 40 ─", got)
	}
}

// END: TestBaselineIsForty
