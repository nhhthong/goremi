// Tests for painting the spectrum rows with the Spectrum colour of the theme.
package ui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"goremi/internal/ui/theme"
)

// fgCode is the 24-bit foreground code of a #rrggbb colour.
func fgCode(hex string) string {
	n, _ := strconv.ParseUint(strings.TrimPrefix(hex, "#"), 16, 32)
	return fmt.Sprintf("38;2;%d;%d;%d", n>>16&0xff, n>>8&0xff, n&0xff)
}

// sampleRows are three rows of bars to paint.
func sampleRows() []string {
	return SpectrumRows(heightsOf(20))
}

// START: TestPaintSpectrumUsesThemeColour

func TestPaintSpectrumUsesThemeColour(t *testing.T) {
	for i, r := range PaintSpectrum(theme.Default(), sampleRows()) {
		plain := regexp.MustCompile("\x1b\\[[0-9;]*m").ReplaceAllString(r, "")
		if strings.TrimSpace(plain) == "" {
			continue
		}
		if !strings.Contains(r, "38;2;45;212;191") {
			t.Fatalf("row %d = %q, want the code 38;2;45;212;191", i, r)
		}
	}
}

// END: TestPaintSpectrumUsesThemeColour

// START: TestPaintSpectrumKeepsText

func TestPaintSpectrumKeepsText(t *testing.T) {
	rows := sampleRows()
	for i, r := range PaintSpectrum(theme.Default(), rows) {
		if got := regexp.MustCompile("\x1b\\[[0-9;]*m").ReplaceAllString(r, ""); got != rows[i] {
			t.Fatalf("row %d text = %q, want %q", i, got, rows[i])
		}
	}
}

// END: TestPaintSpectrumKeepsText

// START: TestPaintSpectrumFollowsEveryTheme

func TestPaintSpectrumFollowsEveryTheme(t *testing.T) {
	for _, e := range theme.All() {
		painted := strings.Join(PaintSpectrum(e.Theme, sampleRows()), "\n")
		if !strings.Contains(painted, fgCode(e.Theme.Spectrum)) {
			t.Fatalf("theme %s: the output has no %s", e.Name, fgCode(e.Theme.Spectrum))
		}
	}
}

// END: TestPaintSpectrumFollowsEveryTheme
