// The spectrum in the player panel: the bar height of a band level.
package ui

import (
	"math"

	"charm.land/lipgloss/v2"

	"goremi/internal/ui/theme"
)

// START: SpectrumHeight

// SpectrumHeight maps a band level in dB onto a bar height: -60 dB is 0, 0 dB is 64 (8 rows of 8 levels), clamped outside.
func SpectrumHeight(db float64) int {
	if math.IsNaN(db) || db <= -60 {
		return 0
	}
	if db >= 0 {
		return 64
	}
	return int(math.Round((db + 60) / 60 * 64))
}

// END: SpectrumHeight

// START: SmoothHeight

// SmoothHeight is the height a bar shows next: it rises at once to a higher target and falls to 85 % of its old height each frame.
func SmoothHeight(prev, target int) int {
	if target >= prev {
		return target
	}
	return max(target, prev*85/100)
}

// END: SmoothHeight

// START: SpectrumRows

// SpectrumBars is the number of bars; SpectrumRowCount their rows and SpectrumOffset the blank cells on each side in the 40 columns of the panel.
const (
	SpectrumBars     = 32
	SpectrumRowCount = 8
	SpectrumOffset   = 4
)

// SpectrumRows draws the bars as 8 lines of PanelWidth columns: each bar is one cell wide, centred; a row holds 8 height levels, a full cell is █ and the top cell of a bar is one of ▁▂▃▄▅▆▇█ by what is left over.
func SpectrumRows(heights [SpectrumBars]int) []string {
	const blocks = "▁▂▃▄▅▆▇█"
	glyphs := []rune(blocks)
	rows := make([]string, SpectrumRowCount)
	for r := range rows {
		line := make([]rune, PanelWidth)
		for i := range line {
			line[i] = ' '
		}
		for b, h := range heights {
			n := min(max(h-(SpectrumRowCount-1-r)*8, 0), 8)
			if n > 0 {
				line[SpectrumOffset+b] = glyphs[n-1]
			}
		}
		rows[r] = string(line)
	}
	return rows
}

// END: SpectrumRows

// START: IdleHeights

// IdleHeights are the bar heights while no audio plays: each bar gets a random level of 0 to 5 % of the 64 levels, that is 0 to 3; intn(n) gives a number from 0 to n-1.
func IdleHeights(intn func(n int) int) [SpectrumBars]int {
	var h [SpectrumBars]int
	for i := range h {
		h[i] = intn(4)
	}
	return h
}

// END: IdleHeights

// START: PaintSpectrum

// PaintSpectrum colours the rows of the spectrum with the theme's Spectrum colour; the text is unchanged.
func PaintSpectrum(t theme.Theme, rows []string) []string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Spectrum))
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = style.Render(r)
	}
	return out
}

// END: PaintSpectrum
