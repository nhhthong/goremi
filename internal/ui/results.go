// Results list rendering: one line per track, a selection marker, and a "load more..." line.
package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

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

// START: Render

// Render renders the list: the tracks, then a "load more..." line while there is more to load, shown as "Loading..." while the next page loads.
func (r Results) Render() string {
	if !r.more {
		return RenderResults(r.tracks, r.selected)
	}
	out := RenderResultsMore(r.tracks, r.selected)
	if r.loading {
		out = strings.TrimSuffix(out, "load more...") + "Loading..."
	}
	return out
}

// END: Render

// START: RenderWidth

// RenderWidth is Render with every line cut to width columns, ending with "…" when it was cut; a width of 0 or less cuts nothing.
func (r Results) RenderWidth(width int) string {
	out := r.Render()
	if width <= 0 {
		return out
	}
	lines := strings.Split(out, "\n")
	for i, l := range lines {
		lines[i] = ansi.Truncate(l, width, "…")
	}
	return strings.Join(lines, "\n")
}

// END: RenderWidth
