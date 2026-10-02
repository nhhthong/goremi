// Tests for PaintLogo: a gradient per column, art text unchanged.
package ui

import (
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"

	"goremi/internal/ui/theme"
)

var redToBlue = theme.Theme{LogoFrom: "#ff0000", LogoTo: "#0000ff"}

// colourCode matches the 24-bit foreground code of a styled cell.
var colourCode = regexp.MustCompile(`38;2;\d+;\d+;\d+`)

// cell is one painted character: its column and its colour code.
type cell struct {
	col  int
	code string
}

// paintedCells returns the painted cells of one row of PaintLogo, in column order.
func paintedCells(row string) []cell {
	var cells []cell
	col, code := 0, ""
	for i := 0; i < len(row); {
		if strings.HasPrefix(row[i:], "\x1b[") {
			end := strings.IndexByte(row[i:], 'm')
			code = colourCode.FindString(row[i+2 : i+end])
			i += end + 1
			continue
		}
		_, size := utf8.DecodeRuneInString(row[i:])
		if code != "" {
			cells = append(cells, cell{col, code})
		}
		col++
		i += size
	}
	return cells
}

// lastRowCells are the painted cells of the equalizer row, which spans the whole width of the art.
func lastRowCells() []cell {
	rows := strings.Split(PaintLogo(redToBlue), "\n")
	return paintedCells(rows[len(rows)-1])
}

// START: TestPaintLogoEnds

func TestPaintLogoEnds(t *testing.T) {
	cells := lastRowCells()
	if first := cells[0].code; first != "38;2;255;0;0" {
		t.Errorf("first cell colour = %q, want 38;2;255;0;0", first)
	}
	if last := cells[len(cells)-1].code; last != "38;2;0;0;255" {
		t.Errorf("last cell colour = %q, want 38;2;0;0;255", last)
	}
}

// END: TestPaintLogoEnds

// START: TestPaintLogoInterpolates

func TestPaintLogoInterpolates(t *testing.T) {
	seen := map[string]int{}
	for _, c := range lastRowCells() {
		if col, dup := seen[c.code]; dup {
			t.Errorf("columns %d and %d share the colour %q", col, c.col, c.code)
		}
		seen[c.code] = c.col
	}
}

// END: TestPaintLogoInterpolates

// START: TestPaintLogoColumnShared

func TestPaintLogoColumnShared(t *testing.T) {
	byColumn := map[int]string{}
	for _, row := range strings.Split(PaintLogo(redToBlue), "\n") {
		for _, c := range paintedCells(row) {
			if code, ok := byColumn[c.col]; ok && code != c.code {
				t.Errorf("column %d has the colours %q and %q", c.col, code, c.code)
			}
			byColumn[c.col] = c.code
		}
	}
}

// END: TestPaintLogoColumnShared

// START: TestPaintLogoKeepsArt

func TestPaintLogoKeepsArt(t *testing.T) {
	plain := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(PaintLogo(redToBlue), "")
	if plain != Logo {
		t.Fatalf("painted art without colour codes = %q, want %q", plain, Logo)
	}
}

// END: TestPaintLogoKeepsArt
