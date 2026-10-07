// Tests of the ASCII banner in the view: above the mascot, at every state, in the theme colour, from 50 columns (ui tasks 3.9.3 to 3.9.5).
package app

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"goremi/internal/ui/theme"
)

// bannerRows are the rows of ascii.txt.
func bannerRows(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile("../../ascii.txt")
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimRight(string(data), "\n"), "\n")
}

// START: TestBannerAboveMascotAndBadge

func TestBannerAboveMascotAndBadge(t *testing.T) {
	lines := plainLines(sizeTo(New(fakeProvider{}), 100, 30))
	for i, want := range bannerRows(t) {
		if got := strings.TrimRight(lines[i], " "); got != strings.TrimRight(want, " ") {
			t.Fatalf("line %d = %q, want banner row %q", i, got, want)
		}
	}
	if b := lineWith(lines, "Goremi v"); b < 6 {
		t.Fatalf("the badge is on line %d, want it under the banner (line 6 or later)", b)
	}
}

// END: TestBannerAboveMascotAndBadge

// START: TestBannerWhilePlaying

func TestBannerWhilePlaying(t *testing.T) {
	got := plain(playedList(100, 40).View().Content)
	if !strings.Contains(got, "██████╗") || !strings.Contains(got, "Created by nhhthong") || !strings.Contains(got, "Loading…") {
		t.Fatalf("view while a track plays:\n%s\nwant the banner, the badge and the panel", got)
	}
}

// END: TestBannerWhilePlaying

// START: TestBannerUsesModelTheme

func TestBannerUsesModelTheme(t *testing.T) {
	for _, th := range []theme.Theme{theme.Default(), theme.Nord()} {
		first := viewLines(sizeTo(New(fakeProvider{}).WithTheme(th), 100, 30))[0]
		if !strings.Contains(first, bannerCode(th.LogoTo)) {
			t.Errorf("banner row %q lacks LogoTo %s of the theme", first, bannerCode(th.LogoTo))
		}
	}
}

// END: TestBannerUsesModelTheme

// START: TestNoBannerBelowFiftyColumns

func TestNoBannerBelowFiftyColumns(t *testing.T) {
	if got := plain(sizeTo(New(fakeProvider{}), 40, 30).View().Content); strings.Contains(got, "██") {
		t.Fatalf("a banner shows at 40 columns:\n%s", got)
	}
}

// END: TestNoBannerBelowFiftyColumns

// START: TestBannerTakesRowsFromTheList

func TestBannerTakesRowsFromTheList(t *testing.T) {
	got := plain(pageAt(100, 25, false).View().Content)
	if !strings.Contains(got, "  T2") || strings.Contains(got, "T3") {
		t.Fatalf("view:\n%s\nwant 3 list rows (25 − 3 bar − 1 hint − 12 mascot − 6 banner)", got)
	}
}

// END: TestBannerTakesRowsFromTheList

// bannerCode is the 24-bit foreground code of a #rrggbb colour.
func bannerCode(hex string) string {
	var r, g, b int
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return fmt.Sprintf("38;2;%d;%d;%d", r, g, b)
}
