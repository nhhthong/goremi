// The spectrum in the player panel: the bar height of a band level.
package ui

import (
	"math"
	"strings"

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

// SpectrumBars is the number of bands the model keeps and SpectrumRowCount the rows of the drawing; SpectrumShown is the number of bars drawn, each one column wide with one blank column after it.
const (
	SpectrumBars     = 32
	SpectrumRowCount = 8
	SpectrumShown    = 20
)

// SpectrumCells is the number of cells a bar of level 0 to 64 lights: one per 8 levels, rounded to the nearest, so a level of 4 lights one cell and a level of 3 none.
func SpectrumCells(level int) int { return min(max((level+4)/8, 0), SpectrumRowCount) }

// SpectrumRows draws the bars as 8 lines of PanelWidth columns: 20 bars, one column wide at the even columns, each a stack of half blocks ▄ that grows from the bottom by SpectrumCells of its level; bar j takes the highest of the bands from j×32/20 up to (j+1)×32/20.
func SpectrumRows(heights [SpectrumBars]int) []string {
	rows := make([]string, SpectrumRowCount)
	for r := range rows {
		line := make([]rune, PanelWidth)
		for i := range line {
			line[i] = ' '
		}
		for j := 0; j < SpectrumShown; j++ {
			level := 0
			for b := j * SpectrumBars / SpectrumShown; b < (j+1)*SpectrumBars/SpectrumShown; b++ {
				level = max(level, heights[b])
			}
			if SpectrumCells(level) > SpectrumRowCount-1-r {
				line[2*j] = '▄'
			}
		}
		rows[r] = string(line)
	}
	return rows
}

// SpectrumBaseline is the row of ─ under the bars, PanelWidth columns wide.
func SpectrumBaseline() string { return strings.Repeat("─", PanelWidth) }

// END: SpectrumRows

// START: NextIdle

// idleFlip is the chance, 1 in idleFlip per frame, that an idle bar changes between no cell and one cell.
const idleFlip = 10

// NextIdle gives the idle levels of the next frame while no audio plays: each bar is 0 or one cell (level 8) and flips with a chance of 1 in 10 per frame, so the bars move gently; intn(n) gives a number from 0 to n-1.
func NextIdle(prev [SpectrumBars]int, intn func(n int) int) [SpectrumBars]int {
	for i := range prev {
		if intn(idleFlip) == 0 {
			prev[i] = 8 - prev[i]
		}
	}
	return prev
}

// END: NextIdle

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
