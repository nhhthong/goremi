// The ASCII art "Goremi" shown in the player panel; made once from resources/ascii_terminal.png.
package ui

import (
	"strings"

	"charm.land/lipgloss/v2"

	"goremi/internal/ui/theme"
)

// START: Logo

// Logo is the art, 40 columns wide, one line per row.
var Logo = strings.Join([]string{
	"        ▟███▙                        ▐█▌",
	"    ▄  ▐█▛▀▜█▌▗▄▄▄ ▄▄▗▄ ▗▄▄▄ ▗▄▄▄▗▄▄ ▄▄",
	"  ▗▖█  ▐█     ████▌▜███▚████▙▜██████▌▜█▌",
	" ▗▐▌█▟ ▐█ ▐██▘█▌ █▌ █▌ ▐█▄▄██▐█▌██ █▌▐█▌",
	"▄▜▐▌██ ▐█  ▜█ █▌ █▌ █▌ ▐█▀▀▀▀▐█▌▜█ █▌▐█▌",
	"▛▐▐▌▌█ ▝██▄█▛ █▙▄█▌▟█▌ ▐█▙▄▄▖▟█▌██ █▌▟█▙",
	"        ▝▀▀▀  ▝▀▀▀ ▝▀▘  ▀▀▀▀ ▝▀▘▝▘ ▀▘▀▀▀",
}, "\n")

// END: Logo

// START: PaintLogo

// PaintLogo colours the art row by row, from t.LogoFrom on the first row to t.LogoTo on the last.
func PaintLogo(t theme.Theme) string {
	rows := strings.Split(Logo, "\n")
	colours := lipgloss.Blend1D(len(rows), lipgloss.Color(t.LogoFrom), lipgloss.Color(t.LogoTo))
	for i, row := range rows {
		rows[i] = lipgloss.NewStyle().Foreground(colours[i]).Render(row)
	}
	return strings.Join(rows, "\n")
}

// END: PaintLogo
