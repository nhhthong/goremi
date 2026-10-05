// Tests for the mapping of a band level in dB onto a bar height of 0 to 64 levels.
package ui

import (
	"math"
	"testing"
)

// START: TestSpectrumHeightFloorIsZero

func TestSpectrumHeightFloorIsZero(t *testing.T) {
	if got := SpectrumHeight(-60); got != 0 {
		t.Fatalf("SpectrumHeight(-60) = %d, want 0", got)
	}
}

// END: TestSpectrumHeightFloorIsZero

// START: TestSpectrumHeightCeilingIs64

func TestSpectrumHeightCeilingIs64(t *testing.T) {
	if got := SpectrumHeight(0); got != 64 {
		t.Fatalf("SpectrumHeight(0) = %d, want 64", got)
	}
}

// END: TestSpectrumHeightCeilingIs64

// START: TestSpectrumHeightClampsOutsideRange

func TestSpectrumHeightClampsOutsideRange(t *testing.T) {
	for db, want := range map[float64]int{-80: 0, math.Inf(-1): 0, 6: 64} {
		if got := SpectrumHeight(db); got != want {
			t.Fatalf("SpectrumHeight(%v) = %d, want %d", db, got, want)
		}
	}
}

// END: TestSpectrumHeightClampsOutsideRange

// START: TestSpectrumHeightNeverFallsAsLevelRises

func TestSpectrumHeightNeverFallsAsLevelRises(t *testing.T) {
	prev := SpectrumHeight(-60)
	for db := -59.0; db <= 0; db++ {
		h := SpectrumHeight(db)
		if h < prev {
			t.Fatalf("SpectrumHeight(%v) = %d after %d: it fell", db, h, prev)
		}
		prev = h
	}
	if mid := SpectrumHeight(-30); mid <= 0 || mid >= 64 {
		t.Fatalf("SpectrumHeight(-30) = %d, want strictly between 0 and 64", mid)
	}
}

// END: TestSpectrumHeightNeverFallsAsLevelRises
