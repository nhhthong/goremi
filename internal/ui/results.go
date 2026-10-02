// Results list rendering: one line per track, a selection marker, and a "load more..." line.
package ui

import (
	"strings"

	"charm.land/lipgloss/v2"

	"goremi/internal/provider"
	"goremi/internal/ui/theme"
)

// START: RenderResults

// RenderResults renders one line per track; the selected line starts with "▶ ".
func RenderResults(tracks []provider.Track, selected int) string {
	lines := make([]string, len(tracks))
	for i, t := range tracks {
		prefix := "  "
		if i == selected {
			prefix = "▶ "
		}
		lines[i] = prefix + t.Title
	}
	return strings.Join(lines, "\n")
}

// END: RenderResults

// START: RenderResultsMore

// RenderResultsMore renders the tracks followed by a selectable "load more..." line.
func RenderResultsMore(tracks []provider.Track, selected int) string {
	if len(tracks) == 0 {
		return ""
	}
	prefix := "  "
	if selected == len(tracks) {
		prefix = "▶ "
	}
	return RenderResults(tracks, selected) + "\n" + prefix + "load more..."
}

// END: RenderResultsMore

// START: PaintResults

// PaintResults colours a rendered list: the selected line bold Accent on Selected, "load more..." Muted, the track lines Foreground.
// tracks is the number of track lines; the lines after them are the "load more..." line.
func PaintResults(t theme.Theme, rendered string, selected, tracks int) string {
	lines := strings.Split(rendered, "\n")
	for i, l := range lines {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Foreground))
		switch {
		case i == selected:
			style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(t.Accent)).Background(lipgloss.Color(t.Selected))
		case i >= tracks:
			style = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted))
		}
		lines[i] = style.Render(l)
	}
	return strings.Join(lines, "\n")
}

// END: PaintResults
