// Tests of the braille mascot in the view: the theme colour and the rows it takes (ui task 3.10.4).
package app

import (
	"fmt"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"

	"goremi/internal/ui/theme"
)

// START: TestMascotColourFollowsTheModelTheme

func TestMascotColourFollowsTheModelTheme(t *testing.T) {
	for _, th := range []theme.Theme{theme.Default(), theme.Nord()} {
		row := viewLines(sizeTo(New(fakeProvider{}).WithTheme(th), 100, 40))[6+6] // 6 banner rows, then row 6 of the mascot, which ends in column 27
		r, g, b, _ := lipgloss.Blend1D(32, lipgloss.Color(th.LogoFrom), lipgloss.Color(th.LogoTo))[27].RGBA()
		if want := fmt.Sprintf("38;2;%d;%d;%d", r>>8, g>>8, b>>8); !strings.Contains(row, want) {
			t.Errorf("mascot row %q lacks the gradient colour %s of column 27", row, want)
		}
	}
}

// END: TestMascotColourFollowsTheModelTheme

// START: TestMascotRowsCountInTheListWindow

func TestMascotRowsCountInTheListWindow(t *testing.T) {
	got := plain(pageAt(100, 25, false).View().Content)
	if !strings.Contains(got, "  T2") || strings.Contains(got, "T3") {
		t.Fatalf("view:\n%s\nwant 3 list rows (25 − 3 bar − 1 hint − 12 mascot − 6 banner)", got)
	}
}

// END: TestMascotRowsCountInTheListWindow
