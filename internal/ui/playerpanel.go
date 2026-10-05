// Parts of the player panel that show the playing track.
package ui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"

	"goremi/internal/ui/theme"
)

// PanelWidth is the width of the player panel in columns, the width of the logo.
const PanelWidth = 40

// START: ArtistLine

// ArtistLine is the artist row of the player panel: the artist, or Loading… while it loads.
func ArtistLine(t theme.Theme, text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(t.Foreground)).Render(text)
}

// END: ArtistLine

// START: TitleLine

// TitleLine is the title row of the player panel.
func TitleLine(t theme.Theme, text string) string {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(t.Accent)).Render(text)
}

// END: TitleLine

// START: ClockText

// ClockText is the clock row: elapsed and total as m:ss, each as h:mm:ss from one hour on.
func ClockText(elapsed, total time.Duration) string {
	return clock(elapsed) + " / " + clock(total)
}

// clock formats one time, in whole seconds.
func clock(d time.Duration) string {
	s := int(d / time.Second)
	if s >= 3600 {
		return fmt.Sprintf("%d:%02d:%02d", s/3600, s%3600/60, s%60)
	}
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

// END: ClockText

// START: ControlsLine

// ControlsText is the plain controls row: previous, seek back, play/pause (▷ while paused), seek forward, next.
func ControlsText(paused bool) string {
	if paused {
		return "◀◀ ◀ ▷ ▶ ▶▶"
	}
	return "◀◀ ◀ ❚❚ ▶ ▶▶"
}

// ControlsLine is the controls row painted with the theme.
func ControlsLine(t theme.Theme, paused bool) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Render(ControlsText(paused))
}

// END: ControlsLine

// START: BarLine

// BarLine is the progress row: panelWidth − 2 cells, ━ for the elapsed part, ● at its head (at most the last cell) and ─ for the rest. A total of 0 reads as no progress.
func BarLine(t theme.Theme, panelWidth int, elapsed, total time.Duration) string {
	cells := panelWidth - 2
	if cells <= 0 {
		return ""
	}
	filled := 0
	if total > 0 {
		filled = int(float64(cells) * float64(elapsed) / float64(total))
	}
	head := min(filled, cells-1)
	done := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Progress)).Render(strings.Repeat("━", head) + "●")
	rest := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Border)).Render(strings.Repeat("─", cells-1-head))
	return done + rest
}

// END: BarLine
