// Tests for PaintLogo: gradient colours on the art rows, art text unchanged.
package ui

import (
	"regexp"
	"strings"
	"testing"

	"goremi/internal/ui/theme"
)

var redToBlue = theme.Theme{LogoFrom: "#ff0000", LogoTo: "#0000ff"}

// colourCode matches the 24-bit foreground code of a styled row.
var colourCode = regexp.MustCompile(`38;2;\d+;\d+;\d+`)

// START: TestPaintLogoEnds

func TestPaintLogoEnds(t *testing.T) {
	rows := strings.Split(PaintLogo(redToBlue), "\n")
	if first := colourCode.FindString(rows[0]); first != "38;2;255;0;0" {
		t.Errorf("first row colour = %q, want 38;2;255;0;0", first)
	}
	if last := colourCode.FindString(rows[len(rows)-1]); last != "38;2;0;0;255" {
		t.Errorf("last row colour = %q, want 38;2;0;0;255", last)
	}
}

// END: TestPaintLogoEnds

// START: TestPaintLogoInterpolates

func TestPaintLogoInterpolates(t *testing.T) {
	seen := map[string]int{}
	for i, row := range strings.Split(PaintLogo(redToBlue), "\n") {
		code := colourCode.FindString(row)
		if j, dup := seen[code]; dup {
			t.Errorf("rows %d and %d share the colour %q", j, i, code)
		}
		seen[code] = i
	}
}

// END: TestPaintLogoInterpolates

// START: TestPaintLogoKeepsArt

func TestPaintLogoKeepsArt(t *testing.T) {
	plain := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(PaintLogo(redToBlue), "")
	if plain != Logo {
		t.Fatalf("painted art without colour codes = %q, want %q", plain, Logo)
	}
}

// END: TestPaintLogoKeepsArt
