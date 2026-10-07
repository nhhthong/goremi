// The ASCII banner of the header: the art of ascii.txt, embedded, painted from the theme.
package ui

import (
	"strings"

	"charm.land/lipgloss/v2"

	"goremi/internal/ui/theme"
)

// START: Banner

// bannerRows is the art of ascii.txt, embedded so the installed binary needs no file.
var bannerRows = []string{
	"   ██████╗  ██████╗ ██████╗ ███████╗███╗   ███╗██╗",
	"  ██╔════╝ ██╔═══██╗██╔══██╗██╔════╝████╗ ████║██║",
	"  ██║  ███╗██║   ██║██████╔╝█████╗  ██╔████╔██║██║",
	"  ██║   ██║██║   ██║██╔══██╗██╔══╝  ██║╚██╔╝██║██║",
	"  ╚██████╔╝╚██████╔╝██║  ██║███████╗██║ ╚═╝ ██║██║",
	"   ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝     ╚═╝╚═╝",
}

// BannerWidth is the width of the art in columns; it is shown from that width.
const BannerWidth = 50

// Banner is the art, six rows of BannerWidth columns.
func Banner() string { return strings.Join(bannerRows, "\n") }

// BannerRows is the number of rows the banner takes at a width: all of them from BannerWidth columns, none below.
func BannerRows(width int) int {
	if width >= BannerWidth {
		return len(bannerRows)
	}
	return 0
}

// END: Banner

// START: PaintBanner

// PaintBanner colours the art column by column, from t.LogoFrom on the first column to t.LogoTo on the last; the cells of one column share a colour and spaces stay plain.
func PaintBanner(t theme.Theme) string {
	colours := lipgloss.Blend1D(BannerWidth, lipgloss.Color(t.LogoFrom), lipgloss.Color(t.LogoTo))
	rows := make([]string, len(bannerRows))
	for i, row := range bannerRows {
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

// END: PaintBanner
