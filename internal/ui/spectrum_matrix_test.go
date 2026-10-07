// Tests of the stacked block matrix of the spectrum and of its idle bars (tasks 6.18 to 6.21).
package ui

import (
	"math/rand"
	"testing"
)

// START: TestSpectrumCellsLow

func TestSpectrumCellsLow(t *testing.T) {
	for level, want := range map[int]int{0: 0, 1: 0, 2: 0, 3: 0, 4: 1} {
		if got := SpectrumCells(level); got != want {
			t.Errorf("SpectrumCells(%d) = %d, want %d", level, got, want)
		}
	}
}

// END: TestSpectrumCellsLow

// START: TestSpectrumCellsHigh

func TestSpectrumCellsHigh(t *testing.T) {
	for level, want := range map[int]int{11: 1, 12: 2, 60: 8, 64: 8} {
		if got := SpectrumCells(level); got != want {
			t.Errorf("SpectrumCells(%d) = %d, want %d", level, got, want)
		}
	}
}

// END: TestSpectrumCellsHigh

// START: TestIdleFlipLightsOneCell

func TestIdleFlipLightsOneCell(t *testing.T) {
	got := NextIdle([SpectrumBars]int{}, func(int) int { return 0 })
	for i, v := range got {
		if SpectrumCells(v) != 1 {
			t.Fatalf("bar %d = level %d (%d cells), want one cell", i, v, SpectrumCells(v))
		}
	}
}

// END: TestIdleFlipLightsOneCell

// START: TestIdleStaysZeroOrOneCell

func TestIdleStaysZeroOrOneCell(t *testing.T) {
	var h [SpectrumBars]int
	for frame := 0; frame < 1000; frame++ {
		h = NextIdle(h, rand.Intn)
		for i, v := range h {
			if c := SpectrumCells(v); c > 1 {
				t.Fatalf("frame %d bar %d has %d cells, want 0 or 1", frame, i, c)
			}
		}
	}
}

// END: TestIdleStaysZeroOrOneCell

// START: TestIdleAsksOneInTen

func TestIdleAsksOneInTen(t *testing.T) {
	asked := 0
	NextIdle([SpectrumBars]int{}, func(n int) int {
		if n != 10 {
			t.Fatalf("bound asked = %d, want 10 (a chance of 1 in 10)", n)
		}
		asked++
		return 1
	})
	if asked != SpectrumBars {
		t.Fatalf("asked %d times, want once per bar (%d)", asked, SpectrumBars)
	}
}

// END: TestIdleAsksOneInTen

// START: TestIdleKeepsStateUnlessDrawnZero

func TestIdleKeepsStateUnlessDrawnZero(t *testing.T) {
	var prev [SpectrumBars]int
	for i := range prev {
		prev[i] = (i % 2) * 8
	}
	if got := NextIdle(prev, func(int) int { return 5 }); got != prev {
		t.Fatalf("heights changed on a draw of 5: %v, want %v", got, prev)
	}
}

// END: TestIdleKeepsStateUnlessDrawnZero

// START: TestIdleFlipsOnDrawnZero

func TestIdleFlipsOnDrawnZero(t *testing.T) {
	var prev [SpectrumBars]int
	for i := range prev {
		prev[i] = (i % 2) * 8
	}
	got := NextIdle(prev, func(int) int { return 0 })
	for i := range got {
		if SpectrumCells(got[i]) == SpectrumCells(prev[i]) {
			t.Fatalf("bar %d stayed at %d cells on a draw of 0, want it flipped", i, SpectrumCells(got[i]))
		}
	}
}

// END: TestIdleFlipsOnDrawnZero
