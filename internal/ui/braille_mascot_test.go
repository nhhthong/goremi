// Tests of the braille mascot and its header (ui tasks 3.10.1 to 3.10.3).
package ui

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"charm.land/lipgloss/v2"

	"goremi/internal/ui/theme"
)

// START: TestMascotIsTwelveByThirtyTwoBraille

func TestMascotIsTwelveByThirtyTwoBraille(t *testing.T) {
	rows := strings.Split(Mascot(), "\n")
	if len(rows) != 12 {
		t.Fatalf("%d rows, want 12", len(rows))
	}
	for i, row := range rows {
		if n := utf8.RuneCountInString(row); n != 32 {
			t.Errorf("row %d is %d columns, want 32", i, n)
		}
		for _, r := range row {
			if r < 0x2800 || r > 0x28FF {
				t.Fatalf("row %d holds %q, want only braille characters", i, r)
			}
		}
	}
}

// END: TestMascotIsTwelveByThirtyTwoBraille

// START: TestMascotHasNoColourCode

func TestMascotHasNoColourCode(t *testing.T) {
	if strings.Contains(Mascot(), "\x1b") {
		t.Fatal("Mascot() carries a colour code, want plain text")
	}
}

// END: TestMascotHasNoColourCode

// START: TestPaintMascotRunsFromLogoFromToLogoTo

// dotted is the 24-bit code of the gradient colour of column col.
func dotted(th theme.Theme, col int) string {
	r, g, b, _ := lipgloss.Blend1D(32, lipgloss.Color(th.LogoFrom), lipgloss.Color(th.LogoTo))[col].RGBA()
	return fmt.Sprintf("38;2;%d;%d;%d", r>>8, g>>8, b>>8)
}

func TestPaintMascotRunsFromLogoFromToLogoTo(t *testing.T) {
	th := theme.Default()
	got := PaintMascot(th, nil)
	if ansiCode.ReplaceAllString(got, "") != Mascot() {
		t.Fatal("painting changed the text")
	}
	rows := strings.Split(got, "\n")
	// row 6 holds a dot in column 27 (the right end of the guitar), row 11 a dot in column 0 (the ground line)
	if !strings.Contains(rows[6], dotted(th, 27)) || !strings.Contains(rows[11], dotted(th, 0)) {
		t.Fatalf("rows lack the colour of column 27 (%s) or column 0 (%s)", dotted(th, 27), dotted(th, 0))
	}
}

// END: TestPaintMascotRunsFromLogoFromToLogoTo

// START: TestPaintMascotFollowsTheme

func TestPaintMascotFollowsTheme(t *testing.T) {
	a, b := PaintMascot(theme.Dracula(), nil), PaintMascot(theme.Nord(), nil)
	if a == b || !strings.Contains(a, dotted(theme.Dracula(), 27)) || !strings.Contains(b, dotted(theme.Nord(), 27)) {
		t.Fatal("the painted mascot does not follow the theme")
	}
}

// END: TestPaintMascotFollowsTheme

// START: TestHeaderWideBadgeBesideBrailleMascot

func TestHeaderWideBadgeBesideBrailleMascot(t *testing.T) {
	rows := plainRows(Header(theme.Default(), 100, "0.1.1", nil))
	if len(rows) != 12 {
		t.Fatalf("%d rows, want the 12 rows of the mascot", len(rows))
	}
	for i, want := range Badge("0.1.1") {
		if got := col(rows[4+i], want); got != 34 {
			t.Errorf("row %d: %q starts at column %d, want 34 (32 of mascot and 2 apart)", 4+i, want, got)
		}
	}
	for _, i := range []int{0, 1, 2, 3, 7, 8, 9, 10, 11} {
		if strings.Contains(rows[i], "Goremi") || strings.Contains(rows[i], "Created") || strings.Contains(rows[i], "Current") {
			t.Errorf("row %d holds badge text, want the badge on rows 4 to 6", i)
		}
	}
}

// END: TestHeaderWideBadgeBesideBrailleMascot

// START: TestHeaderNarrowBadgeUnderBrailleMascot

func TestHeaderNarrowBadgeUnderBrailleMascot(t *testing.T) {
	rows := plainRows(Header(theme.Default(), 60, "0.1.1", nil))
	if len(rows) != 15 || !reflect.DeepEqual(rows[12:], Badge("0.1.1")) {
		t.Fatalf("%d rows, last three %q; want 12 of mascot then the 3 badge lines at column 0", len(rows), rows[max(len(rows)-3, 0):])
	}
}

// END: TestHeaderNarrowBadgeUnderBrailleMascot

// START: TestHeaderRows

func TestHeaderRows(t *testing.T) {
	if HeaderRows(100) != 12 || HeaderRows(60) != 15 {
		t.Fatalf("HeaderRows(100) = %d, HeaderRows(60) = %d; want 12 and 15", HeaderRows(100), HeaderRows(60))
	}
}

// END: TestHeaderRows
