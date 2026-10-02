// The ASCII art "GOREMI" shown in the player panel (variant B of resources/ascii_terminal.png), and its painting.
package ui

import (
	"strings"

	"charm.land/lipgloss/v2"

	"goremi/internal/ui/theme"
)

// START: Logo

// Logo is the art: the capitals on three rows, a blank row, then a two-row equalizer; at most 40 columns wide.
var Logo = strings.Join([]string{
	"     ▄▀▀▀ ▄▀▀▄ █▀▀▄ █▀▀▀ █▄ ▄█ ▀█▀",
	"     █ ▀█ █  █ █▀█  █▀▀  █ ▀ █  █",
	"      ▀▀▀  ▀▀  ▀  ▀ ▀▀▀▀ ▀   ▀ ▀▀▀",
	"",
	"    ▂   ▅ ▁ █ ▃   ▆   ▄   ▇ ▁   ▃",
	"▃ ▆ █ ▇ █ █ █ █ ▆ █ █ █ ▅ █ █ ▇ █ ▄ █ ▃",
}, "\n")

// END: Logo

// START: PaintLogo

// PaintLogo colours the art column by column, from t.LogoFrom on the first column to t.LogoTo on the last;
// the cells of one column share a colour and spaces stay plain.
func PaintLogo(t theme.Theme) string {
	rows := strings.Split(Logo, "\n")
	width := 0
	for _, row := range rows {
		width = max(width, len([]rune(row)))
	}
	colours := lipgloss.Blend1D(width, lipgloss.Color(t.LogoFrom), lipgloss.Color(t.LogoTo))
	for i, row := range rows {
		var b strings.Builder
		for col, r := range []rune(row) {
			if r == ' ' {
				b.WriteRune(r)
				continue
			}
			b.WriteString(lipgloss.NewStyle().Foreground(colours[col]).Render(string(r)))
		}
		rows[i] = b.String()
	}
	return strings.Join(rows, "\n")
}

// END: PaintLogo
