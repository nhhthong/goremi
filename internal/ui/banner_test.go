// Tests of the ASCII banner: its text, its width, its theme colours and when it shows (ui tasks 3.9.1, 3.9.2, 3.9.4).
package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	"charm.land/lipgloss/v2"

	"goremi/internal/ui/theme"
)

// code is the 24-bit foreground code of a #rrggbb colour.
func code(hex string) string {
	var rgb [3]int
	for i := range rgb {
		n, _ := strconv.ParseInt(hex[1+2*i:3+2*i], 16, 0)
		rgb[i] = int(n)
	}
	return fmt.Sprintf("38;2;%d;%d;%d", rgb[0], rgb[1], rgb[2])
}

// START: TestBannerMatchesFile

func TestBannerMatchesFile(t *testing.T) {
	data, err := os.ReadFile("../../ascii.txt")
	if err != nil {
		t.Fatal(err)
	}
	if want := strings.TrimRight(string(data), "\n"); Banner() != want {
		t.Fatalf("Banner() differs from ascii.txt:\n%s\nwant\n%s", Banner(), want)
	}
}

// END: TestBannerMatchesFile

// START: TestBannerIsFiftyWide

func TestBannerIsFiftyWide(t *testing.T) {
	rows := strings.Split(Banner(), "\n")
	if len(rows) != 6 {
		t.Fatalf("%d rows, want 6", len(rows))
	}
	for i, r := range rows {
		if n := utf8.RuneCountInString(r); n != 50 {
			t.Errorf("row %d is %d columns, want 50", i, n)
		}
	}
	if BannerWidth != 50 {
		t.Fatalf("BannerWidth = %d, want 50", BannerWidth)
	}
}

// END: TestBannerIsFiftyWide

// START: TestPaintBannerRunsFromLogoFromToLogoTo

// blend is the 24-bit foreground code of column col of the gradient the banner must follow.
func blend(th theme.Theme, col int) string {
	r, g, b, _ := lipgloss.Blend1D(BannerWidth, lipgloss.Color(th.LogoFrom), lipgloss.Color(th.LogoTo))[col].RGBA()
	return fmt.Sprintf("38;2;%d;%d;%d", r>>8, g>>8, b>>8)
}

func TestPaintBannerRunsFromLogoFromToLogoTo(t *testing.T) {
	th := theme.Default()
	got := PaintBanner(th)
	if ansiCode.ReplaceAllString(got, "") != Banner() {
		t.Fatal("painting changed the text")
	}
	row := strings.Split(got, "\n")[0] // the first block is in column 3, the last cell in column 49
	if !strings.Contains(row, blend(th, 3)) || !strings.HasSuffix(strings.TrimSuffix(row, "\x1b[m"), "╗") || !strings.Contains(row, code(th.LogoTo)) {
		t.Fatalf("row %q lacks the gradient colour of column 3 or LogoTo %s on the last cell", row, code(th.LogoTo))
	}
}

// END: TestPaintBannerRunsFromLogoFromToLogoTo

// START: TestPaintBannerFollowsTheme

func TestPaintBannerFollowsTheme(t *testing.T) {
	for _, th := range []theme.Theme{theme.Dracula(), theme.Nord()} {
		got := PaintBanner(th)
		if !strings.Contains(got, code(th.LogoTo)) {
			t.Errorf("banner lacks the end colour %s", code(th.LogoTo))
		}
	}
}

// END: TestPaintBannerFollowsTheme

// START: TestBannerRowsByWidth

func TestBannerRowsByWidth(t *testing.T) {
	for width, want := range map[int]int{0: 0, 49: 0, 50: 6, 100: 6} {
		if got := BannerRows(width); got != want {
			t.Errorf("BannerRows(%d) = %d, want %d", width, got, want)
		}
	}
}

// END: TestBannerRowsByWidth
